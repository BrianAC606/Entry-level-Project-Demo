package main

import (
	"ProjectDemo/ProducerAndConsumerDemo/ManyToMany"
	"fmt"
)

type T struct{}

func (t *T) f(i int) *T {
	fmt.Println(i)
	return t
}
func test() {
	t := new(T)
	defer t.f(1).f(2)
	defer t.f(3)
	fmt.Println(t.f(4))
}

func main() {
	//itemQueue := items.NewItems()
	//go itemQueue.Working()
	//itemQueue.ProductItem("gdfonesr")
	//itemQueue.ProductItem("htgr")
	//itemQueue.ProductItem("bressterger5ge")
	//itemQueue.ProductItem("juythetg")
	//itemQueue.ProductItem("zgtesr")
	//
	//time.Sleep(time.Second * 4)
	//close(itemQueue.Data)
	//select {}
	//OneToOne.OneToOne()
	//ManyToOne.ManyToOne()
	//OneToMany.OneToMany()
	ManyToMany.ManyToMany()

	//test()

}
