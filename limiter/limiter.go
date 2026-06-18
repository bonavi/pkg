package limiter

type Limiter struct {
	ch chan struct{}
}

func New(limit int) *Limiter {
	if limit <= 0 {
		panic("limiter: limit must be > 0")
	}
	return &Limiter{
		ch: make(chan struct{}, limit),
	}
}

func (l *Limiter) TryAcquire() bool {
	select {
	case l.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *Limiter) Release() {
	<-l.ch
}
