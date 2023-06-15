package main

import (
	"fmt"
	"net"

	"gitlab.com/comment-service/config"
	c "gitlab.com/comment-service/genproto/comment"
	"gitlab.com/comment-service/pkg/db"
	"gitlab.com/comment-service/pkg/logger"
	"gitlab.com/comment-service/service"
	gp "gitlab.com/comment-service/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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
		"comment-service",
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

	CommentService := service.NewCommentService(connDB, log, *client)

	lis, err := net.Listen("tcp", cfg.CommentServicePort)

	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	c.RegisterCommentServiceServer(s, CommentService)

	log.Info("main: server running",
		logger.String("port", cfg.CommentServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
