package semaphore

type Semaphore chan struct{}

func NewSemaphore(maxConcurrencyAmount int) *Semaphore {
	s := make(Semaphore, maxConcurrencyAmount)

	return &s
}

func (s *Semaphore) Acquire() {
	*s <- struct{}{}
}

func (s *Semaphore) Release() {
	<-*s
}
