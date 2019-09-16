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
	"sync"
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
	eventChan chan Event
}

func (b BizImplementation) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := status.Error(codes.Unauthenticated, "can't get metadata from context")
		return nil, err
	}
	consumers := md.Get("consumer")
	consumer := consumers[0]
	event := Event{
		Timestamp: 0,
		Consumer: consumer,
		Method: "/main.Biz/Check",
		Host: "127.0.0.1:80821",
	}
	b.eventChan <- event

	return &Nothing{}, nil
}

func (b BizImplementation) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := status.Error(codes.Unauthenticated, "can't get metadata from context")
		return nil, err
	}
	consumers := md.Get("consumer")
	consumer := consumers[0]
	event := Event{
		Timestamp: 0,
		Consumer: consumer,
		Method: "/main.Biz/Add",
		Host: "127.0.0.1:80821",
	}
	b.eventChan <- event

	return &Nothing{}, nil
}

func (b BizImplementation) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	err := checkAccesToResource(ctx, b.ACL, "/main.Biz/Test")
	if err != nil {
		return nil, err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := status.Error(codes.Unauthenticated, "can't get metadata from context")
		return nil, err
	}
	consumers := md.Get("consumer")
	consumer := consumers[0]
	event := Event{
		Timestamp: 0,
		Consumer: consumer,
		Method: "/main.Biz/Check",
		Host: "127.0.0.1:80821",
	}
	b.eventChan <- event

	return &Nothing{}, nil
}

func newBizImplementation(ACL ACLConfig, eventChan chan Event) BizImplementation{
	return BizImplementation{ACL, eventChan}
}

type AdminServerImplementation struct {
	mu *sync.Mutex
	ACL ACLConfig
	eventChannel chan Event
	logListeners []Admin_LoggingServer
}

func (a AdminServerImplementation) Logging(nothing *Nothing, adminLogging Admin_LoggingServer) error {
	err := checkAccesToResource(adminLogging.Context(), a.ACL, "/main.Admin/Logging")
	if err != nil {
		return err
	}

	md, ok := metadata.FromIncomingContext(adminLogging.Context())
	if !ok {
		err := status.Error(codes.Unauthenticated, "can't get metadata from context")
		return err
	}
	consumers := md.Get("consumer")
	consumer := consumers[0]
	event := Event{
		Timestamp: 0,
		Consumer: consumer,
		Method: "/main.Admin/Logging",
		Host: "127.0.0.1:80822",
	}
	a.eventChannel <- event

	a.mu.Lock()
	defer a.mu.Unlock()
	a.logListeners = append(a.logListeners, adminLogging)

	return nil
}

func (a AdminServerImplementation) Statistics(interval *StatInterval, adminStatistic Admin_StatisticsServer) error {
	return nil
}

func newAdminServer(mu *sync.Mutex, ACL ACLConfig, listenAddr string, eventChan chan Event) AdminServerImplementation{
	return AdminServerImplementation{mu, ACL, eventChan, []Admin_LoggingServer{}}
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

	eventChan := make(chan Event)

	mu := &sync.Mutex{}
	adminServer := newAdminServer(mu, ACL, listenAddr, eventChan)
	RegisterAdminServer(server, adminServer)

	RegisterBizServer(server, newBizImplementation(ACL, eventChan))

	go server.Serve(lis)
	go adminServer.publishLogs(ctx, eventChan)
	<-ctx.Done()
}

func (a AdminServerImplementation)publishLogs(ctx context.Context, eventChan chan Event) {
	for{
		select {
		case <-ctx.Done():
			return
		case event := <- eventChan:
			a.sendToLogListners(event)
		}
	}
}

func (a AdminServerImplementation)sendToLogListners(event Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, listener := range a.logListeners{
		listener.Send(&event)
	}
}