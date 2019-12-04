// Account Payment Service.
// main.go main entry point into application.
package main

import (
	"context"
	"flag"
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gopayment/account"
)

// Database ENV variables
const (
	dbhost = "POSTGRES_HOST"
	dbport = "POSTGRES_PORT"
	dbuser = "POSTGRES_USER"
	dbpass = "POSTGRES_PASSWORD"
	dbname = "POSTGRES_DB"
)

// initDb setup database connection and return it.
func initDb() *sql.DB {
	var (
		err error
		db  *sql.DB
	)

	config := dbConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")
	return db
}

// dbConfig parse ENV variables and return config key/value for database.
func dbConfig() map[string]string {
	conf := make(map[string]string)

	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}

	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

// main starts WEB server on port 8888.
// Setups logger, database connection, servises and handlers
// for application.
func main() {
	// base declarations
	var (
		logger    log.Logger
		ctx       context.Context
		db        *sql.DB
		pg        account.Dal
		srv       account.Service
		endpoints account.Endpoints
		handler   http.Handler
		httpAddr  = flag.String("http", ":8888", "http listen address")
	)

	// logger staff
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	// defer level and db connection
	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// create database connection and graceful shutdown for it
	db = initDb()
	defer db.Close()

	// setup service
	ctx = context.Background()
	pg = account.NewDal(db, logger)
	srv = account.NewService(pg, logger)
	errs := make(chan error)

	// POSIX events
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// setup endpoints
	endpoints = account.MakeEndpoints(srv)
	flag.Parse()

	// start application
	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler = account.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
