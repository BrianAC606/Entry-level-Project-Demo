package OneToMany

import (
	"fmt"
	"sync"
)

type task struct {
	data int
}

type taskQueue struct {
	Data chan interface{}
}

var (
	mutex    sync.Mutex
	instance *taskQueue
)

func newTaskQueue() *taskQueue {
	if instance == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if instance == nil {
			instance = &taskQueue{
				Data: make(chan interface{}, 10),
			}
		}
	}
	return instance
}

var (
	taskNum     = 100
	consumerNum = 10
)

func producer(item chan<- interface{}) {
	for i := 0; i < taskNum; i++ {
		item <- task{data: i}
	}
	fmt.Println(">>>>>>>>>>>>producer done<<<<<<<<<<<<<<<<")
	defer close(item)
}

func consumer(item <-chan interface{}) {
	for i := range item {
		if task, ok := i.(task); ok {
			fmt.Println(task.data)
		}
	}
	fmt.Println("-------------consumer done---------------")
}

func OneToMany() {
	wg := sync.WaitGroup{}
	queue := newTaskQueue()

	go func() {
		wg.Add(1)
		defer wg.Done()
		producer(queue.Data)
	}()

	for i := 0; i < consumerNum; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			consumer(queue.Data)
		}(&wg)
	}

	wg.Wait()
}
