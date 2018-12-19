# fuse-go

一个简单的libfuse实现，用来学习Golang的语法和好奇libfuse是怎么工作的

相当于只实现了libfuse的low level接口

## 原理

通过对`/dev/fuse`读写，获得内核fuse传过来的文件操作。

传过来的前40bytes为Header:

```
Len     uint32
Opcode  uint32
Unique  uint64
Nodeid  uint64
Uid     uint32
Gid     uint32
Pid     uint32
Padding uint32
```

* Opcode: 指明对文件的操作
* Nodeid: 针对哪个文件的nodeid
* Len: 整个命令的长度(包括前40bytes的header和后面的内容)
* Unique: 是命令的唯一id

详细可以看看`fuse/kernel/`中的代码。

Mount部分的代码，为了实现普通用户可以直接Mount，调用了`/bin/fusermount`（可能需要安装`libfuse`后才有），Mount的代码主要复制[go-fuse](https://github.com/hanwen/go-fuse)，然后自己做了小修改

## 例子

可以参考:[https://github.com/mingforpc/hadoop-fs](https://github.com/mingforpc/hadoop-fs)

## TODO

* 修改好注释和调整结构
* 完成例子和测试
* 实现一个高级接口