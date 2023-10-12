package main

import (
	"math/rand"

	"example.com/fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	go_ssr "github.com/natewong1313/go-react-ssr"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/react"
)

var APP_ENV string

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(favicon.New(favicon.Config{
		File: "../frontend/public/favicon.ico",
		URL:  "/favicon.ico",
	}))
	app.Static("/assets", "../frontend/public/")

	go_ssr.Init(config.Config{
		AppEnv:             APP_ENV,
		AssetRoute:         "/assets",
		FrontendDir:        "../frontend/src",
		GeneratedTypesPath: "../frontend/src/generated.d.ts",
		PropsStructsPath:   "./models/props.go",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		response := react.RenderRoute(react.Config{
			File:  "Home.tsx",
			Title: "Fiber example app",
			MetaTags: map[string]string{
				"og:title":    "Fiber example app",
				"description": "Hello world!",
			},
			Props: &models.IndexRouteProps{
				InitialCount: rand.Intn(100),
			},
		})
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(string(response))
	})

	app.Listen(":8080")
}
