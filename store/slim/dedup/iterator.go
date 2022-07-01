package dedup

import (
	"github.com/pokt-network/pocket-core/types"
	dbm "github.com/tendermint/tm-db"
)

var _ types.Iterator = &DedupIterator{}
var _ types.KVStore = &Store{}
var _ types.CommitStore = &Store{}

type DedupIterator struct {
	parent dbm.DB
	it     dbm.Iterator
}

func NewDedupIterator(parent dbm.DB, height int64, prefix string, startKey, endKey []byte, isReverse bool) (dedupIterator *DedupIterator, err error) {
	start := HeightKey(height, prefix, startKey)
	end := HeightKey(height, prefix, endKey)
	dedupIterator = &DedupIterator{parent: parent}
	if isReverse {
		dedupIterator.it, err = parent.ReverseIterator(start, end)
	} else {
		dedupIterator.it, err = parent.Iterator(start, end)
	}
	return
}

// NewDedupIteratorForHeight : is the same as the above function, except exclusively for iterating an entire
// 'height'. It doesn't have any special functionality, rather a convenient way to initialize the iterator.
func NewDedupIteratorForHeight(parentDB dbm.DB, height int64, prefix string) *DedupIterator {
	startKey := HeightKey(height, prefix, nil)
	endKey := types.PrefixEndBytes(startKey)
	it := &DedupIterator{parent: parentDB}
	it.it, _ = parentDB.Iterator(startKey, endKey)
	return it
}

func (d *DedupIterator) Value() (value []byte) {
	dataStoreKey := d.it.Value()
	if dataStoreKey == nil {
		return
	}
	value, err := d.parent.Get(dataStoreKey)
	if err != nil {
		panic("an error occurred in dedup iterator value call: " + err.Error())
	}
	return
}

func (d *DedupIterator) Next()                              { d.it.Next() }
func (d *DedupIterator) Key() (key []byte)                  { return KeyFromHeightKey(d.it.Key()) }
func (d *DedupIterator) Error() error                       { return d.it.Error() }
func (d *DedupIterator) Close()                             { d.it.Close() }
func (d *DedupIterator) Valid() bool                        { return d.it.Valid() }
func (d *DedupIterator) Domain() (start []byte, end []byte) { panic("domain() not implemented") }

// NewOrphanIteratorForHeight : iterators for orphans (only needed when pruning), the only magic here is the Key() function
// which unwraps the actual key() value from the orphan key
func NewOrphanIteratorForHeight(parentDB dbm.DB, height int64, prefix string) *orphanIt {
	startKey := OrphanPrefixKey(HeightKey(height, prefix, nil))
	endKey := types.PrefixEndBytes(startKey)
	it := &orphanIt{parent: parentDB}
	it.it, _ = parentDB.Iterator(startKey, endKey)
	return it
}

func (d *orphanIt) Next()                    { d.it.Next() }
func (d *orphanIt) Key() (key []byte)        { return KeyFromHeightKey(KeyFromOrphanKey(d.it.Key())) }
func (d *orphanIt) Value() (value []byte)    { panic("should never call value on orphan iterator") }
func (d *orphanIt) Error() error             { return d.it.Error() }
func (d *orphanIt) Close()                   { d.it.Close() }
func (d *orphanIt) Valid() bool              { return d.it.Valid() }
func (d *orphanIt) Domain() ([]byte, []byte) { panic("domain not impelented()") }

type orphanIt struct {
	parent dbm.DB
	it     dbm.Iterator
}
