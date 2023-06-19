package semaphore

import "fmt"

type Semaphore chan struct{}

func NewSemaphore(maxConcurrencyAmount int) Semaphore {
	s := make(Semaphore, maxConcurrencyAmount)

	return s
}

func (s *Semaphore) Acquire() {
	*s <- struct{}{}
	fmt.Println("Acquire")
}

func (s *Semaphore) Release() {
	<-*s
	fmt.Println("Release")
}
