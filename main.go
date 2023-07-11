package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	collection := mg.Db.Collection(employeesNameInMongoDB)

	var employee Employee
	if err := ctx.BodyParser(employee); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	employee.Id = ""
	insertionResult, err := collection.InsertOne(ctx.Context(), employee)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(ctx.Context(), filter)

	createdEmployee := &Employee{}
	createdRecord.Decode(createdEmployee)
	return ctx.Status(http.StatusOK).JSON(createdEmployee)
}

func putEmployee(ctx *fiber.Ctx) error {
	idParam := ctx.Params("_id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	var employee Employee
	if err := ctx.BodyParser(&employee); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	filterEmployeeByIdQuery := bson.D{{Key: "_id", Value: id}}
	updateEmployeeQuery := bson.D{{
		Key: "$set", Value: bson.D{
			{Key: "name", Value: employee.Name},
			{Key: "salary", Value: employee.Salary},
			{Key: "age", Value: employee.Age},
		},
	}}

	if err := mg.Db.Collection(employeesNameInMongoDB).FindOneAndUpdate(ctx.Context(), filterEmployeeByIdQuery, updateEmployeeQuery).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.Status(http.StatusBadRequest).SendString(err.Error())
		}
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	employee.Id = idParam
	return ctx.Status(http.StatusOK).JSON(employee)
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
