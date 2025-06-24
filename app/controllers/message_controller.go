package controllers

import (
	"SimpleMessaging/app/repository"
	"SimpleMessaging/pkg/response"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetHistory(ctx *fiber.Ctx) error {
	resp, err := repository.GetAllMessage(ctx.Context())
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada server", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
