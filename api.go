// package main

// import (
// 	"fmt"

// 	"github.com/gofiber/fiber/v2"
// )

// type SignInRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// func SignInWeb(c *fiber.Ctx) {
// 	signInRequestObject := SignInRequest{}
// 	if err := c.BodyParser(&signInRequestObject); err != nil {
// 		c.Status(500)
// 		return nil
// 	}
// 	fmt.Println(signInRequestObject)
// 	c.JSON(("its working"))
// 	return nil
// }
