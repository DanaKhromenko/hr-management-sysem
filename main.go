package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName                 = "hr-management-system"
	mongoUri               = "mongodb://localhost:27017/" + dbName
	employeesNameInMongoDB = "employees"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type Employee struct {
	// json for frontend or postman, bson for MongoDB
	Id     string  `json: "id,omitempty" bson: "_id,omitempty"`
	Name   string  `json: "name"`
	Salary float64 `json: salary`
	Age    float64 `json: age`
}

var mg MongoInstance

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		return err
	}

	// To avoid blocking the entire program because of the MongoDB blocking functions (like Insert)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	if err := client.Connect(ctx); err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     client.Database(dbName),
	}
	return nil
}

func getEmployee(ctx *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := mg.Db.Collection(employeesNameInMongoDB).Find(ctx.Context(), query)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	var employees []Employee = make([]Employee, 0)
	if err = cursor.All(ctx.Context(), &employees); err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.JSON(employees)
}

func postEmployee(ctx *fiber.Ctx) error {
	//TODO: implement function
	return nil
}

func putEmployee(ctx *fiber.Ctx) error {
	//TODO: implement function
	return nil
}

func deleteEmployee(ctx *fiber.Ctx) error {
	//TODO: implement function
	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", getEmployee)
	app.Post("/employee", postEmployee)
	app.Put("/employee/:id", putEmployee)
	app.Delete("/employee/:id", deleteEmployee)

	log.Fatal(app.Listen("3000"))
}
