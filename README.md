# fuse-go

[![License](https://img.shields.io/badge/license-Apache%202-green.svg)](https://www.apache.org/licenses/LICENSE-2.0)

一个纯Golang的[libfuse](https://github.com/libfuse/libfuse)实现，目前相当于实现了libfuse中的low level接口。

## 如何使用

* `fuse.Session`是启动和管理整个Fuse运行周期的结构体，所以这个类的实例是必须要有的。
* `fuse.mount`中是挂载目录的函数封装，主要是`mount`和`umount`，具体实现其实是使用了`fusermount`的。
* `fuse.FuseOpt`是保存用户实现的方法的结构体，然后将其传入`fuse.Session`中。
* `fuse.util`中目前提供了两工具类:
    * `FusePathManager`是一个key: inode，val: filepath的字典
    * `NotExistManager`是用来缓存那些文件路径不存在的，可以设置一个超时时间

要实现的文件操作接口，可以查看[opt_h.go](./fuse/opt_h.go)，如果有些接口不需要实现，则直接不赋值(`nil`)即可。

### 示例代码

```golang
import (
    "github.com/mingforpc/fuse-go/fuse"
    "github.com/mingforpc/fuse-go/fuse/mount"
)

......

var getattr = func(req fuse.Req, nodeid uint64) (fsStat *fuse.Stat, result int32) {

	......

	return fsStat, result
}

......

func main() {
    ......
    opts := fuse.Opt{}
    opts.Getattr = &getattr
    opts.Opendir = &opendir
    opts.Readdir = &readdir
    opts.Releasedir = &release
    opts.Release = &release
    opts.Lookup = &lookup
    opts.Open = &open
    opts.Read = &read
    opts.Mkdir = &mkdir
    opts.Create = &create
    opts.Setattr = &setattr
    opts.Write = &write
    opts.Unlink = &unlink
    opts.Rmdir = &rmdir
    opts.Rename = &rename
    opts.Setxattr = &setxattr
    opts.Getxattr = &getxattr
    opts.Listxattr = &listxattr
    opts.Removexattr = &removexattr
    ......

    se := fuse.NewFuseSession(cg.Mountpoint, &opts, 1024)

    ......

    err := mount.Mount(se)

    ......

    se.FuseLoop()

    ......
}
```

## 原理

通过对`/dev/fuse`读写，获得内核fuse传过来的文件操作。

详细可以看看`fuse/kernel/`中的代码。

与`/dev/fuse`交互的请求响应格式，可以参考我的博文[FUSE协议解析](http://blog.mingforpc.me/2018/11/30/FUSE%E5%8D%8F%E8%AE%AE%E8%A7%A3%E6%9E%90/#more)。

Mount部分的代码，为了实现普通用户可以直接Mount，调用了`/bin/fusermount`（可能需要安装`libfuse`后才有），Mount的代码主要复制[go-fuse](https://github.com/hanwen/go-fuse)，然后自己做了小修改

## 例子

`example`中有一个`hello`的简单例子，具体使用是`go build`后，`hello -mp {挂载的目录}`。

更详细的例子可以参考我下面的`hadoop-fs`

[https://github.com/mingforpc/hadoop-fs](https://github.com/mingforpc/hadoop-fs)

## TODO(按优先级排序)

* ~~`errno`中的错误代码需要完善(急)~~(2018/12/21,已在`errno`中添加基本的错误代码)
* ~~修改好注释和调整结构~~(2018/12/24)
* ~~完成`example`中的例子~~(2018/12/29)
* ~~发现`Close()`会存在block的情况，需要解决(急)~~(2019/1/3)
* 往`test`中添加测试用例
* 测试各个系统上的兼容性(目前只在64位Ubuntu上进行)
* 完善`evloop`相关的操作
* 提供一些管理`inode`和`path`的工具类
* 实现一个高级接口

## 联系我

如果有什么问题直接联系我指出，万分感谢！

邮箱: jianmingforpc@gmail.com