package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/catalogue"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

func main() {
	var (
		//	port   = flag.String("port", "8081", "Port to bind HTTP listener") // TODO(pb): should be -addr, default ":8081"
		//	images = flag.String("images", "./images/", "Image path")
		dsn = flag.String("DSN", "catalogue_user:default_password@tcp(172.17.0.2:3306)/cataloguedb", "Data Source Name: [username[:password]@][protocol[(address)]]/dbname")
	//	zip    = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	db, err := sqlx.Open("mysql", *dsn)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Check if DB connection can be made, only for logging purposes, should not fail/exit
	err = db.Ping()
	if err != nil {
		logger.Log("Error", "Unable to connect to Database", "DSN", dsn)
	}

	var svc catalogue.Service
	svc = catalogue.NewCatalogueService(db, logger)
	endpoint := catalogue.Endpoints{
		GetListEndpoint:       catalogue.MakeGetListEndpoint(svc),
		GetProductEndpoint:    catalogue.MakeGetProductEndpoint(svc),
		GetCategoriesEndpoint: catalogue.MakeGetCategoriesEndpoint(svc),
	}

	r := catalogue.MakeHttpHandler(ctx, endpoint, logger)

	// HTTP transport
	go func() {
		fmt.Println("Starting server at port 8080")
		handler := r
		errChan <- http.ListenAndServe(":8080", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}
