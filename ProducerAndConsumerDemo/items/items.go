package items

import (
	"fmt"
	"sync"
)

type Items struct {
	Data chan interface{}
}

var (
	mutex    sync.Mutex
	instance *Items
)

// 单例设计
func NewItems() *Items {
	if instance == nil {
		mutex.Lock()
		instance = new(Items)
		instance.Data = make(chan interface{}, 10)
		mutex.Unlock()
	}
	return instance
}

func (i *Items) ProductItem(item interface{}) {
	i.Data <- item
}

func (i *Items) Working() {
	for {
		select {
		case item := <-i.Data:
			fmt.Println(item)
		}
	}
}
