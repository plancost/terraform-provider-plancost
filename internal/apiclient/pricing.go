// Copyright 2021 Infracost Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiclient

import (
	"fmt"
	"math"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/schema"

	"github.com/tidwall/gjson"
)

type PricingAPIClient struct {
	APIClient

	cache *lru.TwoQueueCache[uint64, cacheValue]
}

type cacheValue struct {
	Result    gjson.Result
	ExpiresAt time.Time
}

type PriceQueryKey struct {
	Resource      *schema.Resource
	CostComponent *schema.CostComponent
}

type PriceQueryResult struct {
	PriceQueryKey
	Result gjson.Result
	Query  GraphQLQuery

	filled bool
}

type BatchRequest struct {
	keys    []PriceQueryKey
	queries []GraphQLQuery
}

func (c *PricingAPIClient) buildQuery(product *schema.ProductFilter, price *schema.PriceFilter, currency string) GraphQLQuery {
	if currency == "" {
		currency = "USD"
	}

	v := map[string]interface{}{}
	v["productFilter"] = product
	v["priceFilter"] = price

	query := fmt.Sprintf(`
		query($productFilter: ProductFilter!, $priceFilter: PriceFilter) {
			products(filter: $productFilter) {
				prices(filter: $priceFilter) {
					priceHash
					%s
					termLength
				}
			}
		}
	`, currency)

	return GraphQLQuery{query, v}
}

// BatchRequests batches all the queries for these resources so we can use less GraphQL requests
// Use PriceQueryKeys to keep track of which query maps to which sub-resource and price component.
func (c *PricingAPIClient) BatchRequests(resources []*schema.Resource, batchSize int, currency string) []BatchRequest {
	reqs := make([]BatchRequest, 0)

	keys := make([]PriceQueryKey, 0)
	queries := make([]GraphQLQuery, 0)

	for _, r := range resources {
		for _, component := range r.CostComponents {
			keys = append(keys, PriceQueryKey{r, component})
			queries = append(queries, c.buildQuery(component.ProductFilter, component.PriceFilter, currency))
		}

		for _, subresource := range r.FlattenedSubResources() {
			for _, component := range subresource.CostComponents {
				keys = append(keys, PriceQueryKey{subresource, component})
				queries = append(queries, c.buildQuery(component.ProductFilter, component.PriceFilter, currency))
			}
		}
	}

	for i := 0; i < len(queries); i += batchSize {
		keysEnd := int64(math.Min(float64(i+batchSize), float64(len(keys))))
		queriesEnd := int64(math.Min(float64(i+batchSize), float64(len(queries))))

		reqs = append(reqs, BatchRequest{keys[i:keysEnd], queries[i:queriesEnd]})
	}

	return reqs
}

type pricingQuery struct {
	hash  uint64
	query GraphQLQuery

	result gjson.Result
}

// PerformRequest sends a batch request to the Pricing API endpoint to fetch
// pricing details for the provided queries. It optimizes the API call by
// checking a local cache for previous results. If the results of a given query
// are cached, they are used directly; otherwise, a request to the API is made.
func (c *PricingAPIClient) PerformRequest(req BatchRequest) ([]PriceQueryResult, error) {
	logging.Logger.Debug().Msgf("Getting pricing details for %d cost components from %s", len(req.queries), c.endpoint)
	res := make([]PriceQueryResult, len(req.keys))
	for i, key := range req.keys {
		res[i].PriceQueryKey = key
	}

	queries := make([]pricingQuery, len(req.queries))
	for i, query := range req.queries {
		key, err := hashstructure.Hash(query, hashstructure.FormatV2, nil)
		if err != nil {
			logging.Logger.Debug().Err(err).Msgf("failed to hash query %s will use nil hash", query)
		}

		queries[i] = pricingQuery{
			hash:  key,
			query: query,
		}

		res[i].Query = query
	}

	// first filter any queries that have been stored in the cache. We don't need to
	// send requests for these as we already have the results in memory.
	var serverQueries []pricingQuery
	if c.cache == nil {
		serverQueries = queries
	} else {
		var hit int
		for i, query := range queries {
			v, ok := c.cache.Get(query.hash)
			if ok {
				logging.Logger.Debug().Msgf("cache hit for query hash: %d", query.hash)
				hit++
				res[i].Result = v.Result
				res[i].filled = true
			} else {
				serverQueries = append(serverQueries, query)
			}
		}

		logging.Logger.Debug().Msgf("%d/%d queries were built from cache", hit, len(queries))
	}

	// now we deduplicate the queries, ensuring that a request for a price only happens once.
	deduplicatedServerQueries := make([]pricingQuery, 0, len(serverQueries))
	seenQueries := map[uint64]bool{}
	for _, query := range serverQueries {
		if seenQueries[query.hash] {
			continue
		}

		deduplicatedServerQueries = append(deduplicatedServerQueries, query)
		seenQueries[query.hash] = true
	}

	// send the deduplicated queries to the pricing API to fetch live prices.
	rawQueries := make([]GraphQLQuery, len(deduplicatedServerQueries))
	for i, query := range deduplicatedServerQueries {
		rawQueries[i] = query.query
	}
	resultsFromServer, err := c.DoQueries(rawQueries)
	if err != nil {
		return []PriceQueryResult{}, err
	}

	// if the cache is enabled lets store each pricing result returned in the cache.
	if c.cache != nil {
		for i, query := range deduplicatedServerQueries {
			if len(resultsFromServer)-1 >= i {
				(*c.cache).Add(query.hash, cacheValue{Result: resultsFromServer[i], ExpiresAt: time.Now().Add(time.Hour * 24)})
			}
		}
	}

	// now lets match the results from the server to their initial deduplicated queries.
	for i, result := range resultsFromServer {
		deduplicatedServerQueries[i].result = result
	}

	// Then we match deduplicated server queries to the initial list using the unique
	// query hash to tie a query to it's deduped query.
	resultMap := make(map[uint64]gjson.Result, len(deduplicatedServerQueries))
	for _, query := range deduplicatedServerQueries {
		resultMap[query.hash] = query.result
	}

	for i, query := range serverQueries {
		serverQueries[i].result = resultMap[query.hash]
	}

	// finally let's use the server queries to fill any results that haven't been
	// already populated from the cache.
	var x int
	for i, re := range res {
		if !re.filled {
			res[i].Result = serverQueries[x].result
			x++
		}
	}

	return res, nil
}
