package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gitlab.com/final_project1240930/booking_service/internal/logs"
	"gitlab.com/final_project1240930/booking_service/internal/repository"
	"gitlab.com/final_project1240930/booking_service/internal/services"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// โหลดค่าการตั้งจากไฟล์ .env
	if err := godotenv.Load("/app/.env"); err != nil {
		logs.Fatal("Error loading .env file", zap.Error(err))
		return
	}
	// if err := godotenv.Load("../.env"); err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		logs.Fatal("Error converting DB_PORT to int", zap.Error(err))
	}

	// Connect to the database
	db, err := repository.NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		logs.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		if err := repository.CloseDatabase(db); err != nil {
			logs.Error("Failed to close database", zap.Error(err))
		}
	}()

	logs.Info("Connected to the database successfully!")

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		logs.Fatal("APP_PORT is not set in .env file")
	}

	// Start the gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", appPort))
	if err != nil {
		logs.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()

	// --------------------------- Menu -------------------------------
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	menuRepositoryDB := repository.NewMenuRepository(db, cloudinaryURL)
	services.RegisterMenuServiceServer(s, services.NewMenuServer(menuRepositoryDB))

	// --------------------------- Table -------------------------------

	tableRepositoryDB := repository.NewTableRepository(db)
	services.RegisterTableServiceServer(s, services.NewTableServer(tableRepositoryDB))

	// -----------------------------------------------------------------

	logs.Info("Server running on port " + appPort)

	if err := s.Serve(listener); err != nil {
		logs.Fatal("Failed to serve", zap.Error(err))
	}

}
