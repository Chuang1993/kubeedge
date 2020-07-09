package watch

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/kubeedge/beehive/pkg/core/model"
)

func newCache() *cache {
	return &cache{
		elems:    make([]model.Message, 100),
		startIdx: 0,
		endIdx:   0,
		capacity: 100,
	}
}

type cache struct {
	sync.RWMutex
	elems                []model.Message
	startIdx             int
	endIdx               int
	capacity             int
	oldestResouceVersion atomic.Value
}

func (c *cache) isFull() bool {
	return c.endIdx == c.startIdx+c.capacity
}
func (c *cache) Add(obj model.Message) {
	c.Lock()
	if c.isFull() {
		oldestobj := c.elems[c.startIdx%c.capacity]
		oldestRV, _ := strconv.ParseUint(oldestobj.Header.ResourceVersion, 0, 64)
		c.oldestResouceVersion.Store(oldestRV)
		c.startIdx++
	}
	c.elems[c.endIdx%c.capacity] = obj
	c.endIdx++
	c.Unlock()
}

func (c *cache) GetCacheMsg(resourceVersion uint64) ([]model.Message, error) {
	var oldestRV uint64
	v := c.oldestResouceVersion.Load()
	if v == nil {
		oldestRV = 0
	} else {
		oldestRV = v.(uint64)
	}

	if oldestRV > resourceVersion {
		return nil, fmt.Errorf("cache lost resouceversion relist")
	}
	size := c.endIdx - c.startIdx
	f := func(i int) bool {
		rv, _ := strconv.ParseUint(c.elems[(c.startIdx+i)%c.capacity].Header.ResourceVersion, 0, 64)
		return rv > resourceVersion
	}
	first := sort.Search(size, f)
	result := make([]model.Message, size-first)
	for i := 0; i < size-first; i++ {
		result[i] = c.elems[(c.startIdx+first+i)%c.capacity]
	}
	return result, nil
}
