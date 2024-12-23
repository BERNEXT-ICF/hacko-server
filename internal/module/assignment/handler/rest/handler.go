package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/middleware"
	"hacko-app/internal/module/assignment/entity"
	"hacko-app/internal/module/assignment/ports"
	"hacko-app/internal/module/assignment/repository"
	"hacko-app/internal/module/assignment/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type assignmentHandler struct {
	service ports.AssignmentService
}

func NewAssignmentHandler() *assignmentHandler {
	var handler = new(assignmentHandler)

	repo := repository.NewAssignmentRepository(adapter.Adapters.HackoPostgres)
	assignmentService := service.NewAssignmentService(repo)

	handler.service = assignmentService
	return handler
}

func (h *assignmentHandler) Register(router fiber.Router) {
	// admin routes
	router.Post("/class/:classId/assignment", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.CreateAssignment)
	router.Get("teacher/class/:classId/assignment", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetAllAssignmentByClassIdAdmin)
	router.Get("teacher/class/:classId/assignment/:assignmentId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetAssignmentDetailsAdmin)

	// user routes
	router.Get("/class/:classId/assignment", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetAllAssignmentByClassId)
	router.Get("/class/assignment/:assignmentId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetAssignmentDetails)
}

func (h *assignmentHandler) CreateAssignment(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateAssignmentRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("classId")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Failed to parsing id class")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.ClassId = reqId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateAssignment(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *assignmentHandler) GetAllAssignmentByClassId(c *fiber.Ctx) error{
	var (
		req = new(entity.GetAllAssignmentByClassIdRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	req.ClassId = c.Params("classId")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAllAssignmentByClassId - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetAllAssignmentByClassId(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))	
}

func (h *assignmentHandler) GetAssignmentDetails(c *fiber.Ctx) error{
	var (
		req = new(entity.GetAssignmentDetailsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("assignmentId")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::GetAssignmentDetails - Failed to parsing id class")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.AssignmentId = reqId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAssignmentDetails - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetAssignmentDetails(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))	
}

func (h *assignmentHandler) GetAllAssignmentByClassIdAdmin(c *fiber.Ctx) error{
	var (
		req = new(entity.GetAllAssignmentByClassIdAdminRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	req.ClassId = c.Params("classId")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAllAssignmentByClassIdAdmin - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetAllAssignmentByClassIdAdmin(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))	
}

func (h *assignmentHandler) GetAssignmentDetailsAdmin(c *fiber.Ctx) error{
	var (
		req = new(entity.GetAssignmentDetailsAdminRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("assignmentId")

	reqId, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::GetAssignmentDetailsAdmin - Failed to parsing id class")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse params id class"))))
	}

	req.AssignmentId = reqId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAssignmentDetailsAdmin - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetAssignmentDetailsAdmin(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))	
}