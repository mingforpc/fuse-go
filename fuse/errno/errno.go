package errno

const SUCCESS = 0

const FAILED = -1
const ENOENT = -2  /* No such file or directory */
const EACCES = -13 /* Permission denied */
const ENOTDIR = 20 /* Not a directory */
const ENOSYS = -38 /* Invalid system call number */
