package service

import (
	"context"
	integOauth "hacko-app/internal/integration/oauth2google"
	oauthgoogleent "hacko-app/internal/integration/oauth2google/entity"
	"hacko-app/internal/module/user/entity"
	"hacko-app/internal/module/user/ports"
	"hacko-app/pkg"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/jwthandler"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.UserService = &userService{}

type userService struct {
	repo ports.UserRepository
	o    integOauth.Oauth2googleContract
}

func NewUserService(repo ports.UserRepository, o integOauth.Oauth2googleContract) *userService {
	return &userService{
		repo: repo,
		o:    o,
	}
}

func (s *userService) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {

	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to hash password")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to hash password"))
	}

	req.HassedPassword = hashed

	result, err := s.repo.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	var res = new(entity.LoginResponse)

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !pkg.ComparePassword(user.Pass, req.Password) {
		log.Warn().Any("payload", req).Msg("service::Login - Password not match")
		return nil, errmsg.NewCustomErrors(401, errmsg.WithMessage("Wrong email or password"))
	}

	// Generate Access Token
	accessToken, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(20 * time.Minute), // Validity period 20 minutes
	})
	if err != nil {
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed generate access token"))
	}

	// Generate Refresh Token
	refreshToken, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(14 * 24 * time.Hour), // Validity period 14 days
	})
	if err != nil {
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed generate refresh token"))
	}

	// Update Refresh Token in Database
	if err := s.repo.UpdateRefreshToken(ctx, user.Id, refreshToken); err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("service::Login - Failed to update refresh token")
		return nil, err
	}

	// Set Tokens in Response
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken

	return res, nil
}

func (s *userService) Profile(ctx context.Context, req *entity.ProfileRequest) (*entity.ProfileResponse, error) {
	user, err := s.repo.FindById(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetOauthGoogleUrl(ctx context.Context) (string, error) {
	url := s.o.GetUrl("state")

	return url, nil
}

func (s *userService) LoginGoogle(ctx context.Context, req *oauthgoogleent.UserInfoResponse) (*entity.LoginResponse, error) {
	var res = new(entity.LoginResponse)

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Handle custom errors first
		if errCustom, ok := err.(*errmsg.CustomError); ok {
			if errCustom.Code == 400 {
				// Create a request for registration
				registerReq := &entity.RegisterByGoogleRequest{
					Email:          req.Email,
					Name:           req.Name,
					GoogleId:       req.Id,
					ImageUrl:       *req.PicURL,
					Password:       "",
					HassedPassword: "",
				}

				_, err := s.repo.RegisterByGoogle(ctx, registerReq)
				if err != nil {
					log.Error().Err(err).Msg("service::loginGoogle - Failed to register a new user")
					return nil, err
				}

				user, err = s.repo.FindByEmail(ctx, req.Email)
				if err != nil {
					log.Error().Err(err).Msg("service::loginGoogle - Failed to find user after registration")
					return nil, err
				}
			} else {
				return nil, err
			}
		} else if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				log.Warn().Msg("service::loginGoogle - Unique violations, no additional actions")
			default:
				log.Error().Err(err).Any("payload", req).Msg("service::loginGoogle - Unhandled pq.Error")
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	accessToken, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(20 * time.Minute), // Validity period 20 minutes
	})

	if err != nil {
		log.Error().Err(err).Msg("service::loginGoogle - Failed to generate access tokens")
		return nil, err
	}

	refreshToken, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(14 * 24 * time.Hour), // Validity period 14 days
	})

	if err != nil {
		log.Error().Err(err).Msg("service::loginGoogle - Failed to generate refresh tokens")
		return nil, err
	}

	if err := s.repo.UpdateRefreshToken(ctx, user.Id, refreshToken); err != nil {
		log.Error().Err(err).Str("userId", user.Id).Msg("service::Login - Failed to update refresh token")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to save refresh token"))
	}

	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	return res, nil
}

func (s *userService) RefreshTokenService(ctx context.Context, refreshToken string) (string, error) {
	payload, err := s.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          payload.UserID,
		Role:            payload.Role,
		TokenExpiration: time.Now().Add(20 * time.Minute), // Validity period 20 minutes
	})
	if err != nil {
		log.Error().Err(err).Msg("service::RefreshTokenService - Failed to generate access token")
		return "", errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to generate access token"))
	}

	return accessToken, nil
}
