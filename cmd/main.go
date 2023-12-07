package main

import (
	"context"
	"fmt"
	"github.com/digitalhouse-dev/dh-kit/logger"
	"github.com/joho/godotenv"
	authentication "github.com/ncostamagna/axul_auth/auth"
	"github.com/ncostamagna/axul_contact/internal/contact"
	"github.com/ncostamagna/axul_contact/pkg/bootstrap"
	"github.com/ncostamagna/axul_contact/pkg/handler"
	"github.com/starry-axul/notifit-go-sdk/notify"
	"net/http"
	"os"
	"time"
)

func main() {

	fmt.Println("Initial")
	var log = logger.New(logger.LogOption{Debug: true})
	_ = godotenv.Load()

	db, err := bootstrap.DBConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	ctx := context.Background()

	token := os.Getenv("TOKEN")
	auth, err := authentication.New(token)
	if err != nil {
		_ = log.CatchError(fmt.Errorf("err: %v", err))
		os.Exit(-1)
	}

	var service contact.Service
	{
		notifTran := notify.NewHttpClient(os.Getenv("PUSH_URL"), "")
		repository := contact.NewRepo(db, log)
		service = contact.NewService(repository, notifTran, auth, log)
	}


	h := handler.NewHTTPServer(ctx, auth, contact.MakeEndpoints(service))
	url := os.Getenv("APP_URL")
	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         url,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  4 * time.Second,
	}

	errSrv := make(chan error)

	go func() {
		fmt.Println("listening on", url)
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
