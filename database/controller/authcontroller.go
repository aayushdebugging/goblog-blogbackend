package controller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aayushdebugging/blogbackend/database"
	"github.com/aayushdebugging/blogbackend/database/models"
	"github.com/dgrijalva/jwt-go"

	"github.com/aayushdebugging/blogbackend/util"

	"github.com/gofiber/fiber/v2"
)
 

func validateEmail(email string)bool{
	Re:=regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z0-9._%+\-]`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err:=c.BodyParser(&data);err!=nil{
		fmt.Println("Unable to prase body")
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}	
	//check if password if lesser that 6 characters
	if len(data["password"].(string))<=6{
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"Password must be greater than 6 character",
		})
	}

	if !validateEmail(strings.TrimSpace(data["email"].(string))){
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"Invalid Email Address",
		})

	}
	//check if email already exist in database
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id!=0{
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":" Email Address already exist",
		})


	}
	user :=models.User{
		FirstName: data["first_name"].(string),
		LastName: data["last_name"].(string),  
		Phone: data["phone"].(string),
		Email: strings.TrimSpace(data["email"].(string)),
	}
	user.SetPassword(data["password"].(string))
	err:=database.DB.Create(&user)
	if err!=nil{
		log.Println(err)
	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":user,
		"message":"Account created Sucessfully",
	})

}

func Login(c *fiber.Ctx)error {
	var data map[string] string
	
	if err:=c.BodyParser(&data);err!=nil{
		fmt.Println("Unable to prase body")
	}
	var user models.User
	database.DB.Where("email=?",data["email"]).First(&user)
	if user.Id == 0{
		c.Status(404)
		return c.JSON(fiber.Map{
			"message":"Email Address doesn't exist ,kindly create an account",
		})
	}
	if err := user.ComparePassword(data["password"]); err!= nil{
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"Incorrect Password",
		})
	}
	token,err:=util.GenerateJwt(strconv.Itoa(int(user.Id)),)
	if err!= nil{
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:        "jwt",
		Value:       token,
		Path:        "",
		Domain:      "",
		MaxAge:      0,
		Expires:     time.Now().Add(time.Hour * 24),
		Secure:      false,
		HTTPOnly:    true,
		SameSite:    "",
		SessionOnly: false,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"you have sucessfully login",
		"user": user,

	})
}



type claims struct{
	jwt.StandardClaims
}