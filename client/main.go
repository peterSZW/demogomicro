package main

import (
	"context"
	"fmt"
	_ "github.com/micro/go-plugins/registry/consul"

	proto "demogomicro/greeter" //这里写你的proto文件放置路劲
	micro "github.com/micro/go-micro"
	"time"
)

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()

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
