package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runs_adapter/adapter"
	"runs_adapter/handlers"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	pemPath, keyPath, proto, host, listenport, logLevel, dbuser, dbpassword, dbhost, dbname string
	dbport                                                                                  int
	wait, writeTimeout, readTimeout, idleTimeout                                            time.Duration
	debug, logger                                                                           bool
)

func init() {
	runtime.GOMAXPROCS(1)
}
func main() {
	flag.StringVar(&pemPath, "pempath", os.Getenv("SERVERREM"), "path to pem file")
	flag.StringVar(&keyPath, "keypath", os.Getenv("SERVERKEY"), "path to key file")
	flag.StringVar(&listenport, "port", "9443", "port to Listen")
	flag.StringVar(&proto, "proto", "http", "http or https")
	flag.StringVar(&logLevel, "loglevel", "INFO", "log level, default INFO")
	flag.BoolVar(&logger, "logger", false, "Gin logger Mode")
	flag.BoolVar(&debug, "debug", false, "Gin debug Mode")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully")
	flag.DurationVar(&readTimeout, "read-timeout", time.Second*60, "read server timeout")
	flag.DurationVar(&writeTimeout, "write-timeout", time.Second*60, "write server timeout")
	flag.DurationVar(&idleTimeout, "idle-timeout", time.Second*60, "idle server timeout")
	flag.StringVar(&dbuser, "user", "postgres", "db user")
	flag.StringVar(&dbpassword, "password", `postgres`, "db user password")
	flag.StringVar(&dbhost, "dbhost", "172.27.192.59", "db host")
	flag.IntVar(&dbport, "dbport", 5432, "db port")
	flag.StringVar(&dbname, "dbname", "test", "db name")
	flag.Parse()

	if !debug {
		gin.SetMode(gin.ReleaseMode)
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.Error("Get interface adress error: ", err.Error())
			os.Exit(1)
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					host = fmt.Sprintf("%s:%s", ipnet.IP.String(), listenport)
				}
			}
		}

	} else {
		host = fmt.Sprintf("%s:%s", "0.0.0.0", listenport)
	}
	var router *gin.Engine
	if logger {
		log.Info("Use gin Default")
		router = gin.Default()
	} else {
		log.Info("Use gin new and recovery")
		router = gin.New()
		router.Use(gin.Recovery())
	}
	router.Use(handlers.Middleware())
	ppof := router.Group("/pprof")
	ppof.GET("/", gin.WrapF(pprof.Index))
	ppof.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	ppof.GET("/profile", gin.WrapF(pprof.Profile))
	ppof.POST("/symbol", gin.WrapF(pprof.Symbol))
	ppof.GET("/symbol", gin.WrapF(pprof.Symbol))
	ppof.GET("/trace", gin.WrapF(pprof.Trace))
	ppof.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	ppof.GET("/block", gin.WrapH(pprof.Handler("block")))
	ppof.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	ppof.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	ppof.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	ppof.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	router.GET("/metrics", handlers.PrometheusHandler())
	v1 := router.Group("/api/v1")
	{
		v1.Any("/spaces", handlers.Spaces)
		v1.Any("/projects", handlers.Projects)
		v1.Any("/releases", handlers.Releases)
		v1.Any("/versions", handlers.Versions)
		v1.Any("/testtypes", handlers.TestTypes)
		v1.Any("/profiles", handlers.Profiles)
		runs := v1.Group("/runs")
		{
			runs.Any("/", handlers.Runs)
			runs.GET("/runid", handlers.LastRunID)
			runs.GET("/space", handlers.SpaceRuns)
			runs.GET("/project", handlers.ProjectRuns)
			runs.GET("/release", handlers.ReleaseRuns)
			runs.GET("/version", handlers.VersionRuns)
			runs.GET("/testtype", handlers.TestTypeRuns)
		}
	}
	go func() {
		for {
			err := adapter.InitDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname))
			if err != nil {
				log.Error(err)
			} else {
				log.Info("Connected db complete")
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()

	srv := &http.Server{
		Addr:         host,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      router,
	}
	go func() {
		switch proto {
		case "https":
			log.Info("Server is run, proto: https, address: %s ", srv.Addr)
			if err := srv.ListenAndServeTLS(pemPath, keyPath); err != nil {
				log.Println(err)
			}
		case "http":
			log.Info("Server is run, proto: http, address: %s ", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Info("server shutting down")
	os.Exit(0)
}
