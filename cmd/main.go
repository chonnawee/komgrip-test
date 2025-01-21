package main

import (
	"context"
	"encoding/json"
	"fmt"
	"komgrip-test/adapters"
	"komgrip-test/entities"
	"komgrip-test/usecases"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var port string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("cannot read environment variables")
	}
}

func initializeConfig() {
	port = os.Getenv("APP_PORT")
}

func connectMariaDB() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db connected")
	return db
}

func connectMongo() *mongo.Client {
	uri := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?authSource=admin&authMechanism=SCRAM-SHA-256",
		os.Getenv("MONGO_DB_USERNAME"),
		os.Getenv("MONGO_DB_PASSWORD"),
		os.Getenv("MONGO_DB_HOST"),
		os.Getenv("MONGO_DB_PORT"),
		os.Getenv("MONGO_DB_NAME"),
	)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB successfully!")
	return client
}

func main() {
	initializeConfig()
	db := connectMariaDB()
	db.AutoMigrate(entities.Beers{})
	mongoClient := connectMongo()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
		}
	}()
	mongoDB := mongoClient.Database("example")

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	uow := adapters.NewUnitOfWorkDB(db)
	beerLogRepo := adapters.NewBeerLogsRepositoryDB(mongoDB.Collection("beer_logs"))
	usecase := usecases.NewBeersService(uow, beerLogRepo)
	handler := adapters.NewBeersHandler(usecase)

	app.Post("/beer", handler.CreateBeer)
	app.Get("/beer", handler.GetBeers)
	app.Put("/beer/:id", handler.UpdateBeer)
	app.Delete("/beer/:id", handler.DeleteBeer)
	app.Listen(":" + port)
}
