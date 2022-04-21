package signal

type Sub struct {
	ch chan struct{}
}

func newSub() *Sub {
	return &Sub{
		ch: make(chan struct{}, 1),
	}
}

func (s *Sub) Channel() <-chan struct{} {
	return s.ch
}
