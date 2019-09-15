package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

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
	go startMicroservise(ctx, listenAddr, ACLData)

	//lis, err := net.Listen("tcp", listenAddr)
	//if err != nil {
		//return err
		//log.Fatalln("cant listet port", err)
	//}
	//defer lis.Close()
	//
	//server := grpc.NewServer()
	//
	//RegisterAdminServer(server, newAdminServer())
	//
	//go server.Serve(lis)

	return nil
}

func startMicroservise (ctx context.Context, listenAddr string, ACLData string) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port:", err)
	}
	defer lis.Close()

	server := grpc.NewServer()

	RegisterAdminServer(server, newAdminServer())

	go server.Serve(lis)

	<-ctx.Done()
}