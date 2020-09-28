package main

import (
	"context"
	proto "demogomicro/greeter" //这里写你的proto文件放置路劲
	"fmt"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/registry/consul"
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

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()

	broker_start()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	// Call the greeter
	var s string
	stime := time.Now()
	for i := 0; i < 10000; i++ {

		rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
		if err != nil {
			fmt.Println(err)
		}
		// Print response
		//fmt.Println(rsp.Greeting)
		s = rsp.Greeting

	}
	etime := time.Now()
	fmt.Println(etime.Sub(stime))
	fmt.Println(s)
}
