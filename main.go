package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/matscus/technical_interview/handlers"

	_ "github.com/lib/pq"
)

var (
	pemPath, keyPath, proto, listenport, host, dbuser, dbpassword, dbhost, dbname string
	wait, writeTimeout, readTimeout, idleTimeout                                  time.Duration
	dbport                                                                        int
	DB                                                                            *sql.DB
)

func main() {
	flag.StringVar(&pemPath, "pempath", "/application/server.pem", "path to pem file")
	flag.StringVar(&keyPath, "keypath", "/application/server.key", "path to key file")
	flag.StringVar(&listenport, "port", "9444", "port to Listen")
	flag.StringVar(&proto, "proto", "https", "http or https")
	flag.StringVar(&dbuser, "user", "postgres", "db user")
	flag.StringVar(&dbpassword, "password", `postgres`, "db user password")
	flag.StringVar(&dbhost, "host", "localhost", "db host")
	flag.IntVar(&dbport, "dbport", 5432, "db port")
	flag.StringVar(&dbname, "dbname", "postgres", "db name")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully")
	flag.DurationVar(&readTimeout, "read-timeout", time.Second*15, "read server timeout")
	flag.DurationVar(&writeTimeout, "write-timeout", time.Second*15, "write server timeout")
	flag.DurationVar(&idleTimeout, "idle-timeout", time.Second*60, "idle server timeout")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/book/getbooks", handlers.GetBooks).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	http.Handle("/", r)

	databaseConnection(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname))
	err := createScheme()
	if err != nil {
		log.Println("createScheme error ", err)
		os.Exit(1)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("Get interface adres error: ", err.Error())
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				host = ipnet.IP.String()
			}
		}
	}
	srv := &http.Server{
		Addr:         host + ":" + listenport,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      r,
	}

	go func() {
		switch proto {
		case "https":
			log.Printf("Server is run, proto: https, address: %s ", srv.Addr)
			if err := srv.ListenAndServeTLS(pemPath, keyPath); err != nil {
				log.Println(err)
			}
		case "http":
			log.Printf("Server is run, proto: http, address: %s ", srv.Addr)
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
	log.Println("server shutting down")
	os.Exit(0)
}

func createScheme() (err error) {
	_, err = DB.Query("select * from tUser limit 1")
	log.Println("teststs", err)
	return err
}

func databaseConnection(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Database connection error %s\n", err)
	}
	log.Printf("Database connection")
	go func() {
		for {
			err := DB.Ping()
			if err != nil {
				log.Printf("database ping error %s\n", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
