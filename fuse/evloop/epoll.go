package evloop

import (
	"errors"
	"syscall"
)

// Epoll register type, add or modify
const (
	EpollRegisterAdd = syscall.EPOLL_CTL_ADD
	EpollRegisterMod = syscall.EPOLL_CTL_MOD
)

const (
	EPOLLIN  = syscall.EPOLLIN
	EPOLLOUT = syscall.EPOLLOUT
	EPOLLERR = syscall.EPOLLERR
)

var errRegisterType = errors.New("worng register type, should be EpollRegisterAdd or EpollRegisterMod")

type epoll struct {
	fd int // epoll fd

	maxSize int                  // the max number in epoll
	events  []syscall.EpollEvent // the array to receive events
}

// newEpoll create epoll object
func newEpoll(maxSize int) epoll {
	fd, err := syscall.EpollCreate(1)
	if err != nil {
		panic(err)
	}

	ep := epoll{fd: fd, maxSize: maxSize, events: make([]syscall.EpollEvent, maxSize)}

	return ep
}

// register register event in epoll
// reType is EpollRegisterAdd or EpollRegisterMod
func (ep *epoll) register(fd int, eventmask int, reType int) error {
	event := syscall.EpollEvent{Fd: int32(fd), Events: uint32(eventmask)}

	if reType != EpollRegisterAdd && reType != EpollRegisterMod {
		panic(errRegisterType)
	}

	err := syscall.EpollCtl(ep.fd, reType, fd, &event)

	return err
}

// unregister fd in epoll
func (ep *epoll) unregister(fd int) error {

	err := syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_DEL, fd, nil)
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
