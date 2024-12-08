package rest

import (
	"fmt"
	"hacko-app/internal/adapter"
	"hacko-app/internal/infrastructure/config"
	integOauth "hacko-app/internal/integration/oauth2google"
	"hacko-app/internal/middleware"
	"hacko-app/internal/module/user/entity"
	"time"

	"hacko-app/internal/module/user/ports"
	"hacko-app/internal/module/user/repository"
	"hacko-app/internal/module/user/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	service     ports.UserService
	integration integOauth.Oauth2googleContract
}

func NewUserHandler(o integOauth.Oauth2googleContract) *userHandler {
	var handler = new(userHandler)

	repo := repository.NewUserRepository(adapter.Adapters.HackoPostgres)
	service := service.NewUserService(repo, o)

	handler.integration = o

	handler.service = service

	return handler
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/login", h.login)
	router.Post("/refresh", h.refresh)
	router.Delete("/logout", h.logout)

	// route login || register by google
	router.Get("/oauth/google/url", h.oauthGoogleUrl)
	router.Get("/signin/callback", h.callbackSigninGoogle)

	// route user service
	router.Get("/profile", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "teacher"}), h.profile)
	router.Get("/profile/:user_id", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.profileByUserId)
}

func (h *userHandler) register(c *fiber.Ctx) error {
	var (
		req = new(entity.RegisterRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Register(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *userHandler) login(c *fiber.Ctx) error {
	var (
		req = new(entity.LoginRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	// Parsing request body
	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	// Validasi request body
	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	// Memanggil service untuk login
	res, err := h.service.Login(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	// Set cookie for accessToken
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(20 * time.Minute), // Validity period 20 minutes
		HTTPOnly: true,
		Secure:   false,
		SameSite: "none",
	})

	// Only set refresh token cookie if req.Remember is true
	if req.Remember {
		c.Cookie(&fiber.Cookie{
			Name:     "refreshToken",
			Value:    res.RefreshToken,
			Expires:  time.Now().Add(14 * 24 * time.Hour), // Validity period 14 days
			HTTPOnly: true,
			Secure:   false,
			SameSite: "none",
		})
	}

	// Return response without token
	return c.Status(fiber.StatusOK).JSON(response.Success(nil, "Login successful"))
}

func (h *userHandler) profileByUserId(c *fiber.Ctx) error {
	var (
		req = new(entity.ProfileRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.UserId = c.Params("user_id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::profileByUserId - Invalid Request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Profile(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) profile(c *fiber.Ctx) error {
	var (
		req = new(entity.ProfileRequest)
		ctx = c.Context()
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	res, err := h.service.Profile(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) oauthGoogleUrl(c *fiber.Ctx) error {
	referer := c.Get("Referer")
	if referer == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request: Referer is missing from the request headers"))))
	}
	return c.Redirect(h.integration.GetUrl(referer), http.StatusTemporaryRedirect)
}

func (h *userHandler) callbackSigninGoogle(c *fiber.Ctx) error {
	var ctx = c.Context()

	state, code := c.FormValue("state"), c.FormValue("code")
	if state == "" || code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request: state or code missing"))))
	}

	token, err := h.integration.Exchange(ctx, code)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	// verify token with google
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.Envs.Oauth.Google.ClientId,
	})
	_, err = verifier.Verify(ctx, token.Extra("id_token").(string))
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	userInfo, err := h.integration.GetUserInfo(ctx, token)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.LoginGoogle(ctx, &userInfo)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	accessToken := res.AccessToken
	refreshToken := res.RefreshToken

	// Set cookie for accessToken
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Expires:  time.Now().Add(20 * time.Minute), // Validity period 20 minutes
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	// Set cookie for refreshToken
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(14 * 24 * time.Hour), // Validity period 14 days
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	// Dekode URL state
	redirectURL, err := url.QueryUnescape(state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid state parameter"))
	}

	// Redirect to frontend with path dashboard
	finalRedirect := fmt.Sprintf("%s/dashboard", redirectURL)
	return c.Redirect(finalRedirect, fiber.StatusTemporaryRedirect)
}

func (h *userHandler) refresh(c *fiber.Ctx) error {
	var ctx = c.Context()
	refreshToken := c.Cookies("refreshToken")

	if refreshToken == "" {
		log.Warn().Msg("handler::refresh - Refresh token not provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Refresh token not provided",
			"success": false,
		})
	}

	// service refresh token
	accessToken, err := h.service.RefreshTokenService(ctx, refreshToken)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Expires:  time.Now().Add(20 * time.Minute), // Validity period 20 minutes
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, "Refresh access token successful"))
}

func (h *userHandler) logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  time.Now().Add(-time.Hour),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, "Successfully logged out"))
}
