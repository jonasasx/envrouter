package utils

type Observer interface {
	Subscribe(s *chan ObserverEvent)
	Unsubscribe(s *chan ObserverEvent)
	Publish(e *ObserverEvent)
}

type observer struct {
	subscribers []*chan ObserverEvent
}

func NewObserver() Observer {
	return &observer{}
}

func (o *observer) Subscribe(s *chan ObserverEvent) {
	o.subscribers = append(o.subscribers, s)
}

func (o *observer) Unsubscribe(s *chan ObserverEvent) {
	for i, v := range o.subscribers {
		if *s == *v {
			o.subscribers = append(o.subscribers[:i], o.subscribers[i+1:]...)
			break
		}
	}
}

func (o *observer) Publish(e *ObserverEvent) {
	for _, s := range o.subscribers {
		s := s
		go func() {
			*s <- *e
		}()
	}
}

type ObserverEvent struct {
	Item  interface{} `json:"item"`
	Event string      `json:"event"`
}
