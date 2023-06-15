package main

import (
	"fmt"
	"net"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go"

	"gitlab.com/user-service/config"
	u "gitlab.com/user-service/genproto/user"
	"gitlab.com/user-service/pkg/db"
	"gitlab.com/user-service/pkg/logger"
	"gitlab.com/user-service/service"
	gp "gitlab.com/user-service/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	conf := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 10,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	closer, err := conf.InitGlobalTracer(
		"user-service",
	)
	if err != nil {
		fmt.Println(err)
	}
	defer closer.Close()


	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "golang")
	defer logger.Cleanup(log)

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		fmt.Println("Error connect to postgres", err.Error())
	}
	
	client, err := gp.New(cfg)
	if err != nil {
		fmt.Println("Error while connecting grpc client", err.Error())
	}

	UserService := service.NewUserService(connDB, log, *client)

	lis, err := net.Listen("tcp", cfg.UserServicePort)

	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	u.RegisterUserServiceServer(s, UserService)

	log.Info("main: server running",
		logger.String("port", cfg.UserServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
