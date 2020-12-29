package brainless

type Brain interface {
	Setup()
	Step()
	CheckDone() bool
	GetTask([][]int)
	ToResponse() interface{}
}
