package main

import (
	"fmt"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"log"
	"os"
	"time"
)

type Config struct {
	Port string `env:"FIBER_HTTP_PORT"`
	Root string `env:"FIBER_HTTP_ROOT"`
}

var Version string = "0.3.1"

var err error

var Cfg Config

func main() {

	loader := aconfig.LoaderFor(&Cfg, aconfig.Config{
		SkipFlags:          false,
		AllowUnknownEnvs:   true,
		AllowUnknownFields: true,
		SkipEnv:            false,
		//	EnvPrefix:          "FIBERHTTP",
		FileDecoders: map[string]aconfig.FileDecoder{
			".env": aconfigdotenv.New(),
		},
		Files: []string{".env"},
	})

	if err = loader.Load(); err != nil {
		return
	}
	if Cfg.Port == "" {
		Cfg.Port = "8080"
	}
	if Cfg.Root == "" {
		Cfg.Root = "./www"
	}
	fmt.Fprintf(os.Stderr, "lvl=info, message=\"Starting Fiber HTTP v%s\", root=\"%s\", port=%s\n", Version, Cfg.Root, Cfg.Port)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               "http-static",
		EnablePrintRoutes:     false,
	})
	app.Use(logger.New(logger.Config{
		Format: "lvl=info, status=${status}, latency=${latency}, ip=${ip} method=${method}, path=${path}\n",
	}))
	app.Use(cache.New(cache.Config{
		Expiration: 10 * time.Minute,
	}))
	app.Static("/", Cfg.Root)
	log.Fatal(app.Listen(":" + Cfg.Port))
}
