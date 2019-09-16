package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"regexp"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type ACLConfig struct {
	ACLs map[string][]string `json:"-"`
}

func checkAccesToResource(ctx context.Context, ACL ACLConfig, resource string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := status.Error(codes.Unauthenticated, "can't get metadata from context")
		return err
	}
	consumers := md.Get("consumer")
	consumerHasAccess := false
	for _, consumer := range consumers{
		regs, ok := ACL.ACLs[consumer]
		if !ok {
			err := status.Error(codes.Unauthenticated, "no consumers in context")
			return err
		}
		for _, reg := range  regs{
			ok, _ := regexp.MatchString(reg, resource)
			if ok {
				consumerHasAccess = true
			}
		}
	}

	if consumerHasAccess{
		return nil
	}
	err := status.Error(codes.Unauthenticated, "consumer doesn't have access")
	return err

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
	err := checkAccesToResource(ctx, b.ACL, "/main.Biz/Test")
	if err != nil {
		return nil, err
	}
	return &Nothing{}, nil
}

func newBizImplementation(ACL ACLConfig) BizImplementation{
	return BizImplementation{ACL}
}

type AdminServerImplementation struct {
	ACL ACLConfig
}

func (a AdminServerImplementation) Logging(nothing *Nothing, adminLogging Admin_LoggingServer) error {
	err := checkAccesToResource(adminLogging.Context(), a.ACL, "/main.Admin/Logging")
	if err != nil {
		return err
	}



	return nil
}

func (a AdminServerImplementation) Statistics(interval *StatInterval, adminStatistic Admin_StatisticsServer) error {
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