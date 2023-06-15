package main

import (
	"fmt"
	"net"

	"gitlab.com/post-service/config"
	p "gitlab.com/post-service/genproto/post"
	"gitlab.com/post-service/pkg/db"
	"gitlab.com/post-service/pkg/logger"
	"gitlab.com/post-service/service"
	gp "gitlab.com/post-service/service/grpc_client"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

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
		"post-service",
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
		fmt.Println("Error connect postgres", err.Error())
	}

	client, err := gp.New(cfg)
	if err != nil {
		fmt.Println("Error while grpc client", err.Error())
	}

	PostService := service.NewPostService(connDB, log, *client)

	lis, err := net.Listen("tcp", cfg.PostServicePort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	p.RegisterPostServiceServer(s, PostService)

	log.Info("main: server running",
		logger.String("port", cfg.PostServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
