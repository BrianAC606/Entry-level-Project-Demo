package OneToOne

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

var taskNum = 10000

func producer(tq chan<- interface{}) {
	for i := 0; i < taskNum; i++ {
		tq <- task{data: i}
	}
	defer close(tq)
}

func consumer(tq <-chan interface{}) {
	//labal:
	//	for {
	//		select {
	//		case t := <-tq:
	//			fmt.Println(t)
	//		case <-time.After(time.Second * 1):
	//			fmt.Println("time out, tasks are finished...")
	//			break labal
	//		}
	//	}

	for t := range tq {
		fmt.Println(t)
	}

	fmt.Println("consumer exit")
}

func OneToOne() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	queue := newTaskQueue()
	go func() {
		defer wg.Done()
		consumer(queue.Data)
	}()
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		producer(queue.Data)
	}(&wg)
	wg.Wait()
	//time.Sleep(time.Second * 3)
}
