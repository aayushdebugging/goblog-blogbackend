package main

import (
	"log"
	"os"

	"github.com/aayushdebugging/blogbackend/database"
	"github.com/aayushdebugging/blogbackend/database/controller/route"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	database.Connect()
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error loading .env files")
	}
	port:=os.Getenv("PORT")
	app := fiber.New()
	route.Setup(app)
	app.Listen(":"+port)

}