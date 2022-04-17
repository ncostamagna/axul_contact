package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/digitalhouse-dev/dh-kit/logger"
	"github.com/joho/godotenv"
	authentication "github.com/ncostamagna/axul_auth/auth"
	"github.com/ncostamagna/axul_contact/pkg/client"

	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"os"
	"os/signal"
	"syscall"

	"github.com/ncostamagna/axul_contact/contacts"
	"github.com/ncostamagna/streetflow/slack"
	"github.com/ncostamagna/streetflow/telegram"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	fmt.Println("Initial")
	var log = logger.New(logger.LogOption{Debug: true})

	err := godotenv.Load()
	if err != nil {
		_ = log.CatchError(err)
		//os.Exit(-1)
	}

	var httpAddr = flag.String("http", ":5000", "http listen address")

	fmt.Println("DataBases")
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = log.CatchError(err)
		os.Exit(-1)
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		err := db.AutoMigrate(contacts.Contact{})
		_ = log.CatchError(err)
	}

	flag.Parse()
	ctx := context.Background()

	token := os.Getenv("TOKEN")
	auth, err := authentication.New(token)
	if err != nil {
		_ = log.CatchError(fmt.Errorf("%s err: %v", dsn, err))
		os.Exit(-1)
	}

	var srv contacts.Service
	{
		tempTran := contacts.NewClient(os.Getenv("USER_GRPC_URL"), "", contacts.GRPC)
		slackTran, _ := slack.NewSlackBuilder(os.Getenv("SLACK_CHANNEL"), os.Getenv("SLACK_TOKEN")).Build()
		telegTran := telegram.NewClient("1536608370:AAErsMmopurv4JhVp1ondOuld8GRUJxohOY", telegram.HTTP)
		userTran := client.NewClient(os.Getenv("USER_GRPC_URL"), "", client.GRPC)
		repository := contacts.NewRepo(db, log)
		srv = contacts.NewService(repository, slackTran, &telegTran, tempTran, userTran, auth, log)
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

	go func() {
		fmt.Println("listening on port", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)

	}()

	err = <-errs

	if err != nil {
		_ = log.CatchError(err)
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
