package evloop

import "syscall"

type epoll struct {
	fd int // epoll fd

	listenFd map[int]int // to save the fd

	maxSize int                  // the max number in epoll
	events  []syscall.EpollEvent // the array to receive events
}

type event struct {
	events uint32
	fd     int32
}

// newEpoll create epoll object
func newEpoll(maxSize int) epoll {
	fd, err := syscall.EpollCreate(1)
	if err != nil {
		panic(err)
	}

	ep := epoll{fd: fd, listenFd: make(map[int]int), maxSize: maxSize, events: make([]syscall.EpollEvent, maxSize)}

	return ep
}

// register register event in epoll
func (ep *epoll) register(fd int, eventmask int) error {
	event := syscall.EpollEvent{Fd: int32(fd), Events: uint32(eventmask)}

	op := syscall.EPOLL_CTL_ADD
	if ep.listenFd[fd] == 0 {
		ep.listenFd[fd] = eventmask
	} else {
		op = syscall.EPOLL_CTL_MOD
	}

	err := syscall.EpollCtl(ep.fd, op, fd, &event)

	return err
}

// unregister fd in epoll
func (ep *epoll) unregister(fd int) error {

	err := syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_DEL, fd, nil)

	if err == nil {
		delete(ep.listenFd, fd)
	}

	return err
}

// wait events
// timeout argument specifies the number of milliseconds that wait() will block.
func (ep *epoll) wait(timeout int) ([]event, error) {

	n, err := syscall.EpollWait(ep.fd, ep.events, timeout)

	if err != nil {
		return nil, err
	}

	if n > 0 {

		events := make([]event, n)

		for i := 0; i < n; i++ {
			writeEvent(ep.events[i], &events[i])
		}

		return events, nil
	} else {
		return nil, err
	}
}

// write syscall.EpollEvent value to event struct
func writeEvent(sysev syscall.EpollEvent, ev *event) {

	ev.fd = sysev.Fd
	ev.events = sysev.Events
}
