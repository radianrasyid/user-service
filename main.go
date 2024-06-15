package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prithuadhikary/user-service/controller"
	"github.com/prithuadhikary/user-service/domain"
	"github.com/prithuadhikary/user-service/repository"
	"github.com/prithuadhikary/user-service/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func main() {
	db, err := InitialiseDB(&DbConfig{
		User:     "postgres",
		Password: "200875",
		DbName:   "go-learn",
		Host:     "localhost",
		Port:     "5432",
		Schema:   "public",
	})
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	err = db.AutoMigrate(&domain.User{}, &domain.Session{})
	if err != nil {
		log.Fatalf("failed to perform database operation: %v", err)
	}
	DB = db
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	engine := gin.Default()

	controller.NewUserController(engine, userService)

	server := &http.Server{
		Addr:    ":8088",
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to gracefully shutdown server: %v", err)
	}

	log.Println("Server gracefully stopped")

}

func InitialiseDB(dbConfig *DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v TimeZone=Asia/Kolkata", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DbName, dbConfig.Port)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dbConfig.Schema + ".",
			SingularTable: false,
		},
	})
	if err != nil {
		return nil, err
	}
	return db, err
}

type DbConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbName"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Schema   string `mapstructure:"schema"`
}
