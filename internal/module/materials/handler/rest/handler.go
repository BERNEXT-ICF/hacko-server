package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/middleware"
	"hacko-app/internal/module/materials/entity"
	"hacko-app/internal/module/materials/ports"
	"hacko-app/internal/module/materials/repository"
	"hacko-app/internal/module/materials/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type materialsHandler struct {
	service ports.MaterialsService
}

func NewMaterialsHandler() *materialsHandler {
	var handler = new(materialsHandler)

	repo := repository.NewMaterialsRepository(adapter.Adapters.HackoPostgres)
	materialsService := service.NewMaterialsService(repo)

	handler.service = materialsService
	return handler
}

func (h *materialsHandler) Register(router fiber.Router) {
	router.Post("/class/:classId/materials", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.CreateMaterials)
	router.Patch("/class/materials/:materialsId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.UpdateMaterials)
	router.Delete("/class/materials/:materialsId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.DeleteMaterials)
}

func (h *materialsHandler) CreateMaterials(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateMaterialsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("classId")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::CreateMaterials - Failed to parsing id class")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.ClassId = reqId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateMaterials - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateMaterials - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateMaterials(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *materialsHandler) UpdateMaterials(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateMaterialsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	materialsId := c.Params("materialsId")

	reqMaterialId, err := strconv.Atoi(materialsId)
	if err != nil {
		log.Warn().Err(err).Msg("handler::UpdateMaterials - Failed to parsing id material")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id material"))))
	}

	req.MaterialId = reqMaterialId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateMaterials - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateMaterials - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateMaterials(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *materialsHandler) DeleteMaterials(c *fiber.Ctx) error{
	var (
		req = new(entity.DeleteMaterialsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	materialsId := c.Params("materialsId")

	reqMaterialId, err := strconv.Atoi(materialsId)
	if err != nil {
		log.Warn().Err(err).Msg("handler::UpdateClass - Failed to parsing id material")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.MaterialId = reqMaterialId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateMaterials - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err = h.service.DeleteMaterials(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(nil, ""))
}