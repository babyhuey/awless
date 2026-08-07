package awsfetch

// maxParallelAWSCalls bounds the fan-out of per-resource AWS calls made while
// fetching.
//
// Previously these helpers spawned one goroutine per item with no limit, so an
// account with thousands of buckets or clusters issued thousands of
// simultaneous API calls and reliably hit throttling (RequestLimitExceeded).
//
// 20 is a conservative default that keeps sync fast without tripping per-region
// request limits.
const maxParallelAWSCalls = 20
