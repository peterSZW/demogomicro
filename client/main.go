package main

import (
	"context"
	proto "demogomicro/greeter" //这里写你的proto文件放置路劲
	"flag"
	"fmt"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/registry/consul"
	"sync"
	"time"
)

var (
	topic = "mu.micro.book.topic.payment.done"
)

func broker_start() {
	//broker初始化
	if err := broker.Init(); err != nil {
		panic(err.Error())
	}
	if err := broker.Connect(); err != nil {
		panic(err.Error())
	}
	//异步调用broker的发布与订阅
	go publish()
	go subscribe()
}
func publish() {
	t := time.NewTicker(time.Second)

	for times := range t.C {
		err := broker.Publish(topic, &broker.Message{
			Header: map[string]string{
				"name": "gangan",
				"age":  "19",
			},
			Body: []byte(times.String()),
		})

		if err != nil {
			panic(err.Error())
		}
	}
}

func subscribe() {
	_, err := broker.Subscribe(topic, func(event broker.Event) error {
		b := event.Message().Body
		fmt.Println(string(b))

		h := event.Message().Header
		fmt.Println(h)
		return nil
	})

	if err != nil {
		panic(err.Error())
	}
}

var c int
var n int

var wg sync.WaitGroup

func main() {

	flag.IntVar(&c, "c", 10, "go routine number")
	flag.IntVar(&n, "n", 10000, "call times")

	// Create a new service. Optionally include some options here.
	// service := micro.NewService(micro.Name("greeter.client"))
	// service.Init()

	broker_start()

	// Create new greeter client

	fmt.Println("c:", c, "n:", n)

	for j := 0; j < c; j++ {
		wg.Add(1)
		go test(j)
	}
	wg.Wait()
	// Call the greeter

}
func test(j int) {

	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()

	greeter := proto.NewGreeterService("greeter", service.Client())
	var s string
	stime := time.Now()
	for i := 0; i < n; i++ {

		rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
		if err != nil {
			fmt.Println("greeter.Hello:", err)
		} else {
			s = rsp.Greeting
		}
		// Print response
		//fmt.Println(rsp.Greeting)
		// fmt.Println(j, i)

	}
	etime := time.Now()
	fmt.Println(etime.Sub(stime))
	fmt.Println(s)
	wg.Done()
}
