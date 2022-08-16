package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

type Resource struct {
	Id        string  `json:"id,omitempty" bson:"_id,omitempty"`
	Category  string  `json:"category"`
	Org       string  `json:"org"`
	Postal    int     `json:"postal"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	err = client.Connect(ctx)
	db := client.Database(os.Getenv("DB_NAME"))

	if err != nil {
		defer cancel()
		return err
	}
	defer cancel()

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func main() {
	port := ":" + os.Getenv("PORT")
	if port == "" {
		port = ":5000"

	}

	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	godotenv.Load()

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		query := bson.D{{}}

		cursor, err := mg.Db.Collection("resources").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var resources []Resource = make([]Resource, 0)

		if err := cursor.All(c.Context(), &resources); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(resources)
	})
	log.Fatal(app.Listen(port))
}
