package errno

const SUCCESS = 0

const FAILED = -1
const ENOENT = -2          /* No such file or directory */
const EAGAIN = -11         /* Try again */
const EACCES = -13         /* Permission denied */
const EEXIST = -17         /* File exists */
const ENOTDIR = 20         /* Not a directory */
const ERANGE = -34         /* Math result not representable, in getxattr: The size of the value buffer is too small to hold the result. */
const ENOSYS = -38         /* Invalid system call number */
const ENODATA = -61        /* No data available */
const ENOATTR = ENODATA    /* XATTR_REPLACE was specified, and the attribute does not exist. */
const EOPNOTSUPP = -95     /* Operation not supported on transport endpoint */
const ENOTSUP = EOPNOTSUPP /* setxattr: The namespace prefix of name is not valid. */
