package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"net/http"
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
)

func main() {
	flag.StringVar(&pemPath, "pempath", os.Getenv("SERVERREM"), "path to pem file")
	flag.StringVar(&keyPath, "keypath", os.Getenv("SERVERKEY"), "path to key file")
	flag.StringVar(&listenport, "port", "9443", "port to Listen")
	flag.StringVar(&proto, "proto", "https", "http or https")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully")
	flag.DurationVar(&readTimeout, "read-timeout", time.Second*15, "read server timeout")
	flag.DurationVar(&writeTimeout, "write-timeout", time.Second*15, "write server timeout")
	flag.DurationVar(&idleTimeout, "idle-timeout", time.Second*60, "idle server timeout")
	flag.StringVar(&dbuser, "user", "postgres", "db user")
	flag.StringVar(&dbpassword, "password", `postgres`, "db user password")
	flag.StringVar(&dbhost, "host", "postgres:5433", "db host")
	flag.StringVar(&dbname, "dbname", "postgres", "db name")
	flag.Parse()
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/book/getbooks", handlers.GetBooks).Methods(http.MethodGet, http.MethodOptions)
	http.Handle("/api/v1/", r)

	//connStr := "postgres://" + dbuser + ":" + dbpassword + "@" + dbhost + "/" + dbname + "?sslmode=disable"
	//databaseConnection(connStr)

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

func databaseConnection(connStr string) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Database connectin error %s\n", err)
	}
	go func() {
		for {
			err := db.Ping()
			if err != nil {
				log.Printf("database ping error %s\n", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

}
