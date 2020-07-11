package stl

//Item ..
type Item interface {
}

//Queue ..
type Queue struct {
	items []Item
}

//QueueInterface ..
type QueueInterface interface {
	New() *Queue
	Enqueue(t Item)
	Dequeue() *Item
	IsEmpty() bool
	Size() int
}

//New 新建一个队列
func (q *Queue) New() *Queue {
	q = &Queue{}
	return q
}

//Enqueue 入队
func (q *Queue) Enqueue(t Item) {
	q.items = append(q.items, t)
}

//Dequeue 出队
func (q *Queue) Dequeue() *Item {
	items := q.items[0]
	q.items = q.items[1:len(q.items)]
	return &items
}

//IsEmpty 是否为空
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

//Size 大小
func (q *Queue) Size() int {
	return len(q.items)
}

//Clear ..
func (q *Queue) Clear() {
	q.items = q.items[0:]
}

//Get ..
func (q *Queue) Get(i uint32) Item {
	return q.items[i]
}
