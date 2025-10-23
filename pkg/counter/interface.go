package counter

// ICounter defines a thread-safe counter interface
type ICounter interface {
	Increment(amount int) int // Increment counter by amount and return new value
	Counter() int             // Get current counter value
}
