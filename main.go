package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ServerGroup struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Servers []Server `json:"servers"`
}

type Server struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Response struct {
	Data string `json:"data"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&Response{Data: "Monit v2 - online"})
	})

	app.Get("/servers", func(c *fiber.Ctx) error {
		file, err := ioutil.ReadFile("./server.json")
		if err != nil {
			panic(err)
		}
		groups := []ServerGroup{}

		err = json.Unmarshal([]byte(file), &groups)
		if err != nil {
			panic(err)
		}

		return c.Status(200).JSON(groups)
	})

	app.Get("/ping/:ip", func(c *fiber.Ctx) error {
		ip := c.Params("ip")
		err := ping(ip)
		if err != nil {
			return c.Status(400).JSON(&Response{Data: ip + " - offline"})
		} else {
			return c.Status(200).JSON(&Response{Data: ip + " - online"})
		}
	})

	app.Listen(":12900")
}

func ping(host string) error {
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return errors.New("connection timeout")
	} else {
		return nil
	}
}
