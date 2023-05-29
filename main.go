package main

import (
	"mall/process"
	"os"
	"os/signal"
	"syscall"
)

// var configFile = flag.String("f", "config/config.yml", "the config file")

// func main() {
// 	flag.Parse()
// 	c, err := config.Init(*configFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	serviceCtx := svc.NewServiceContext(c)
// 	e := routes.NewRouter(serviceCtx)
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%d", c.HttpServerConf.Port),
// 		Handler: e,
// 	}

// 	go func() {
// 		// 开启一个goroutine启动服务
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			serviceCtx.Log.Fatalf("ListenAndServe panic [%s]", err)
// 		}
// 	}()

// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
// 	<-sig

// 	_ = srv.Shutdown(context.TODO())
// }

func main() {
	process.Init()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sig
		process.Exit()
	}()

	process.Run()
}
