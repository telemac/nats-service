package counter

import (
	"sync"
)

// Counter is a thread-safe counter implementation
type Counter struct {
	counter int          // Current counter value
	mutex   sync.RWMutex // Read-write mutex for thread safety
}

// Increment adds amount to the counter and returns the new value
func (c *Counter) Increment(amount int) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.counter += amount
	return c.counter
}

// Counter returns the current counter value
func (c *Counter) Counter() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.counter
}
