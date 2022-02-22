package main

import (
	"time"

	"github.com/ncostamagna/axul_contact/pkg/client"

	"github.com/joho/godotenv"

	"context"
	"flag"
	"fmt"

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/go-kit/kit/log/level"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ncostamagna/axul_contact/contacts"
	"github.com/ncostamagna/streetflow/slack"
	"github.com/ncostamagna/streetflow/telegram"
)

func main() {

	fmt.Println("Initial")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "postapp",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	_ = level.Info(logger).Log("msg", "service started")
	defer func() {
		_ = level.Info(logger).Log("msg", "service ended")
	}()

	err := godotenv.Load()
	if err != nil {
		_ = level.Info(logger).Log("Error loading .env file", err)
		//os.Exit(-1)
	}

	var httpAddr = flag.String("http", ":"+os.Getenv("APP_PORT"), "http listen address")

	mysqlInfo := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
	time.Sleep(8 * time.Second)
	db, err := gorm.Open("mysql", mysqlInfo)
	if err != nil {
		_ = level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	db.AutoMigrate(contacts.Contact{})

	flag.Parse()
	ctx := context.Background()

	var srv contacts.Service
	{
		tempTran := contacts.NewClient(os.Getenv("USER_GRPC_URL"), "", contacts.GRPC)
		slackTran, _ := slack.NewSlackBuilder(os.Getenv("SLACK_CHANNEL"), os.Getenv("SLACK_TOKEN")).Build()
		telegTran := telegram.NewClient("1536608370:AAErsMmopurv4JhVp1ondOuld8GRUJxohOY", telegram.HTTP)
		userTran := client.NewClient(os.Getenv("USER_GRPC_URL"), "", client.GRPC)
		repository := contacts.NewRepo(db, logger)
		srv = contacts.NewService(repository, slackTran, &telegTran, tempTran, userTran, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	mux := http.NewServeMux()

	mux.Handle("/", contacts.NewHTTPServer(ctx, contacts.MakeEndpoints(srv)))

	http.Handle("/", cors.AllowAll().Handler(accessControl(mux)))

	http.Handle("/metrics", promhttp.Handler())

	/* fmt.Println()
	fmt.Println("Same Value")

	start1 := time.Now()

	recur := 5000
	for i := 1; i < recur; i++ {
		contacts.GetTemplate(1)
	}

	elapsed1 := time.Since(start1)

	fmt.Printf("gRPC took %s", elapsed1)

	start2 := time.Now()

	fmt.Println()
	for i := 1; i < recur; i++ {
		contacts.GetTemplateHTTP(1)
	}

	elapsed2 := time.Since(start2)

	fmt.Printf("HTTP took %s", elapsed2)

	fmt.Println()
	fmt.Println()
	fmt.Println("Others Values")
	start1 = time.Now()

	for i := 1; i < recur; i++ {
		contacts.GetTemplate(uint(i))
	}

	elapsed1 = time.Since(start1)

	fmt.Printf("gRPC took %s", elapsed1)

	start2 = time.Now()

	fmt.Println()


	elapsed2 = time.Since(start2)

	fmt.Printf("HTTP took %s", elapsed2)
	fmt.Println()
	fmt.Println() */
	go func() {
		fmt.Println("listening on port", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)

	}()

	err = <-errs

	if err != nil {
		_ = level.Error(logger).Log("exit", err)
	}

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
