package network

import (
	"feng/internal/fc/stl"
	"feng/internal/log"
)

//小知识
// s := "A1" // 分配存储"A1"的内存空间，s结构体里的str指针指向这快内存
// s = "A2" // 重新给"A2"的分配内存空间，s结构体里的str指针指向这快内存
// 其实[]byte和string的差别是更改变量的时候array的内容可以被更改。

// s := []byte{1} // 分配存储1数组的内存空间，s结构体的array指针指向这个数组。
// s = []byte{2} // 将array的内容改为2
// 因为string的指针指向的内容是不可以更改的，所以每更改一次字符串，就得重新分配一次内存，之前分配空间的还得由gc回收，这是导致string操作低效的根本原因。

//MessageBuffer ..
type MessageBuffer struct {
	Index stl.Pair
	//:message_buffer<1024*1024> 1048576
	//BufferLen ..
	BufferLen uint32
	//std::deque<std::array<char, buffer_len>* > buffers; 用双向链表代替合适？好像也不是很合适，因为用到了下标操作
	//[]byte 和string 该怎么选， 选[]byte 更高效灵活
	//[][]byte 跟 []bytes.Buffer 比较呢
	buffers [][]byte
	//pair 的用法 这里的pair, 一开始想到这里的下标识用来标识三维数组，但是好像又不太对
	//也有可能是二维数组的，[first][second]，因为second可能没写完
	readInd     stl.Pair
	writeInd    stl.Pair
	sanityCheck uint32
}

//BytesToWrite ..
//返回可写入多少数据，可以用 AddBufferToChain AddSpace 扩容
//比如有 4 * 1024*1024 ，如果第一个写满了， first 1, second 是1024*1024中的某个位置，如998
//可写的则为 4*1024*1024 - 1024*1024 - 998
func (m *MessageBuffer) BytesToWrite() uint32 {
	return m.totalBytes() - (m.readInd.First).(uint32)*m.BufferLen - (m.readInd.Second).(uint32)
}

//*buffer_len - write_ind.second
//总共的buf 长度 为 1024*1024 * []byte
func (m *MessageBuffer) totalBytes() uint32 {
	return m.BufferLen * uint32(len(m.buffers))
}

func (m *MessageBuffer) advanceIndex(index *stl.Pair, bytes uint32) {
	first := (index.First).(uint32)
	second := (index.Second).(uint32)
	first += (bytes + second) / m.BufferLen
	index.SetFirst(first)
	second = (bytes + second) % m.BufferLen
	index.SetSecond(second)
}

//下面这步应该就是用来设置first和second 位置的
// //AdvanceWritePtr ..
// func (m *MessageBuffer) AdvanceWritePtr(bytes uint32) {
// 	m.advanceIndex(&m.writeInd, bytes)
// 	for {
// 		if (m.writeInd.First).(uint32) < (uint32)(len(m.buffers)) {
// 			break
// 		}

// 		m.sanityCheck++
// 		m.buffers = append(m.buffers, m.malloc())
// 	}
// }

//AdvanceWritePtr ..
func (m *MessageBuffer) AdvanceWritePtr(l uint32) {
	if l <= 0 {
		log.Assert("AsyncReadToWrite write is zero")
	}

	first := (m.writeInd.First).(uint32)
	if int(first) > len(m.buffers) {
		log.Assert("AsyncReadToWrite buffer size is error")
	}

	second := (m.writeInd.Second).(uint32)
	iSecond := second
	lenB := l
	if lenB < second {
		iSecond += uint32(lenB)
	} else {
		i := first
		for {
			if lenB <= 0 {
				break
			}

			var end uint32 = 0
			var start uint32
			start = end
			end = ((m.BufferLen - second) + m.BufferLen*(i-first))
			if end > l {
				end = l
				iSecond += end
			} else {
				i++
				iSecond = end
				//开辟一个新内存空间出来
				m.buffers = append(m.buffers, m.malloc())
				m.sanityCheck++
			}

			lenB -= end - start
		}

		m.writeInd.SetFirst(i)
	}

	m.writeInd.SetSecond(iSecond)
}

//AddBufferToChain ..
func (m *MessageBuffer) AddBufferToChain() {
	m.sanityCheck++
	m.buffers = append(m.buffers, m.malloc())
}

func (m *MessageBuffer) malloc() []byte {
	return make([]byte, m.BufferLen)
}

//BytesToRead ..
func (m *MessageBuffer) BytesToRead() uint32 {
	return m.bytesToReadFromIndex(m.readInd)
}

func (m *MessageBuffer) bytesToReadFromIndex(ind stl.Pair) uint32 {
	return ((m.writeInd.First).(uint32)-(ind.First).(uint32))*m.BufferLen + (m.writeInd.Second).(uint32) - (ind.Second).(uint32)
}

//ReadIndex ..
func (m *MessageBuffer) ReadIndex() stl.Pair {
	return m.readInd
}

//memcpy(s, get_ptr(index), size);
//拷贝给一个uint32 用意是什么？ 是为了和1024*1024*4*2 比较 8388608

//Peek ..
//大概知道peek 为什么传len 进来， 包头有4个字节（int）是来存长度的，估计就是这个值
func (m *MessageBuffer) Peek(c *[]byte, size uint32, index *stl.Pair) bool {
	if m.bytesToReadFromIndex(*index) < size {
		log.Assert("tried to peek %d but only %d left", size, m.bytesToReadFromIndex(*index))
	}

	if (index.Second).(uint32)+size < m.BufferLen {
		//memcpy(s, get_ptr(index), size);
		copy(*c, m.getPtr(*index)[:size])
		m.advanceIndex(index, size)
	} else {
		numInBuffer := m.BufferLen - (index.Second).(uint32)
		copy(*c, m.getPtr(*index)[:numInBuffer])
		m.advanceIndex(index, numInBuffer)
		m.Peek(c, size-numInBuffer, index)
	}

	return true
}

// std::deque<std::array<char, buffer_len>* > buffers;
//return &buffers[index.first]->at(index.second);
//返回位置为n的元素的引用
func (m *MessageBuffer) getPtr(index stl.Pair) []byte {
	first := (index.First).(uint32)
	second := (index.Second).(uint32)
	a := m.buffers[first][second:]
	return a
}

//MbPeekDatastream ..
type MbPeekDatastream struct {
	buffeLen uint32
	mb       *MessageBuffer
	index    stl.Pair
}

//CreatePeekDatastream ..
func (m *MessageBuffer) CreatePeekDatastream() MbPeekDatastream {
	t := MbPeekDatastream{mb: m, index: m.ReadIndex(), buffeLen: m.BufferLen}
	return t
}

//AddSpace ..
func (m *MessageBuffer) AddSpace(bytes uint32) {
	buffersToAdd := bytes/m.BufferLen + 1
	for i := 0; i < int(buffersToAdd); i++ {
		m.sanityCheck++
		m.buffers = append(m.buffers, m.malloc())
	}
}

// //GetBufferSequenceForAsyncRead ..
// func (m *MessageBuffer) GetBufferSequenceForAsyncRead() *[]byte {
// 	//常见的转换方法，boost::asio::buffer 是配合起来使用的，这里有必要这么写么
// 	// boost::asio::mutable_buffer b1 =boost::asio::buffer(str)；
// 	// unsigned char* p1 = boost::asio::buffer_cast<unsigned char*>(b1);
// 	if (m.writeInd.First).(uint32) > uint32(m.buffers.Size()) {
// 		log.Assert("buffer size is error")
// 	}

// }

//AsyncReadToWrite .. 返回一个boost::asio::buffer 缓冲区去读取数据，但是go 没办法去这么操作
//比如 第N个slice， 把第N个slice 的第 N个剩余字节 加上后面的slice 去读取TCP数据，
//只能反向思维，把读到的数据写到这里。。。
func (m *MessageBuffer) AsyncReadToWrite(b []byte) {
	if len(b) <= 0 {
		log.Assert("AsyncReadToWrite write is zero")
	}

	first := (m.writeInd.First).(uint32)
	if int(first) > len(m.buffers) {
		log.Assert("AsyncReadToWrite buffer size is error")
	}

	second := (m.writeInd.Second).(uint32)
	iSecond := second
	lenB := len(b)
	if lenB < int(second) {
		copy(m.buffers[first], b[0:])
		iSecond += uint32(lenB)
	} else {
		i := first
		for {
			if lenB <= 0 {
				break
			}

			var end uint32 = 0
			var start uint32
			start = end
			end = ((m.BufferLen - second) + m.BufferLen*(i-first))
			if int(end) > len(b) {
				end = uint32(len(b))
				iSecond += end
			} else {
				i++
				iSecond = end
				//开辟一个新内存空间出来
				m.buffers = append(m.buffers, m.malloc())
				m.sanityCheck++
			}

			copy(m.buffers[i], b[start:end])
			lenB -= int(end - start)
		}

		// m.writeInd.SetFirst(i)
	}

	// m.writeInd.SetSecond(iSecond)
}

//InitIndex ..
func (m *MessageBuffer) InitIndex() {
	var zero uint32 = 0
	m.writeInd.SetFirst(zero)
	m.writeInd.SetSecond(zero)
	m.readInd.SetFirst(zero)
	m.readInd.SetSecond(zero)
	m.Index.SetFirst(zero)
	m.Index.SetSecond(zero)
	m.sanityCheck = 1
	m.buffers = append(m.buffers, m.malloc())
}

//Get ..
func (m *MbPeekDatastream) Get(c *[]byte) bool {
	return m.mb.Peek(c, 1, &m.index)
}
