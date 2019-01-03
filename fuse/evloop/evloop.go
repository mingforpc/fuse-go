package evloop

// Handler event loop handler
type Handler func(el *EvLoop, fd int, eventmask int, privdata interface{})

// The struct to save handler and private data
type loopEvent struct {
	fd       int
	handler  Handler
	privdata interface{}
}

// EvLoop the event loop struct
type EvLoop struct {
	maxSize int

	epoll *epoll

	fdMap map[int]*loopEvent // key: fd, value: *loopEvent
}

// event struct
type event struct {
	events uint32
	fd     int32
}

// NewEvLoop create EvLoop
func NewEvLoop(maxSize int) EvLoop {

	ep := newEpoll(maxSize)
	evloop := EvLoop{maxSize: maxSize, epoll: &ep, fdMap: make(map[int]*loopEvent)}

	return evloop
}

// Register register event with fd
func (el *EvLoop) Register(fd int, eventmask int, handler Handler, privdata interface{}) error {

	op := EpollRegisterAdd
	if el.fdMap[fd] != nil {
		op = EpollRegisterMod
	}

	err := el.epoll.register(fd, eventmask, op)

	if err == nil {
		el.fdMap[fd] = &loopEvent{fd: fd, handler: handler, privdata: privdata}
	}

	return err
}

// UnRegister unregister event
func (el *EvLoop) UnRegister(fd int) error {

	if el.fdMap[fd] == nil {
		return nil
	}

	err := el.epoll.unregister(fd)

	if err == nil {
		delete(el.fdMap, fd)
	}

	return err
}

// Process event
// timeout argument specifies the number of milliseconds that wait() will block.
// the return value int is the number of event raised
func (el *EvLoop) Process(timeout int) (int, error) {

	events, err := el.epoll.wait(timeout)

	if err != nil {
		return 9, err
	}

	n := len(events)

	if events == nil || n == 0 {
		return 0, nil
	}

	// raise earch handler
	for _, event := range events {
		loopev := el.fdMap[int(event.fd)]
		loopev.handler(el, loopev.fd, int(event.events), loopev.privdata)
	}

	return n, nil
}
