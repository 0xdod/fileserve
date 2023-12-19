package main

import (
	"fmt"
	"log"

	"github.com/0xdod/fileserve/filestorage"
	"github.com/0xdod/fileserve/http"
	"github.com/0xdod/fileserve/sqlite"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	viper.AutomaticEnv()
	viper.SetDefault("PORT", 8000)
	viper.SetDefault("DB_NAME", "sqlite3.db")

	var (
		dbName   = viper.GetString("DB_NAME")
		httpPort = viper.GetInt("PORT")
	)

	db := sqlite.DB{
		DSN: dbName,
	}

	if err := db.Open(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("Connected to database")

	port := fmt.Sprintf(":%d", httpPort)
	server := http.NewServer(http.NewServerOpts{
		DB:   &db,
		Addr: &port,
		FileStorage: filestorage.NewS3StorageBackend(&filestorage.S3BackendConfig{
			AccessKeyID:     viper.GetString("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: viper.GetString("AWS_SECRET_ACCESS_KEY"),
			Region:          viper.GetString("AWS_REGION"),
			BucketName:      viper.GetString("AWS_BUCKET_NAME"),
		}),
		FileService: sqlite.NewFileService(&db),
	})

	fmt.Printf("Server running on %s\n", port)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
