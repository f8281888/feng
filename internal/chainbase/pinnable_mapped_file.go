package chainbase

import (
	"feng/internal/fc/common"

	"golang.org/x/exp/mmap"
)

const (
	//Mapped ..
	Mapped = iota
	//Heap ..
	Heap = iota
	//Locked ..
	Locked = iota
)

//PinnableMappedFile ..
type PinnableMappedFile struct {
	mappedFileLock common.FileLock
	dataFilePath   string
	databaseName   string
	writable       bool
	//bip::file_mapping  _file_mapping; 文件内存映射
	fileMapping mmap.ReaderAt
	//Boost的提供了一套ipc的接口，内存映射文件将文件的内容映射到进程的地址空间。segment_manager字面意思为段管理器
	segmentManager *mmap.ReaderAt
}

//GetSegmentManager ..
func (c PinnableMappedFile) GetSegmentManager() *mmap.ReaderAt {
	// at, err := mmap.Open("/Users/xyz/test_mmap_data.txt")
	// p := []byte("nihao")
	// at.ReadAt(p, len(p))
	return c.segmentManager
}

//共享内存和文件内存映射的区别
// 首先是共享内存和文件内存映射的接口、用法不一样，GLIBC的POSIX的共享内存实现会默认把共享内存文件放在/dev/shm/分区下，
//如果没有这个分区，需要手动挂载一下。然后就是共享内存和文件内存映射，在内核中的实现原理，都使用了内核的cache和swap机制，完全没有区别。

// 我们的系统中，本来使用的是共享内存来进行进程间通信，共享内存文件的位置在/dev/shm/.但是后面要迁移至某平台，
// 而据说平台的容器中并没有这个分区，所以无法使用共享内存。大佬给支招，“可以用mmap内存映射，把文件映射到内存中，
// 和原来用共享内存差不多”。经过我一番折腾，发现，把原来使用boost库共享内存的接口改为使用文件内存映射，
// 一共改动不超过10行。用文件内存映射的方式运行程序，未见任何异常，而要映射的文件，却可以不用放在/dev/shm/下。
// 大佬又说：“我见过的系统基本都用的文件内存映射，没怎么见过用共享内存的，不知道为什么这个系统开始设计时候一定要用共享内存”。
