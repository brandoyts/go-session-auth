package auth

import "time"

type SessionMetadata struct {
	Key   string
	Value interface{}
	TTL   time.Time
}
