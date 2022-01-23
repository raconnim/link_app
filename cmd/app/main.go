package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"link_app/pkg/handler"
	"link_app/pkg/item"
	"link_app/pkg/middleware"
	"log"
	"net/http"
	"os"
)

func main() {

	dbflag := "" // "postgres" or "memory"
	flag.StringVar(&dbflag, "flag", "memory", "The name of database")
	flag.Parse()
	log.Println(")))))))))", dbflag)
	if dbflag != "postgres" && dbflag != "memory" {
		log.Println("неверный флаг")
	}
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Println("logger error")
		return
	}

	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	var repo handler.ItemRepositoryInterface
	if dbflag == "memory" {
		repo = item.NewItemRepositoryInMemory()
	} else if dbflag == "postgres" {
		name := os.Getenv("PG_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("PG_DB")

		connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
			name, password, dbname, "disable", "db", "5432")
		fmt.Println(connStr)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			fmt.Println("no open bd:", err)
			return
		}

		repo = item.NewItemRepository(db)
	}

	handler := handler.ItemHandler{
		Logger:   logger,
		ItemRepo: repo,
	}
	r := mux.NewRouter()
	r.HandleFunc("/add/", handler.CreateShortLink)
	r.HandleFunc("/link/{SHORT_LINK}", handler.GetLongLink)

	mux := middleware.AccessLog(logger, r)

	addr := ":8000"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}
