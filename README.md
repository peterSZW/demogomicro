# demo go-micro

## 安装 protoc
```
下载 https://github.com/google/protobuf/releases
make && make install 
编译安装
```

## protoc-gen-go
```
go get -u github.com/golang/protobuf/protoc-gen-go
```

## protoc-gen-micro
```
go get github.com/micro/protoc-gen-micro
```

## 生成 pb.go, pb.micro.go
```
cd D:\source.go\demogomicro\greeter
protoc --proto_path=D:\source.go\demogomicro\greeter --micro_out=. --go_out=. greeter.proto
```

## 启动consul（可忽略,新版使用 mdns）
```
consul agent -server  -node demogomicro-server -bind=127.0.0.1  -data-dir D:\tmp  -ui
```

## 启动 server
```
go build
./demogomicro
.\demogomicro.exe --registry=consul
 #consul

浏览器访问 http://127.0.0.1:8181/qps 查看qps
```

## 启动 client

```
cd client
go build
./client
``` 