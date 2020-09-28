package main

import (
	"context"
	proto "demogomicro/greeter" //这里写你的proto文件放置路劲
	"fmt"
	micro "github.com/micro/go-micro"

	_ "github.com/micro/go-plugins/registry/consul"

	mqps "github.com/zengming00/go-qps"
	"math/rand"
	"net/http"
	"time"
)

var qps *mqps.QP

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = req.Name
	qps.Count()
	return nil
}

func main() {

	//Just for qps
	go qps_http_server()

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func qps_http_server() {
	rand.Seed(time.Now().UnixNano())

	// Statistics every second, a total of 3600 data
	qps = mqps.NewQP(time.Second, 3600)

	// Add a route to get HTML, for example /qps
	http.HandleFunc("/qps", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Get the raw HTML, you can gzip it
		s, err := qps.Show()
		if err != nil {
			panic(err)
		}
		w.Write([]byte(s))
	})
	// Add a route to get json report, The name is the same as getting the HTML routing, but you need to add the '_json' suffix
	http.HandleFunc("/qps_json", func(w http.ResponseWriter, r *http.Request) {
		// Get the json report
		bts, err := qps.GetJson()
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bts)
	})
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
}
