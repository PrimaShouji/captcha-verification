package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/PrimaShouji/captcha-verification/pkg/captchagen"
)

type VerifyResult struct {
	Result bool `json:"result"`
}

func intializeLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s |", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	return zerolog.New(output).With().Timestamp().Logger()
}

func main() {
	// Initialize logger
	log := intializeLogger()

	// Initialize API
	app := fiber.New()

	app.Get("/generate/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		r, err := captchagen.Generate(id)
		if err != nil {
			log.Info().Str("id", id).Err(err).Msg("failed to generate CAPTCHA image")
			return c.SendStatus(500)
		}

		log.Info().Str("id", id).Msg("generated CAPTCHA image")

		c.Context().SetContentType("image/png")
		return c.SendStream(r)
	})

	app.Get("/verify/:id/:test", func(c *fiber.Ctx) error {
		id := c.Params("id")
		test := c.Params("test")

		verified := captchagen.Verify(id, test)

		log.Info().Str("id", id).Str("test", test).Bool("verified", verified).Msg("got verification request for CAPTCHA image")

		return c.JSON(&VerifyResult{
			Result: verified,
		})
	})

	// Listen for requests
	log.Info().Msg("application started")
	app.Listen(":2539")
}
