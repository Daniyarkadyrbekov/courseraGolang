package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type ACLConfig struct {
	ACLs map[string][]string `json:"-"`
}

type BizImplementation struct {
	ACL ACLConfig
}

func (b BizImplementation) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (b BizImplementation) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (b BizImplementation) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	fmt.Printf("test method is worked")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("Unauthenticated bs no metadata")
		err := status.Error(codes.Unauthenticated, "can't get metadata from context")
		return nil, err //codes.Unauthenticated//errors.New("can't get metadata from context")
	}
	fmt.Printf("consumerVAlues = %v \n", md.Get("consumer"))

	//_, ok := b.ACL[ctx.Value()]
	return &Nothing{}, nil
}

func newBizImplementation(ACL ACLConfig) BizImplementation{
	return BizImplementation{ACL}
}

type AdminServerImplementation struct {
	ACL ACLConfig
}

func (a AdminServerImplementation) Logging(*Nothing, Admin_LoggingServer) error {
	return nil
}

func (a AdminServerImplementation) Statistics(*StatInterval, Admin_StatisticsServer) error {
	return nil
}

func newAdminServer(ACL ACLConfig) AdminServerImplementation{
	return AdminServerImplementation{ACL}
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
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("cant listen port:", err)
	}
	defer lis.Close()

	server := grpc.NewServer()

	RegisterAdminServer(server, newAdminServer(ACL))
	RegisterBizServer(server, newBizImplementation(ACL))

	go server.Serve(lis)
	<-ctx.Done()
}