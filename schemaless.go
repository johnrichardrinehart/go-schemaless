package schemaless

import (
	"context"
	"github.com/dgryski/go-metro"
	jh "github.com/dgryski/go-shardedkv/choosers/jump"
	"github.com/rbastic/go-schemaless/core"
	"github.com/rbastic/go-schemaless/models"
	"sync"
)

// Storage is a key-value storage backend
type Storage interface {
	// Get the cell designated (row key, column key, ref key)
	Get(ctx context.Context, rowKey string, columnKey string, refKey int64) (cell models.Cell, found bool, err error)

	// GetLatest returns the latest value for a given rowKey and columnKey, and a bool indicating if the key was present
	GetLatest(ctx context.Context, rowKey string, columnKey string) (cell models.Cell, found bool, err error)

	// PartitionRead returns 'limit' cells after 'location' from shard 'shard_no'
	PartitionRead(ctx context.Context, partitionNumber int, location string, value uint64, limit int) (cells []models.Cell, found bool, err error)

	// Put inits a cell with given row key, column key, and ref key
	Put(ctx context.Context, rowKey string, columnKey string, refKey int64, cell models.Cell) (err error)

	// ResetConnection reinitializes the connection for the shard responsible for a key
	ResetConnection(ctx context.Context, key string) error

	// Cleans up any resources, etc.
	Destroy(ctx context.Context) error
}

// DataStore is our overall datastore structure, backed by at least one
// KVStore.
type DataStore struct {
	source *core.KVStore
	// we avoid holding the lock during a call to a storage engine, which
	// may block
	mu sync.Mutex
}

// Chooser maps keys to shards
type Chooser interface {
	// SetBuckets sets the list of known buckets from which the chooser should select
	SetBuckets([]string) error
	// Choose returns a bucket for a given key
	Choose(key string) string
	// Buckets returns the list of known buckets
	Buckets() []string
}

// Shard is a named storage backend
type Shard struct {
	Name    string
	Backend Storage
}

func hash64(b []byte) uint64 { return metro.Hash64(b, 0) }

func (ds *DataStore) WithSource(shards []core.Shard) *DataStore {
	chooser := jh.New(hash64)
	kv := core.New(chooser, shards)
	ds.source = kv
	return ds
}

func New() *DataStore {
	return &DataStore{}
}

func (ds *DataStore) Get(ctx context.Context, rowKey string, columnKey string, refKey int64) (cell models.Cell, found bool, err error) {
	return ds.source.Get(ctx, rowKey, columnKey, refKey)
}

func (ds *DataStore) GetLatest(ctx context.Context, rowKey string, columnKey string) (cell models.Cell, found bool, err error) {
	return ds.source.GetLatest(ctx, rowKey, columnKey)
}

func (ds *DataStore) PartitionRead(ctx context.Context, partitionNumber int, location string, value uint64, limit int) (cells []models.Cell, found bool, err error) {
	return ds.source.PartitionRead(ctx, partitionNumber, location, value, limit)
}

// Put
func (ds *DataStore) Put(ctx context.Context, rowKey string, columnKey string, refKey int64, cell models.Cell) error {
	return ds.source.Put(ctx, rowKey, columnKey, refKey, cell)
}

// ResetConnection implements Storage.ResetConnection()
func (ds *DataStore) ResetConnection(ctx context.Context, key string) error {
	return ds.source.ResetConnection(ctx, key)
}

// Destroy implements Storage.Destroy()
func (ds *DataStore) Destroy(ctx context.Context) error {
	return ds.source.Destroy(ctx)
}
