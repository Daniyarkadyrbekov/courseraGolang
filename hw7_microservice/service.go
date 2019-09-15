package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	"net"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type ACLConfig struct {
	ACLs map[string][]string `json:"-"`
}

type BizImplementation struct {}

func (b BizImplementation) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (b BizImplementation) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (b BizImplementation) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func newBizImplementation() BizImplementation{
	return BizImplementation{}
}

type AdminServerImplementation struct {

}

func (a AdminServerImplementation) Logging(*Nothing, Admin_LoggingServer) error {
	return nil
}

func (a AdminServerImplementation) Statistics(*StatInterval, Admin_StatisticsServer) error {
	return nil
}

func newAdminServer() AdminServerImplementation{
	return AdminServerImplementation{}
}

func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	ACL := ACLConfig{}
	if err := json.Unmarshal([]byte(ACLData), &ACL.ACLs); err != nil {
		return err
	}
	go startMicroservise(ctx, listenAddr, ACL)

	return nil
}

func startMicroservise (ctx context.Context, listenAddr string, ACL ACLConfig) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port:", err)
	}
	defer lis.Close()

	server := grpc.NewServer()

	RegisterAdminServer(server, newAdminServer())
	RegisterBizServer(server, newBizImplementation())

	go server.Serve(lis)

	<-ctx.Done()
}