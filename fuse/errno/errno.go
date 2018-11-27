package errno

const SUCCESS = 0

const FAILED = -1
const ENOENT = -2  /* No such file or directory */
const EAGAIN = -11 /* Try again */
const EACCES = -13 /* Permission denied */
const EEXIST = -17 /* File exists */
const ENOTDIR = 20 /* Not a directory */
const ENOSYS = -38 /* Invalid system call number */
