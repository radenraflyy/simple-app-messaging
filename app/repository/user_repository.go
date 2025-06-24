package repository

import (
	"SimpleMessaging/app/models"
	"SimpleMessaging/pkg/database"
	"context"

	"go.elastic.co/apm"
)

func InsertNewUser(ctx context.Context, user *models.User) error {
	span, _ := apm.StartSpan(ctx, "InsertNewUser", "repository")
	defer span.End()

	return database.DB.Create(user).Error
}

func InsertNewUserSession(ctx context.Context, session *models.UserSession) error {
	return database.DB.Create(session).Error
}

func GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error) {
	var (
		resp models.UserSession
		err  error
	)

	err = database.DB.Where("token = ?", token).Last(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, err
}

func DeleteUserSessionByToken(ctx context.Context, token string) error {
	return database.DB.Exec("DELETE FROM user_session WHERE token = ?", token).Error
}

func UpdateRefreshTokenSession(ctx context.Context, token, refreshToken string) error {
	return database.DB.Exec("UPDATE user_session SET token = ? WHERE refresh_token = ?", token, refreshToken).Error
}

func GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var (
		resp models.User
		err  error
	)
	err = database.DB.Where("username = ?", username).Last(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, err
}
