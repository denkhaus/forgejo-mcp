package label

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"codeberg.org/goern/forgejo-mcp/v2/pkg/forgejo"
	"codeberg.org/goern/forgejo-mcp/v2/pkg/log"

	forgejo_sdk "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	forgejo_models "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// LabelResolver provides label name-to-ID resolution functionality
type LabelResolver interface {
	ResolveLabelIDs(ctx context.Context, owner, repo string, labels []string) ([]int64, []ResolvedLabel, error)
}

type labelResolverImpl struct {
	// cache stores fetched labels per repository to avoid repeated API calls
	// map key: "owner/repo" -> cachedLabels with both name->ID and ID->name mappings
	cache map[string]*cachedLabels
	mu    sync.RWMutex
}

// cachedLabels stores both name-to-ID and ID-to-name mappings for a repository
type cachedLabels struct {
	nameToID map[string]int64
	idToName map[int64]string
}

// NewLabelResolver creates a new LabelResolver instance
func NewLabelResolver() LabelResolver {
	return &labelResolverImpl{
		cache: make(map[string]*cachedLabels),
	}
}

// ResolveLabelIDs resolves a list of label strings to their numeric IDs.
// Each label string can be either:
//   - A numeric ID (e.g., "47")
//   - A label name (e.g., "ready-to-merge")
// Returns the resolved IDs and the resolved label metadata
func (p *labelResolverImpl) ResolveLabelIDs(ctx context.Context, owner, repo string, labels []string) ([]int64, []ResolvedLabel, error) {
	log.Debugf("Resolving %d labels for %s/%s", len(labels), owner, repo)

	labelIDs := make([]int64, 0, len(labels))
	resolvedLabels := make([]ResolvedLabel, 0, len(labels))

	// Fetch all labels for this repository once
	nameToID, idToName, err := p.fetchRepositoryLabels(ctx, owner, repo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch repository labels: %w", err)
	}

	for _, labelStr := range labels {
		labelStr = strings.TrimSpace(labelStr)

		// Skip empty strings
		if labelStr == "" {
			return nil, nil, NewResolutionError("", "empty label not allowed", nil)
		}

		// Try parsing as numeric ID first
		labelID, err := strconv.ParseInt(labelStr, 10, 64)
		if err == nil {
			// It's a valid numeric ID
			if labelID <= 0 {
				return nil, nil, NewResolutionError(labelStr, "ID must be positive", nil)
			}
			labelIDs = append(labelIDs, labelID)

			// Try to find the name for this ID for the response
			if name := p.findNameByID(idToName, labelID); name != "" {
				resolvedLabels = append(resolvedLabels, ResolvedLabel{ID: labelID, Name: name})
			} else {
				resolvedLabels = append(resolvedLabels, ResolvedLabel{ID: labelID, Name: labelStr})
			}
			continue
		}

		// Not a numeric ID, try lookup by name (case-insensitive)
		labelID, found := p.lookupByName(nameToID, labelStr)
		if !found {
			// Label not found, provide suggestions
			suggestions := p.getSuggestions(nameToID, labelStr)
			return nil, nil, NewResolutionError(
				labelStr,
				fmt.Sprintf("not found in repository %s/%s", owner, repo),
				suggestions,
			)
		}

		labelIDs = append(labelIDs, labelID)
		resolvedLabels = append(resolvedLabels, ResolvedLabel{ID: labelID, Name: labelStr})
	}

	log.Debugf("Resolved %d labels successfully", len(labelIDs))
	return labelIDs, resolvedLabels, nil
}

// fetchRepositoryLabels fetches all labels for a repository, using cache if available
func (p *labelResolverImpl) fetchRepositoryLabels(ctx context.Context, owner, repo string) (map[string]int64, map[int64]string, error) {
	cacheKey := fmt.Sprintf("%s/%s", owner, repo)

	// Check cache first (read lock)
	p.mu.RLock()
	if cached, exists := p.cache[cacheKey]; exists {
		p.mu.RUnlock()
		log.Debugf("Using cached labels for %s", cacheKey)
		return cached.nameToID, cached.idToName, nil
	}
	p.mu.RUnlock()

	// Not in cache, fetch from API (write lock)
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if cached, exists := p.cache[cacheKey]; exists {
		return cached.nameToID, cached.idToName, nil
	}

	log.Debugf("Fetching labels from API for %s", cacheKey)

	// Fetch all labels (no pagination limit to get all labels)
	labels, _, err := forgejo.Client().ListRepoLabels(owner, repo, forgejo_sdk.ListLabelsOptions{
		ListOptions: forgejo_sdk.ListOptions{
			Page:     1,
			PageSize: 100, // Use max page size
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list labels: %w", err)
	}

	// Build label maps for both directions
	nameToID := make(map[string]int64)
	idToName := make(map[int64]string)
	for _, label := range labels {
		nameToID[strings.ToLower(label.Name)] = label.ID
		idToName[label.ID] = label.Name
	}

	// Cache the result
	p.cache[cacheKey] = &cachedLabels{
		nameToID: nameToID,
		idToName: idToName,
	}

	log.Debugf("Cached %d labels for %s", len(labels), cacheKey)
	return nameToID, idToName, nil
}

// lookupByName looks up a label ID by name (case-insensitive)
func (p *labelResolverImpl) lookupByName(labelMap map[string]int64, name string) (int64, bool) {
	id, found := labelMap[strings.ToLower(name)]
	return id, found
}

// findNameByID finds a label name by ID (reverse lookup)
func (p *labelResolverImpl) findNameByID(idToName map[int64]string, id int64) string {
	if name, exists := idToName[id]; exists {
		return name
	}
	return ""
}

// getSuggestions provides label name suggestions for a given input
func (p *labelResolverImpl) getSuggestions(nameToID map[string]int64, input string) []string {
	inputLower := strings.ToLower(input)
	var suggestions []string

	for name := range nameToID {
		// Check for partial match or similar spelling
		if strings.Contains(name, inputLower) || strings.Contains(inputLower, name) {
			suggestions = append(suggestions, name)
		}
	}

	// Limit suggestions to 3 most relevant
	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return suggestions
}

// ClearCache clears the label cache (useful for testing or force refresh)
func (p *labelResolverImpl) ClearCache() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache = make(map[string]*cachedLabels)
}

// GetCacheSize returns the number of repositories currently cached
func (p *labelResolverImpl) GetCacheSize() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.cache)
}

// FetchRepositoryLabels is a helper function that fetches all labels for a repository
// This can be used externally if needed
func FetchRepositoryLabels(ctx context.Context, owner, repo string) ([]*forgejo_models.Label, error) {
	labels, _, err := forgejo.Client().ListRepoLabels(owner, repo, forgejo_sdk.ListLabelsOptions{
		ListOptions: forgejo_sdk.ListOptions{
			Page:     1,
			PageSize: 100,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list labels: %w", err)
	}
	return labels, nil
}
