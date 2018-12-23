package kernel

import "errors"

// ErrDataLen the lenght of binary not enough to parse to data struct
var ErrDataLen = errors.New("Data len not enough")

// ErrNoNeedReply the kernel request no need to reply
var ErrNoNeedReply = errors.New("No need to reply")

// ErrNotInit fuse session not inited
var ErrNotInit = errors.New("Fuse session not inited")
