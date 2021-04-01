# Sample Microservice

Sample Golang Microservice

It comes pre-configured with :

1. Protocol Buffers (https://google.golang.org/protobuf)
2. gRPC(https://google.golang.org/grpc)
3. gRPC Open Tracing(https://github.com/grpc-ecosystem/grpc-opentracing)
4. Open Tracing (https://github.com/opentracing/opentracing-go)
5. Jaeger (https://github.com/uber/jaeger-lib)
6. Jaeger Client for Go (https://github.com/uber/jaeger-client-go)
7. Mongo DB (https://go.mongodb.org/mongo-driver)
8. Redigo (https://github.com/gomodule/redigo)
9. Cobra (https://github.com/spf13/cobra)
10. Viper (https://github.com/spf13/viper)


## Setup

Use this command to install the blueprint

```bash
go get github.com/rifqiakrm/go-microservice
```

or manually clone the repo and then run `go run main.go`.

## Quick Note

Before you start the main service, you may want to set your environtment variables. You can choose it on the config, fill the env key and then set the env path file on `cmd/root.go`

```bash
EXPORT PROJECT_ENV=config/config.{the choosen env}.toml
```

## Generate Protocol Buffer

Use this command to install generate code from protocol buffers :

```
protoc -I $GOPATH/src --go_out=$GOPATH/src $GOPATH/src/github.com/rifqiakrm/{project_name}/pb/{proto_dir}/{your_proto}.proto
protoc -I $GOPATH/src --go-grpc_out=$GOPATH/src $GOPATH/src/github.com/rifqiakrm/{project_name}/pb/{proto_dir}/{your_proto}.proto
```

protoc -I pb/user/ --go_out=pb/user/ pb/user/user.proto
## Step by step deploying to server?

1. Build the application
* The first step is build your application. You can build it manually with `go build`.
2. Setup the config file
* Create a directory on the server. You can create the directory anywhere but i suggest to put it on `/var/www/{your_directory}`. After creating the directory put `configs` files under your directory. 
3. Setup the service
* Create a file under `/etc/systemd/system` and name it `{application_name}.service`. Fill the file with this code below
* ```
    [Unit]
    Description= instance to serve api
    Requires=mongodb.service
    After=network.target
    After=mongodb.service
    [Service]
    WorkingDirectory=/var/www/{your_directory}
    User=root
    Group=www-data
    Environment="GOPATH=/var/www/"
    Environment="PROJECT_ENV=configs/{your_config_env}"
    ExecStart=/var/www/{application_build_file}
    [Install]
    WantedBy=multi-user.target
  ```
4. Allow UFW Port
* After Setting up your nginx conf don't forget to allow the port that not open on linux by default. You can run commad `ufw allow {port}` on the terminal. In this case we have to open port 50052, so the command will be `ufw allow 50052`
5. Run it
* Run the application with this command : `service {application_name} start` or `systemctl start {application_name}` and you're good to go!


Also Thanks to the internet and Tabvn for tutorial on how to deploy golang project to VPS server.

> https://medium.com/@tabvn/deploy-golang-application-on-digital-ocean-server-ubuntu-16-04-b7bf5340ccd9

