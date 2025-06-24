package controllers

import (
	"SimpleMessaging/app/models"
	"SimpleMessaging/app/repository"
	jwttoken "SimpleMessaging/pkg/jwt_token"
	"SimpleMessaging/pkg/response"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()
	userBody := new(models.User)
	err := ctx.BodyParser(userBody)
	if err != nil {
		errResp := fmt.Errorf("error parsing body: %v", err)
		log.Println("error parsing body: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResp.Error(), nil)
	}

	err = userBody.Validate()
	if err != nil {
		errResp := fmt.Errorf("error validating body: %v", err)
		log.Println("error validating body: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResp.Error(), nil)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if err != nil {
		errResp := fmt.Errorf("error hashing password: %v", err)
		log.Println("error hashing password: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResp.Error(), nil)
	}
	userBody.Password = string(hashPassword)
	err = repository.InsertNewUser(spanCtx, userBody)
	if err != nil {
		errResp := fmt.Errorf("error inserting user: %v", err)
		log.Println("error inserting user: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResp.Error(), nil)
	}

	userBody.Password = ""
	return response.SendSuccessResponse(ctx, userBody)
}

func Login(ctx *fiber.Ctx) error {
	loginReq := &models.LoginRequest{}
	resp := models.LoginResponse{}

	err := ctx.BodyParser(loginReq)
	if err != nil {
		errResp := fmt.Errorf("failed to parser request: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResp.Error(), nil)
	}

	err = loginReq.Validate()
	if err != nil {
		errResp := fmt.Errorf("failed to validate request: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResp.Error(), nil)
	}

	user, err := repository.GetUserByUsername(ctx.Context(), loginReq.Username)
	if err != nil {
		errResp := fmt.Errorf("failed to get username: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, errResp.Error(), nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		errResp := fmt.Errorf("failed to check password: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username/password salah", nil)
	}

	token, err := jwttoken.GenerateToken(ctx.Context(), user.Username, user.Fullname, "token")
	if err != nil {
		errResp := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	refresh_token, err := jwttoken.GenerateToken(ctx.Context(), user.Username, user.Fullname, "token")
	if err != nil {
		errResp := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	userSession := &models.UserSession{
		UserId:       user.ID,
		Token:        token,
		RefreshToken: refresh_token,
	}

	err = repository.InsertNewUserSession(ctx.Context(), userSession)
	if err != nil {
		errResp := fmt.Errorf("failed to insert user session: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	resp.Username = user.Username
	resp.Fullname = user.Fullname
	resp.Token = token
	resp.RefreshTokenToken = refresh_token

	return response.SendSuccessResponse(ctx, resp)
}

func Logout(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	err := repository.DeleteUserSessionByToken(ctx.Context(), token)
	if err != nil {
		erResp := fmt.Errorf("failed delete user session: %v", err)
		log.Println(erResp)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, erResp.Error(), nil)
	}
	return response.SendSuccessResponse(ctx, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Get("Authorization")
	username := ctx.Get("username")
	fullname := ctx.Get("fullname")

	token, err := jwttoken.GenerateToken(ctx.Context(), username, fullname, "token")
	if err != nil {
		errResp := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	err = repository.UpdateRefreshTokenSession(ctx.Context(), token, refreshToken)
	if err != nil {
		errResp := fmt.Errorf("failed to update token: %v", err)
		log.Println(errResp)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
