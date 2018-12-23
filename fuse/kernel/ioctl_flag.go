package kernel

// FuseIoctlCompat : 32bit compat ioctl on 64bit machine
const FuseIoctlCompat = (1 << 0)

// FuseIoctlUnrestricted : not restricted to well-formed ioctls, retry allowed
const FuseIoctlUnrestricted = (1 << 1)

// FuseIoctlRetry : retry with new iovecs
const FuseIoctlRetry = (1 << 2)

// FuseIoctl32bit : 32bit ioctl
const FuseIoctl32bit = (1 << 3)

// FuseIoctlDir : is a directory
const FuseIoctlDir = (1 << 4)

// FuseIoctlMaxIov : maximum of in_iovecs + out_iovecs
const FuseIoctlMaxIov = 256
