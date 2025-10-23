package counter

import (
	"sync"
)

type Counter struct {
	counter int
	mutex   sync.RWMutex
}

func (c *Counter) Increment(amout int) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.counter += amout
	return c.counter
}

func (c *Counter) Counter() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.counter
}
