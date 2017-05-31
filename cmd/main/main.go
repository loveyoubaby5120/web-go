package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"jdy/pkg/backend/httpserv"
	"jdy/pkg/site/store"

	"github.com/gorilla/mux"
)

var (
	addr           = flag.String("addr", ":8000", "Binding address of the proxy server.")
	ataAddr        = flag.String("ata_addr", "127.0.0.1:8899", "Binding address of the proxy server.")
	sqlParams      = flag.String("sql_params", "user=vagrant password=vagrant dbname=demo host=127.0.0.1 port=5432 sslmode=disable", "Parameters of SQL database.")
	dir            = flag.String("dir", ".", "the directory to serve files from. Defaults to the current dir")
	insecureCookie = flag.Bool("insecurecookie", true, "Allow insecure cookie (non-https).")
	sentryDSN      = flag.String("sentry_dsn", "", "Sentry ata project dsn.")
)

var db *sql.DB

func main() {
	flag.Parse()

	// pgEnv := pgdb.NewEnv(*sqlParams)

	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	httpserv.AddRoutes(r, store.Routes())

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*dir))))

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
		// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		// p, err := route.GetPathRegexp()
		// if err != nil {
		// 	return err
		// }
		// m, err := route.GetMethods()
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(strings.Join(m, ","), t, p)

		fmt.Println("tpl: ", t)
		return nil
	})
	http.Handle("/", r)

	log.Fatal(srv.ListenAndServe())
}
