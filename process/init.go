package process

import (
	"context"
	"flag"
	"fmt"
	"mall/config"
	"mall/routes"
	"mall/svc"
	"net/http"
	"os"
	"sync"

	skillconsumer "mall/service/skill_consumer"
)

var (
	server *http.Server

	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	serviceCtx *svc.ServiceContext
)

var configFile = flag.String("f", "config/config.yml", "the config file")

func Init() {
	flag.Parse()
	c, err := config.Init(*configFile)
	fmt.Println(c.HttpServerConf)
	if err != nil {
		panic(err)
	}

	serviceCtx = svc.NewServiceContext(c)
	ctx, cancel = context.WithCancel(context.TODO())
}

func Run() {
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		runHttpServer()
	}(&wg)

	wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		runOrtherTask()
	}(&wg)
	wg.Wait()
}

func Exit() {
	exitHttpServer()
	exitOrtherTask()
	serviceCtx.Close()
}

func runHttpServer() {
	e := routes.NewRouter(serviceCtx)

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", serviceCtx.Config.HttpServerConf.Port),
		Handler: e,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		serviceCtx.Log.Fatalf("ListenAndServe panic [%s]", err)
		os.Exit(1)
	}
}

func exitHttpServer() {
	server.Shutdown(context.TODO())
}

func runOrtherTask() {
	skillconsumer.SkillReqHandleTask(ctx, serviceCtx)
}

func exitOrtherTask() {
	cancel()
}
