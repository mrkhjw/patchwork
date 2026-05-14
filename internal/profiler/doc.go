// Package profiler provides lightweight execution profiling for patchwork operations.
//
// It wraps arbitrary functions with timing instrumentation, collecting duration
// and error information for each named operation. Results can be queried for the
// slowest operations, overall failure status, and total elapsed time.
//
// Typical usage:
//
//	p := profiler.New()
//	_ = p.Track("apply patch", func() error { return patch.Apply(...) })
//	_ = p.Track("run hooks", func() error { return hooks.RunAll(...) })
//	fmt.Print(profiler.Format(p))
package profiler
