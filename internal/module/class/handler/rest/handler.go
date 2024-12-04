package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/module/class/entity"
	"hacko-app/internal/module/class/ports"
	"hacko-app/internal/module/class/repository"
	"hacko-app/internal/module/class/service"
	"hacko-app/internal/middleware"
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
	router.Get("/class", h.GetAllClasses)
	router.Get("/class/:id", h.GetClassById)
	
	router.Post("/class", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.CreateClassregister)
	router.Post("/class/enroll", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.EnrollClass)
}

func (h *classHandler) CreateClassregister(c *fiber.Ctx) error {
	var(
		req = new(entity.CreateClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator  
		l = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()
	
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

func (h *classHandler) GetAllClasses(c *fiber.Ctx) error {
	var ctx = c.Context()

	classes, err := h.service.GetAllClasses(ctx)
	if err != nil {
		log.Error().Err(err).Msg("handler::GetAllClasses - Failed to get classes from service")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(classes, "Successfully retrieved all classes"))
}

func (h *classHandler) GetClassById(c *fiber.Ctx) error {
	var ctx = c.Context()
	classId := c.Params("id")

	req := &entity.GetClassByIdRequest{
		Id: classId,
	}

	res, err := h.service.GetClassById(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *classHandler) EnrollClass(c *fiber.Ctx) error {
	var(
		req = new(entity.EnrollClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator  
		l = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request body"))))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.EnrollClass(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, "Successfully enrolled in the class"))
}
