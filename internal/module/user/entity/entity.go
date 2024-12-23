package entity

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	HassedPassword string
}

type RegisterByGoogleRequest struct {
	Email          	string `json:"email" validate:"required,email"`
	Name           	string `json:"name" validate:"required"`
	GoogleId		string `json:"id"`	
	ImageUrl		string `json:"picture"`	
	Password       	string `json:"password,omitempty"`
	HassedPassword 	string
}

type RegisterResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Remember bool	`json:"remember"`
}

type LoginResponse struct {
	AccessToken		string	`json:"accessToken"`
	RefreshToken	string	`json:"refreshToken"`
}

type ProfileRequest struct {
	UserId string `validate:"required"`
}

type ProfileResponse struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Role  string `json:"-" db:"role"`
}

type UserPayload struct {
	UserID       string `json:"user_id" db:"id"`
	Role         string `json:"role" db:"role"`
}
