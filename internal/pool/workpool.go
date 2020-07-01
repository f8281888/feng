package pool

import (
	"context"
	"sync"
)

//Task ..
type Task struct {
	ID  int
	Err error
	f   func() error
}

//Do ..
func (task *Task) Do() error {
	return task.f()
}

//WorkPool ..
type WorkPool struct {
	PoolSize    int
	tasksSize   int
	tasksChan   chan Task
	resultsChan chan Task
	Results     func() []Task
	wg          *sync.WaitGroup
	ioc         context.Context
}

//GetExecutor ..
func (pool WorkPool) GetExecutor() context.Context {
	return pool.ioc
}

func (pool *WorkPool) results() []Task {
	tasks := make([]Task, pool.tasksSize)
	for i := 0; i < pool.tasksSize; i++ {
		tasks[i] = <-pool.resultsChan
	}

	pool.wg.Wait()
	return tasks
}

//Wait ..
func (pool WorkPool) Wait() {
	pool.wg.Wait()
}

//NewPool ..
func NewPool(size int) *WorkPool {
	newWg := &sync.WaitGroup{}
	newWg.Add(size)
	c := context.Background()
	pool := &WorkPool{
		PoolSize: size,
		wg:       newWg,
		ioc:      c,
	}

	return pool
}

//AddTask ..
func (pool *WorkPool) AddTask(tasks []Task) {
	tasksChan := make(chan Task, len(tasks))
	resultsChan := make(chan Task, len(tasks))
	for _, task := range tasks {
		tasksChan <- task
	}

	close(tasksChan)
	pool.tasksSize = len(tasks)
	pool.tasksChan = tasksChan
	pool.resultsChan = resultsChan
	pool.Results = pool.results
}

func (pool *WorkPool) worker() {
	for task := range pool.tasksChan {
		task.Err = task.Do()
		pool.resultsChan <- task
		pool.wg.Done()
	}
}

//Start ..
func (pool *WorkPool) Start() {
	for i := 0; i < pool.PoolSize; i++ {
		go pool.worker()
	}
}

//用法
// t := time.Now()

// tasks := []Task{
// 	{Id: 0, f: func() error { time.Sleep(2 * time.Second); fmt.Println(0); return nil }},
// 	{Id: 1, f: func() error { time.Sleep(time.Second); fmt.Println(1); return errors.New("error") }},
// 	{Id: 2, f: func() error { fmt.Println(2); return errors.New("error") }},
// }
// pool := NewWorkerPool(tasks, 2)
// pool.Start()

// tasks = pool.Results()
// fmt.Printf("all tasks finished, timeElapsed: %f s\n", time.Now().Sub(t).Seconds())
// for _, task := range tasks {
// 	fmt.Printf("result of task %d is %v\n", task.Id, task.Err)
// }
