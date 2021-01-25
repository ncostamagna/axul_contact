package main

import (
	"time"

	"github.com/joho/godotenv"

	"context"
	"flag"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"

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

	b, err11 := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL: "https://149.154.167.40:443",

		Token:  "5101b7cda1e99c3456a56b4753b07afa",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err11 != nil {
		fmt.Println(err11)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	b.Start()

	fmt.Println("Initial")

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

	fmt.Println("Env")
	err := godotenv.Load()
	if err != nil {
		_ = level.Error(logger).Log("Error loading .env file", err)
		os.Exit(-1)
	}

	var httpAddr = flag.String("http", ":"+os.Getenv("APP_PORT"), "http listen address")

	mysqlInfo := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

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
		tempTran := contacts.NewClient(":50055", "", contacts.GRPC)
		slackTran, _ := slack.NewSlackBuilder("birthday", "xoxb-1448869030753-1436532267283-AZoMMLoxODNMC5xydelq1uLP").Build()
		telegTran := telegram.NewClient("1536608370:AAErsMmopurv4JhVp1ondOuld8GRUJxohOY", telegram.HTTP)
		repository := contacts.NewRepo(db, logger)
		srv = contacts.NewService(repository, *slackTran,*telegTran, tempTran, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	mux := http.NewServeMux()

	mux.Handle("/contacts/", contacts.NewHTTPServer(ctx, contacts.MakeEndpoints(srv)))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println()
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
	/* 	for i := 1; i < recur; i++ {
		contacts.GetTemplateHTTP(uint(i))
	} */

	elapsed2 = time.Since(start2)

	fmt.Printf("HTTP took %s", elapsed2)
	fmt.Println()
	fmt.Println()
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
