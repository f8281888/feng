package netplugin

import (
	"sync"
)

type queueWrite struct {
	buffer   *[]byte
	callback func(error, int)
}

//QueueBuffer ..
type QueueBuffer struct {
	queueWrite
	mtx            sync.Mutex
	writeQueueSize uint32
	writeQueue     *[][]byte
	syncWriteQueue *[][]byte
	outQueue       *[][]byte
}

//WriteQueueSize ..
func (q *QueueBuffer) WriteQueueSize() uint32 {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return q.writeQueueSize
}

//ClearWriteQueue ..
func (q *QueueBuffer) ClearWriteQueue() {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.syncWriteQueue = new([][]byte)
}
