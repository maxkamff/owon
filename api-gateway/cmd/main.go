package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/casbin/gorm-adapter/v2"
	"github.com/gomodule/redigo/redis"

	"gitlab.com/api-gateway/api"
	"gitlab.com/api-gateway/config"
	"gitlab.com/api-gateway/pkg/logger"
	"gitlab.com/api-gateway/services"
	r "gitlab.com/api-gateway/storage/redis"


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
		"api-gateway",
	)
	if err != nil {
		fmt.Println(err)
	}
	defer closer.Close()

	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	a, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil{
		log.Error("New Adapter Error", logger.Error(err))
		return
	}

	casbinEnforcer, err := casbin.NewEnforcer(cfg.CasbinConfigPath, a)
	if err != nil{
		log.Error("new enforcer error", logger.Error(err))
		return
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil{
		log.Error("casbin load policy error", logger.Error(err))
		return
	}

	pool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("kayMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("kayMatch3", util.KeyMatch3)

	server := api.New(api.Option{
		Conf:            cfg,
		Logger:          log,
		ServiceManager:  serviceManager,
		InMemoryStorage: r.NewRedisRepo(&pool),
		CasbinEnforcer: casbinEnforcer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
