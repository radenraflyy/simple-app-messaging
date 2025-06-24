package database

import (
	"SimpleMessaging/app/models"
	"SimpleMessaging/pkg/env"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetUpDatabase() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.GetEnv("DB_USER", "root"),
		env.GetEnv("DB_PASSWORD", ""),
		env.GetEnv("DB_HOST", "localhost"),
		env.GetEnv("DB_PORT", "3306"),
		env.GetEnv("DB_NAME", "simple_messaging"),
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database %v", err)
		os.Exit(1)
	}
	DB.Logger = logger.Default.LogMode(logger.Info)

	if err = DB.AutoMigrate(&models.User{}, &models.UserSession{}); err != nil {
		log.Fatalf("Failed to migrate the database %v", err)
	}

	fmt.Printf("Connected to the database %s, host: %s, port: %s\n", env.GetEnv("DB_NAME", "simple_messaging"), env.GetEnv("DB_HOST", "localhost"), env.GetEnv("DB_PORT", "3306"))
	log.Println("====Database Migrated====")
}

func SetupMongoDb() {
	client, err := mongo.Connect(options.Client().
		ApplyURI(env.GetEnv("MONGODB_URI", "")))
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	coll := client.Database("message").Collection("message_history")
	MongoDB = coll

	log.Println("successfully connected to mongoDB")
}
