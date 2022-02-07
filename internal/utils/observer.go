package utils

type Observer interface {
	Subscribe(h ObserverEventHandler)
	Unsubscribe(h ObserverEventHandler)
	Publish(oldObj interface{}, newObj interface{})
}

type observer struct {
	handlers []ObserverEventHandler
}

func NewObserver() Observer {
	return &observer{}
}

func (o *observer) Subscribe(h ObserverEventHandler) {
	o.handlers = append(o.handlers, h)
}

func (o *observer) Unsubscribe(h ObserverEventHandler) {
	for i, v := range o.handlers {
		if h == v {
			o.handlers = append(o.handlers[:i], o.handlers[i+1:]...)
			break
		}
	}
}

func (o *observer) Publish(oldObj interface{}, newObj interface{}) {
	for _, s := range o.handlers {
		s.OnEvent(oldObj, newObj)
	}
}

type ObserverEventHandler interface {
	OnEvent(oldObj interface{}, newObj interface{})
}

type ObserverEventHandlerFuncs struct {
	EventFunc func(oldObj interface{}, newObj interface{})
}

func (f ObserverEventHandlerFuncs) OnEvent(oldObj interface{}, newObj interface{}) {
	f.EventFunc(oldObj, newObj)
}
