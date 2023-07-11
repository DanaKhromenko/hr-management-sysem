package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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
	//TODO: implement function
	return nil
}

func getEmployee(ctx *fiber.Ctx) error {
	//TODO: implement function
	return nil
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
