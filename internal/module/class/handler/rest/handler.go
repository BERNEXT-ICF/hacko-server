package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/middleware"
	"hacko-app/internal/module/class/entity"
	"hacko-app/internal/module/class/ports"
	"hacko-app/internal/module/class/repository"
	"hacko-app/internal/module/class/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type classHandler struct {
	service ports.ClassService
}

func NewClassHandler() *classHandler {
	var handler = new(classHandler)

	repo := repository.NewClassRepository(adapter.Adapters.HackoPostgres)
	classService := service.NewClassService(repo)

	handler.service = classService
	return handler
}

func (h *classHandler) Register(router fiber.Router) {
	// route public
	router.Get("/class", h.GetAllClasses)
	router.Get("/class/:id", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetOverviewClassById)

	// route user
	router.Post("/class/enroll", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.EnrollClass)

	// route teacher, manage class
	router.Post("/class", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.CreateClassregister)
	router.Put("/class/:id", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "teacher"}), h.UpdateClass)
	router.Delete("/class/:id", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "teacher"}), h.DeleteClass)
	router.Patch("/class/:id", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "teacher"}), h.UpdateVisibilityClass)
}

func (h *classHandler) CreateClassregister(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
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
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(classes, "Successfully retrieved all classes"))
}

func (h *classHandler) GetOverviewClassById(c *fiber.Ctx) error {
	var ctx = c.Context()
	classId := c.Params("id")
	var l = middleware.GetLocals(c)

	req := &entity.GetOverviewClassByIdRequest{
		UserId: l.GetUserId(),
		Id:     classId,
	}

	res, err := h.service.GetOverviewClassById(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *classHandler) EnrollClass(c *fiber.Ctx) error {
	var (
		req = new(entity.EnrollClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request body"))))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::EnrollClass - Invalid request body")
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

func (h *classHandler) UpdateClass(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()
	id := c.Params("id")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::UpdateClass - Failed to parsing id class")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.Id = reqId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateClass - Failed to parsing body request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request body"))))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateClass - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateClass(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *classHandler) DeleteClass(c *fiber.Ctx) error {
	var (
		req = new(entity.DeleteClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()
	id := c.Params("id")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::DeleteClass - Failed to parsing id class")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.Id = reqId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::DeleteClass - Failed to parsing body request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request body"))))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::DeleteClass - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err = h.service.DeleteClass(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, "Delete Class Successful"))
}

func (h *classHandler) UpdateVisibilityClass(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateVisibilityClassRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()
	id := c.Params("id")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::DeleteClass - Failed to parsing id class")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.Id = reqId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::DeleteClass - Failed to parsing body request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request body"))))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::DeleteClass - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}
	res, err := h.service.UpdateVisibilityClass(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
