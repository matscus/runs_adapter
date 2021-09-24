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
	"runs_adapter/middleware"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	pemPath, keyPath, proto, listenport, host, dbuser, dbpassword, dbhost, dbname, logLevel string
	dbport                                                                                  int
	wait, writeTimeout, readTimeout, idleTimeout                                            time.Duration
)

func main() {
	flag.StringVar(&pemPath, "pempath", os.Getenv("SERVERREM"), "path to pem file")
	flag.StringVar(&keyPath, "keypath", os.Getenv("SERVERKEY"), "path to key file")
	flag.StringVar(&listenport, "port", "10000", "port to Listen")
	flag.StringVar(&proto, "proto", "http", "http or https")
	flag.StringVar(&dbuser, "user", "postgres", "db user")
	flag.StringVar(&dbpassword, "password", `postgres`, "db user password")
	flag.StringVar(&dbhost, "host", "localhost", "db host")
	flag.StringVar(&logLevel, "loglevel", "INFO", "log level, default INFO")
	flag.IntVar(&dbport, "dbport", 5432, "db port")
	flag.StringVar(&dbname, "dbname", "runs", "db name")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully")
	flag.DurationVar(&readTimeout, "read-timeout", time.Second*15, "read server timeout")
	flag.DurationVar(&writeTimeout, "write-timeout", time.Second*15, "write server timeout")
	flag.DurationVar(&idleTimeout, "idle-timeout", time.Second*60, "idle server timeout")
	flag.Parse()
	log.Info("Parse flag completed")
	setLogLevel(logLevel)
	log.Info("Set log level completed")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/runs/new", middleware.Middleware(handlers.NewHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/runs/get/lasts", middleware.Middleware(handlers.LastsHandler)).Methods(http.MethodGet, http.MethodOptions).Queries("count", "{count}")
	r.HandleFunc("/api/v1/runs/get/range", middleware.Middleware(handlers.RangeHandler)).Methods(http.MethodGet, http.MethodOptions).Queries("starttime", "{starttime}", "endtime", "{endtime}")
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	http.Handle("/", r)
	log.Info("Register handlers and route completed")
	go func() {
		for {
			err := adapter.InitDB(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname))
			if err != nil {
				log.Error(err)
			}
			time.Sleep(10 * time.Second)
		}
	}()
	//time.Sleep(2 * time.Second)
	// scenarios := []adapter.Scenario{adapter.Scenario{
	// 	Name:     "UC01",
	// 	TPS:      1000,
	// 	SLA:      5,
	// 	Duration: 120},
	// }
	// run := adapter.Run{
	// 	ID:        3,
	// 	StartTime: time.Now(),
	// 	EndTime:   time.Now().Add(1 * time.Hour),
	// 	Data: adapter.Data{
	// 		Project:     "TestProject",
	// 		Grafanalink: "http://master:3000/d/In_The_Bus/jmeter-test-overview?orgId=1&refresh=10s&from=now-1h&to=now",
	// 		Description: "test description",
	// 		Status:      "state ok",
	// 		Scenarios:   scenarios,
	// 	},
	// }
	//log.Println(run)
	//err := run.Insert()
	//log.Println(err)
	// res, err := adapter.GetAllRuns()
	// if err != nil {
	// 	log.Error("Get interface adress error: ", err.Error())
	// }

	// log.Println(res)
	// log.Println(res[0].Data.Scenarios[0].Name)
	// template, err := adapter.GetTableHTML(res)
	// if err != nil {
	// 	log.Error("GetTemplate error: ", err.Error())
	// }
	// log.Println(template)
	// client := confluence.Client{Client: http.DefaultClient, URL: "https://100.65.36.14/wiki/rest/api/content/42068579", AuthBase64: "UkVNYXRza3VzOlNiZXJNb2FuNDY5OCEh"}
	// client.Push([]byte(template))
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error("Get interface adress error: ", err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				host = ipnet.IP.String()
			}
		}
	}
	log.Info("Get IPv4 addr completed")
	srv := &http.Server{
		Addr:         host + ":" + listenport,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      r,
	}
	log.Info("Set server params completed")
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
	srv.Shutdown(ctx)
	log.Info("server shutting down")
	os.Exit(0)
}

func setLogLevel(level string) {
	level = strings.ToUpper(level)
	switch level {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	}
}
