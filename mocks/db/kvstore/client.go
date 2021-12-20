package kvstore

import (
	"context"
	"time"
)

const (
	// Maximum channel buffer between publisher/subscriber goroutines
	maxClientChannelBufferSize = 10
)

// These constants represent the event types returned by the KV client
const (
	PUT = iota
	DELETE
	CONNECTIONDOWN
	UNKNOWN
)

// KVPair is a common wrapper for key-value pairs returned from the KV store
type KVPair struct {
	Key     string
	Value   interface{}
	Version int64
	Session string
	Lease   int64
}

// NewKVPair creates a new KVPair object
func NewKVPair(key string, value interface{}, session string, lease int64, version int64) *KVPair {
	kv := new(KVPair)
	kv.Key = key
	kv.Value = value
	kv.Session = session
	kv.Lease = lease
	kv.Version = version
	return kv
}

// Event is generated by the KV client when a key change is detected
type Event struct {
	EventType int
	Key       interface{}
	Value     interface{}
	Version   int64
}

// NewEvent creates a new Event object
func NewEvent(eventType int, key interface{}, value interface{}, version int64) *Event {
	evnt := new(Event)
	evnt.EventType = eventType
	evnt.Key = key
	evnt.Value = value
	evnt.Version = version

	return evnt
}

// Client represents the set of APIs a KV Client must implement
type Client interface {
	List(ctx context.Context, key string) (map[string]*KVPair, error)
	Get(ctx context.Context, key string) (*KVPair, error)
	Put(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	DeleteWithPrefix(ctx context.Context, prefixKey string) error
	Watch(ctx context.Context, key string, withPrefix bool) chan *Event
	IsConnectionUp(ctx context.Context) bool // timeout in second
	CloseWatch(ctx context.Context, key string, ch chan *Event)
	Close(ctx context.Context)

	// These APIs are not used.  They will be cleaned up in release Voltha 2.9.
	// It's not cleaned now to limit changes in all components
	Reserve(ctx context.Context, key string, value interface{}, ttl time.Duration) (interface{}, error)
	ReleaseReservation(ctx context.Context, key string) error
	ReleaseAllReservations(ctx context.Context) error
	RenewReservation(ctx context.Context, key string) error
	AcquireLock(ctx context.Context, lockName string, timeout time.Duration) error
	ReleaseLock(lockName string) error
}