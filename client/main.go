package main

import (
	"context"
	proto "demogomicro/greeter" //这里写你的proto文件放置路劲
	"flag"
	"fmt"
	log "github.com/inconshreveable/log15"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/registry/consul"

	"runtime/debug"
	"sync"
	"time"
)

var (
	topic = "mu.micro.book.topic.payment.done"
)

func broker_start() {

	//broker初始化
	if err := broker.Init(); err != nil {
		fmt.Println(err)
		//panic(err.Error())
	}
	if err := broker.Connect(); err != nil {
		fmt.Println(err)
		//panic(err.Error())
	}
	//异步调用broker的发布与订阅

	subscribe()
	go publish()

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
			log.Crit("subscribe", "err", err)

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
		log.Crit("subscribe", "err", err)
	}
}

var c int
var n int

var wg sync.WaitGroup
var service micro.Service

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Crit("main", "server crash: ", err)
			log.Crit("main", "stack: ", string(debug.Stack()))
		}
	}()

	flag.IntVar(&c, "c", 10, "go routine number")
	flag.IntVar(&n, "n", 10000, "call times")

	// Create a new service. Optionally include some options here.
	service = micro.NewService(micro.Name("greeter.client"))
	service.Init()

	// Start broker subscribe and publish thread
	broker_start()

	if true {
		fmt.Println("c:", c, "n:", n)
		for j := 0; j < c; j++ {
			wg.Add(1)
			// Call the greeter
			go CallMicroSerice(j)
		}
	}

	// wait all process end
	wg.Wait()

}
func CallMicroSerice(j int) {

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())
	var s string
	stime := time.Now()
	for i := 0; i < n; i++ {

		// Call the greeter
		rsp, err := greeter.Hello(context.Background(), &proto.HelloRequest{Name: "John"})
		if err != nil {
			fmt.Println("greeter.Hello:", j, i, err)
		} else {
			// save response
			s = rsp.Greeting
		}
		// fmt.Println(j, i)
	}
	etime := time.Now()
	fmt.Println(etime.Sub(stime))
	fmt.Println(s)
	wg.Done()
}
