package main

import (
	"github.com/digitalhouse-dev/dh-kit/logger"
	"github.com/joho/godotenv"
	authentication "github.com/ncostamagna/axul_auth/auth"
	"github.com/ncostamagna/axul_contact/internal/contact"
	"github.com/ncostamagna/axul_contact/pkg/client"
	"github.com/ncostamagna/axul_contact/pkg/handler"
	"github.com/ncostamagna/streetflow/slack"
	"github.com/ncostamagna/streetflow/telegram"
	"time"

	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func main() {

	fmt.Println("Initial")
	var log = logger.New(logger.LogOption{Debug: true})
	_ = godotenv.Load()

	//var httpAddr = flag.String("http", ":"+os.Getenv("APP_PORT"), "http listen address")

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
		err := db.AutoMigrate(contact.Contact{})
		_ = log.CatchError(err)
	}

	//flag.Parse()
	ctx := context.Background()

	token := os.Getenv("TOKEN")
	auth, err := authentication.New(token)
	if err != nil {
		_ = log.CatchError(fmt.Errorf("%s err: %v", dsn, err))
		os.Exit(-1)
	}

	var service contact.Service
	{
		tempTran := contact.NewClient(os.Getenv("USER_GRPC_URL"), "", contact.GRPC)
		slackTran, _ := slack.NewSlackBuilder(os.Getenv("SLACK_CHANNEL"), os.Getenv("SLACK_TOKEN")).Build()
		telegTran := telegram.NewClient("1536608370:AAErsMmopurv4JhVp1ondOuld8GRUJxohOY", telegram.HTTP)
		userTran := client.NewClient(os.Getenv("USER_GRPC_URL"), "", client.GRPC)
		repository := contact.NewRepo(db, log)
		service = contact.NewService(repository, slackTran, &telegTran, tempTran, userTran, auth, log)
	}

	/*go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()*/

	h := handler.NewHTTPServer(ctx, contact.MakeEndpoints(service))
	port := os.Getenv("APP_PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)
	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  4 * time.Second,
	}

	//http.Handle("/", cors.AllowAll().Handler(accessControl(mux)))
	errSrv := make(chan error)

	go func() {
		fmt.Println("listening on port", address)
		errSrv <- srv.ListenAndServe()

	}()

	err = <-errSrv
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
