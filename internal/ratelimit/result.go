package ratelimit

type ResultRateLimit struct {
	Allowed   bool
	Limit     int
	Remaining int
	ResetAt   int64
}
