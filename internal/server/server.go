package server

import (
	"flag"
	"fmt"
	"log"
	"mini_jira/config"
	"mini_jira/internal/database"
	"mini_jira/migrations"
	"mini_jira/repository"
	"mini_jira/routes"
	"mini_jira/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)

	cfg := config.GetConfig()

	db, sqlDB := database.ConnectDB()

	downFlag := flag.Bool("down", false, "Run database migration down")
	downAllFlag := flag.Bool("down-all", false, "Run all database migrations down")

	flag.Parse()

	if *downFlag {
		log.Println("Running database migration down...")
		migrations.Down(sqlDB)
		log.Println("Successfully run database migration down.")
		return
	}

	if *downAllFlag {
		log.Println("Running all database migrations down...")
		migrations.DownAll(sqlDB)
		log.Println("Successfully run all database migrations down.")
		return
	}

	log.Println("Running database migration up...")

	migrations.Up(sqlDB)

	log.Println("Successfully run database migration up.")

	repo := repository.New(db)

	serv, err := service.New(repo)
	if err != nil {
		log.Fatalf("failed to initialize service: %v", err)
	}

	if cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	gin.DefaultWriter = os.Stdout
	gin.DefaultErrorWriter = os.Stderr

	r := routes.SetupRoutes(serv)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.PORT),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server is running on port %d", cfg.PORT)
	log.Fatal(srv.ListenAndServe())
}
