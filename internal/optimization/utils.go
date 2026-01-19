package optimization

import (
	"sort"
)

// GroupOptimizations groups optimizations by ResourceAddress and returns a sorted list of keys.
func GroupOptimizations(opts []OptimizationRecommendation) (map[string][]OptimizationRecommendation, []string) {
	grouped := make(map[string][]OptimizationRecommendation)
	for _, opt := range opts {
		grouped[opt.ResourceAddress] = append(grouped[opt.ResourceAddress], opt)
	}

	keys := make([]string, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return grouped, keys
}
