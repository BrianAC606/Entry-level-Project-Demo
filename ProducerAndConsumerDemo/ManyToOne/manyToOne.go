package ManyToOne

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
	taskNum    = 100
	produceNum = 10
)

func producer(ch chan<- interface{}, startNum int, num int) {
	for i := 0; i < num; i++ {
		ch <- task{
			data: startNum + i,
		}
	}
	fmt.Println("-----producer done------")
}

func consumer(ch <-chan interface{}) {
	for v := range ch {
		fmt.Println(v)
	}
	fmt.Println(">>>>>consumer done<<<<<")
}

func ManyToOne() {
	pwg := sync.WaitGroup{}
	wg := sync.WaitGroup{}
	queue := newTaskQueue()
	offsets := 0
	if taskNum%produceNum > 0 {
		offsets++
	}
	pwg.Add(taskNum/produceNum + offsets)
	//这里会有一定的全局变量共享问题
	var i int
	for i = 0; i < taskNum; i += produceNum {
		go func(i int) {
			defer pwg.Done()
			producer(queue.Data, i, produceNum)
		}(i)
	}
	go func() {
		wg.Add(1)
		defer wg.Done()
		consumer(queue.Data)
	}()
	pwg.Wait()
	close(queue.Data)
	wg.Wait()
}
