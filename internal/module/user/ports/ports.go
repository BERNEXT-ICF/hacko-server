package ports

import (
	"context"
	oauthgoogleent "hacko-app/internal/integration/oauth2google/entity"
	"hacko-app/internal/module/user/entity"
)

type UserRepository interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	RegisterByGoogle(ctx context.Context, req *entity.RegisterByGoogleRequest) (*entity.RegisterResponse, error)
	FindByEmail(ctx context.Context, email string) (*entity.UserResult, error)
	FindById(ctx context.Context, id string) (*entity.ProfileResponse, error)
	UpdateRefreshToken(ctx context.Context, userId, refreshToken string) error
	FindRefreshToken(ctx context.Context, refreshToken string) (*entity.UserPayload, error) 
}
type UserService interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	Profile(ctx context.Context, req *entity.ProfileRequest) (*entity.ProfileResponse, error)
	GetOauthGoogleUrl(ctx context.Context) (string, error)
	LoginGoogle(ctx context.Context, req *oauthgoogleent.UserInfoResponse) (*entity.LoginResponse, error)
	RefreshTokenService(ctx context.Context, refreshToken string) (string, error)
}
