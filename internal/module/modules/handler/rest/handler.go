package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/middleware"
	"hacko-app/internal/module/modules/entity"
	"hacko-app/internal/module/modules/ports"
	"hacko-app/internal/module/modules/repository"
	"hacko-app/internal/module/modules/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type modulesHandler struct {
	service ports.ModulesService
}

func NewModulesHandler() *modulesHandler {
	var handler = new(modulesHandler)

	repo := repository.NewModulesRepository(adapter.Adapters.HackoPostgres)
	modulesService := service.NewModulesService(repo)

	handler.service = modulesService
	return handler
}

func (h *modulesHandler) Register(router fiber.Router) {
	router.Post("/class/materials/:materialsId/modules", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.CreateModules)
	router.Put("/class/materials/modules/:modulesId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.UpdateModules)
	router.Delete("/class/materials/modules/:modulesId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.DeleteModules)
}

func (h *modulesHandler) CreateModules(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateModulesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("materialsId")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::CreateMaterials - Failed to parsing id materials")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id materials"))))
	}

	req.MaterialsId = reqId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateMaterials - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateMaterials - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateModules(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *modulesHandler) UpdateModules(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateModulesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("modulesId")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::UpdateModules - Failed to parsing id modules")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id modules"))))
	}

	req.ModulesId = reqId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateModules - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateModules - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateModules(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *modulesHandler) DeleteModules(c *fiber.Ctx) error {
	var (
		req = new(entity.DeleteModulesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("modulesId")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::DeleteModules - Failed to parsing id modules")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id modules"))))
	}

	req.ModulesId = reqId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::DeleteModules - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err = h.service.DeleteModules(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(nil, "successfully deleted module"))
}