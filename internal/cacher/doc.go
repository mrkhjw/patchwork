// Package cacher provides a lightweight thread-safe in-memory TTL cache
// intended for short-lived memoisation of patch evaluation results within
// a single patchwork run.
//
// Typical usage:
//
//	c := cacher.New(30 * time.Second)
//
//	if v, ok := c.Get(patchName); ok {
//		return v.(state.Entry), nil
//	}
//	result, err := expensiveOp(patchName)
//	if err == nil {
//		c.Set(patchName, result)
//	}
//
// Entries are evicted lazily on Get or explicitly via Evict. The cache does
// not start a background goroutine, keeping it dependency-free and easy to
// embed in larger pipelines.
package cacher
