package counter

type ICounter interface {
	Increment(amout int) int
	Counter() int
}
