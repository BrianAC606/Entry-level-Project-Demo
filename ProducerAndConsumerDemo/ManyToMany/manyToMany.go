package ManyToMany

import (
	"fmt"
	"sync"
	"time"
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
	taskNum = 50
	worKNum = 10
)

func producer(item chan<- interface{}, Done chan struct{}) {
	var i int
	for {
		if i >= taskNum {
			i = 0
		}
		i++
		t := task{data: i}
		select {
		case item <- t:
		case <-Done:
			fmt.Println("________________producer DONE_________________")
			return
		}
	}
}

func consumer(item <-chan interface{}, Done chan struct{}) {
	for {
		select {
		case i := <-item:
			fmt.Println(i)
		case <-Done:
			for v := range item {
				fmt.Println(v)
				fmt.Println("还剩余item......")
			}
			fmt.Println("________________consumer DONE_________________")
		}
	}
}

func ManyToMany() {
	Done := make(chan struct{})
	queue := newTaskQueue()
	for i := 0; i < taskNum; i++ {
		go producer(queue.Data, Done)
		go consumer(queue.Data, Done)
	}
	time.Sleep(time.Second * 5)
	close(Done)
	close(queue.Data)
}
