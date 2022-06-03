package routes

import (
	"go-url-shortner/database"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error { //passing context
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error cannot parse JSON"})
	}

	//implement rate limiting

	r2 := database.CreateClient(1)
    defer r2.Close()
	r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil{
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		r2.Get(database.Ctx, c.IP()).Result())
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0{
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Rate limit exceeded",
                "rate_limit_rest": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	//check if the input is an actual url
	
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error Invalid URL"})

	}

	//check the domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "domain not found"})
	}

	//enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)
}