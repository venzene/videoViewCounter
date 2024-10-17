package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"view_count/cli"
	"view_count/database.go"
	"view_count/repository/viewrepository"
	"view_count/viewservice"

	"github.com/go-kit/kit/metrics/prometheus"
	kitlog "github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	_ "github.com/lib/pq"
)

func main() {

	var vs viewservice.Service // concrete vs interface type declaration

	// viewRepo := viewrepository.NewInmemoryRepo()

	// TODO: add logging mw @abhishekAK/@abhishekGupta/Rishu
	// TODO: add instumenting mw
	database.CreateDB("view_count")
	database, err := database.Connect("view_count")
	if err != nil {
		log.Fatal("Error in database connection: ", err)
	}

	viewRepo := viewrepository.NewPostgresRepo(database)

	// TODO: add logging mw : Done
	// TODO: add instumenting mw : Done

	vs = viewservice.NewService(viewRepo)

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	vs = viewservice.NewServiceLogging(logger, vs)

	requestCount := prometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "video_service",
		Subsystem: "view_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{"method"})
	requestLatency := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "video_service",
		Subsystem: "view_service",
		Name:      "request_latency_seconds",
		Help:      "Total duration of requests in seconds.",
	}, []string{"method"})

	vs = viewservice.NewInstrumentingService(requestCount, requestLatency, logger, vs)

	// endpoints := viewservice.MakeEndpoints(vs)

	// r := viewservice.MakeHandler(endpoints, logger)

	h := NewHandler(vs)

	r := routeIntialiser(*h)

	// TODO: handle intrupt gracefully

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("Recieved signal %s. Shutdown begins \n", sig)
		// close the db
		// if database != nil {
		// 	log.Println("Closing Database connection...")
		// 	if err := database.Close(); err != nil {
		// 		log.Fatalf("Error closing database %s\n", err)
		// 	}
		// }
		done <- true
	}()

	// TODO: create command line tool -> inc vid, viewall, view vid : Done

	// TODO: run this with postgress in local machine : Done

	// go func() {
	// 	log.Println("Server started on :8080")
	// 	log.Fatal(http.ListenAndServe(":8080", r))
	// }()

	<-done
	fmt.Println("Shutting the system now.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %s", err)
	}
	fmt.Println("Server closed gracefully!!")

	if err := cli.Execute(vs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
