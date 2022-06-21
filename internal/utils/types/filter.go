package types

import "time"

const (
	FilterRateLimiting = "RATE-LIMITING"
	FilterOAuthValid   = "OAUTH-VALID"
	FilterLeaseValid   = "LEASE-VALID"
	FilterFlowThresh   = "FLOW-THRESHOLD"
)

const (
	DefaultQPS      = 100
	DefaultDuration = 1 * time.Second
)
