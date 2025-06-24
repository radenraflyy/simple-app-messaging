package router

import (
	"SimpleMessaging/app/repository"
	jwttoken "SimpleMessaging/pkg/jwt_token"
	"SimpleMessaging/pkg/response"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth != "" {
		log.Println("authorization empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unathorized", nil)
	}

	_, err := repository.GetUserSessionByToken(ctx.Context(), auth)
	if err != nil {
		log.Println("failed to get userSession in DB")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unathorized", nil)
	}

	claim, err := jwttoken.ValidateToken(ctx.Context(), auth)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unathorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unathorized", nil)
	}
	ctx.Set("username", claim.Username)
	ctx.Set("fullname", claim.Fullname)
	return ctx.Next()
}

func AuthMiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth != "" {
		log.Println("authorization empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unathorized", nil)
	}

	claim, err := jwttoken.ValidateToken(ctx.Context(), auth)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unathorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unathorized", nil)
	}
	ctx.Set("username", claim.Username)
	ctx.Set("fullname", claim.Fullname)
	return ctx.Next()
}
