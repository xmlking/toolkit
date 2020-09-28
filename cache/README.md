# Cache

Thread safety and LRU eviction are two basic things we wanted from a cache system. 
As we require an open-source solution –  HashiCorp’s [golang-lru](https://github.com/hashicorp/golang-lru) looked like a good fit. 
It’s an open source cache library which implements a fixed-size thread-safe LRU cache.

We mapped following concerns we wanted to address, building on top of golang-lru:

1. Concurrent misses (i.e. [“Thundering Herd”](https://en.wikipedia.org/wiki/Thundering_herd_problem))
   
    If we miss the cache we must hit the database
   
    Since we know many concurrent requests will likely be made by the same user, this means a cache miss will trigger many concurrent DB round-trips, adversely affecting our overall [tail latency](https://engineering.linkedin.com/performance/who-moved-my-99th-percentile-latency).
   
    We wanted to reduce the load on the database in this case, so we implemented a locking mechanism where for a given key, only one Goroutine will acquire the data from the database while the rest will simply block and wait for that request to return and populate the cache.

2. Eviction time
   
   If a user suddenly generates high load on the system, we don’t want a fixed expiry duration since all servers are likely to evict that user at the same time, which will in turn cause a spike in DB reads (see “thundering herd” above).

Risks we mapped and considered as a non-issue for our usage:

1. When DB lookup fails or returns a Not Found Error, we do not set any state into the cache, which will make the rest of the routines report that the item was not found.
1. Cache thrashing: If the LRU cache is too small, there’s a race condition here: it is possible to successfully add a key to the cache that will be immediately evicted, causing the subsequent lookup to fail.

## Features
- [x] Thread safe
- [x] LRU eviction
- [ ] Eviction time

## Usage

see cache [test-case](./cache_test.go)

```go
func TestCatalogService_Cache(t *testing.T) {
	log.Info().Msg("called TestCatalogService_Cache")
	catSrv := NewCatalogService(ServiceCache{Enabled: true, Size: 5})
	gotProduct, err := catSrv.GetProductByID("abc")
	gotProduct2, err2 := catSrv.GetProductByID("abc")

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.Equal(t, "abc", gotProduct.id)
	assert.Equal(t, gotProduct, gotProduct2)
}
```

## Reference
- [In-process caching in Go: scaling lakeFS to 100k requests/second](https://lakefs.io/2020/09/23/in-process-caching-in-go-scaling-lakefs-to-100k-requests-second/)
- [sync.Map](https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c)
- lakeFS's [cache](https://github.com/treeverse/lakeFS/tree/master/cache)
