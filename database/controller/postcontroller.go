package controller

import (
	"errors"
	"math"
	"strconv"

	"github.com/aayushdebugging/blogbackend/database"
	"github.com/aayushdebugging/blogbackend/database/models"
	"github.com/aayushdebugging/blogbackend/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var blogpost models.Blog
	if err := c.BodyParser(&blogpost); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid payload",
		})
	}

	if err := database.DB.Create(&blogpost).Error; err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Congratulations! Your post is live",
	})
}

func AllPost(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 5
	offset := (page - 1) * limit
	var total int64
	var getblog []models.Blog
	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getblog)
	database.DB.Model(&models.Blog{}).Count(&total)
	return c.JSON(fiber.Map{
		"data": getblog,
		"meta": fiber.Map{
			"total":      total,
			"page":       page,
			"last_page":  math.Ceil(float64(int(total))/float64(limit)),
		},
	})
}

func DetailPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}

	var blogpost models.Blog
	if err := database.DB.Where("id=?", id).Preload("User").First(&blogpost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"data": blogpost,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}

	var blog models.Blog
	if err := c.BodyParser(&blog); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid payload",
		})
	}

	database.DB.Model(&models.Blog{}).Where("id=?", id).Updates(blog)

	return c.JSON(fiber.Map{
		"message": "Post updated successfully",
	})
}

func UniquePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, err := util.ParseJwt(cookie)
	if err != nil {
		c.Status(401)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var blog []models.Blog
	database.DB.Model(&models.Blog{}).Where("user_id", id).Preload("User").Find(&blog)
	return c.JSON(blog)
}

func DeletePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}

	var blog models.Blog
	blog.Id = uint(id)

	deleteResult := database.DB.Delete(&blog)
	if deleteResult.Error != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	if deleteResult.RowsAffected == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Post not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}