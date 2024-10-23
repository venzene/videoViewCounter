package main

import (
	"bufio"
	"context"
	"flag"
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
	httpclient "view_count/viewservice/httpclient"

	"github.com/go-kit/kit/metrics/prometheus"
	kitlog "github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	_ "github.com/lib/pq"
)

func main() {
	mode := flag.String("mode", "server", "Use server, client, cli")
	flag.Parse()
	if *mode == "server" {
		server()
	} else if *mode == "client" {
		client()
	} else {
		log.Fatalf("Unknown mode: %s. Use 'server' or 'client'.", *mode)
	}
}

func client() {

	client := httpclient.NewClient("http://localhost:8080")
	fmt.Println("Client started...")

	for {
		var option int
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Get All Views")
		fmt.Println("2. Get Top Videos")
		fmt.Println("3. Increment Video View")
		fmt.Println("4. Get View by ID")
		fmt.Println("5. Get Recent Incremented Videos")
		fmt.Println("0. Exit")
		fmt.Print("Enter option: ")
		fmt.Scan(&option)

		switch option {
		case 1:
			res, err := client.GetAllViews(context.Background(), nil)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("GetAllViews Response: %+v\n", res)
			}

		case 2:
			var n int
			fmt.Print("Enter number of top videos to fetch: ")
			fmt.Scan(&n)
			res, err := client.GetTopVideos(context.Background(), n)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Top Videos: %+v\n", res)
			}

		case 3:
			var vid string
			fmt.Print("Enter Video ID to increment: ")
			fmt.Scan(&vid)
			_, err := client.Increment(context.Background(), vid)

			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Incremented view count successfully")
			}

		case 4:
			var vid string
			fmt.Print("Enter Video ID: ")
			fmt.Scan(&vid)
			res, err := client.GetView(context.Background(), vid)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Video view count: %v\n", res)
			}

		case 5:
			var n int
			fmt.Print("Enter number of recent incremented videos: ")
			fmt.Scan(&n)
			res, err := client.GetRecentVideos(context.Background(), n)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Recent incremented videos: %+v\n", res)
			}

		case 0:
			fmt.Println("Exiting client...")
			return

		default:
			fmt.Println("Invalid option, try again.")
		}
	}
}

func server() {
	var vs viewservice.Service // concrete vs interface type declaration

	// viewRepo := viewrepository.NewInmemoryRepo()

	database, err := database.Connect("view_count")
	if err != nil {
		log.Fatal("Error in database connection: ", err)
	}

	viewRepo := viewrepository.NewPostgresRepo(database)

	vs = viewservice.NewService(viewRepo)
	cliVs := vs

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

	endpoints := viewservice.MakeEndpoints(vs)

	r := viewservice.MakeHandler(endpoints, logger)

	// h := NewHandler(vs)

	// r := routeIntialiser(*h)

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
		if database != nil {
			log.Println("Closing Database connection...")
			if err := database.Close(); err != nil {
				log.Fatalf("Error closing database %s\n", err)
			}
		}
		done <- true
	}()

	// TODO: create command line tool -> inc vid, viewall, view vid : Done

	// TODO: run this with postgress in local machine : Done

	go func() {
		log.Println("Server started on :8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commands (getView, incre, getAll, top, recent, exit):")

	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			break
		} else {
			if err := cli.Execute(cliVs, input); err != nil {
				fmt.Println(err)
			}
		}
	}

	<-done
	fmt.Println("Shutting the system now.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %s", err)
	}
	fmt.Println("Server closed gracefully!!")
}
