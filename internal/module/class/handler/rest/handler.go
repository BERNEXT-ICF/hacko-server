package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/module/class/entity"
	"hacko-app/internal/module/class/ports"
	"hacko-app/internal/module/class/repository"
	"hacko-app/internal/module/class/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type classHandler struct {
	service ports.ClassService
}

// NewClassHandler creates a new instance of classHandler
func NewClassHandler() *classHandler {
	var handler = new(classHandler)

	repo := repository.NewClassRepository(adapter.Adapters.HackoPostgres)
	classService := service.NewClassService(repo)

	handler.service = classService
	return handler
}

func (h *classHandler) Register(router fiber.Router) {
	router.Post("/class", h.CreateClassregister)
}

func (h *classHandler) CreateClassregister(c *fiber.Ctx) error {
	var(
		req = new(entity.CreateClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateClass(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}
