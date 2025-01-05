package models

type User struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Age        *int    `json:"age"`
	Level      *int    `json:"level"`
	Gender     *string `json:"gender"`
	ClassID    *int    `json:"class_id"`
	AvatarPath *string `json:"avatar_path"`
	ModelsPath *string `json:"models_path"`
	Message    *string `json:"message"`
	RoleID     *int    `json:"role_id"`
}

// LoginResponse defines the successful response for the login request
type LoginResponse struct {
	User         string `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Email        string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type RefreshToken struct {
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}
