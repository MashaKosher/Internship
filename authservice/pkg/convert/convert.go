package convert

import "authservice/internal/entity"

func ConvertUserToResponse(accessToken string, refreshToken string, user *entity.User) entity.UserResponse {
	return entity.UserResponse{
		Role:         user.Role,
		ID:           user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}
}
