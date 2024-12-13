package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/middleware"
	"hacko-app/internal/module/submission/entity"
	"hacko-app/internal/module/submission/ports"
	"hacko-app/internal/module/submission/repository"
	"hacko-app/internal/module/submission/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"
	// "strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type submissionHandler struct {
	service ports.SubmissionService
}

func NewSubmissionHandler() *submissionHandler {
	var handler = new(submissionHandler)

	repo := repository.NewSubmissionRepository(adapter.Adapters.HackoPostgres)
	submissionService := service.NewSubmissionService(repo)

	handler.service = submissionService
	return handler
}



func (h *submissionHandler) Register(router fiber.Router) {
	// user routes
	router.Post("/class/assignment/:assignmentId/submission", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.SubmitAssignment)

	// admin routes
	router.Get("/class/assignment/submission/:submissionId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetSubmissionDetails)
	router.Post("/class/assignment/submission/:submissionId", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GradingSubmission)
}

func (h *submissionHandler) SubmitAssignment(c *fiber.Ctx) error {
	var (
		req = new(entity.SubmitRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("assignmentId")

	req.AssignmentId = id

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.SubmitAssignment(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *submissionHandler) GetSubmissionDetails(c *fiber.Ctx) error {
	var (
		req = new(entity.GetSubmissionDetailsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("submissionId")

	req.SubmissionId = id

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetSubmissionDetails - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetSubmissionDetails(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *submissionHandler) GradingSubmission(c *fiber.Ctx) error {
	var (
		req = new(entity.GradingSubmissionRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("submissionId")

	req.SubmissionId = id

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GradingSubmission - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::GradingSubmission - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GradingSubmission(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}
