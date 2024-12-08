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
	router.Post("/class/assignment/:assignmentId/submission", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.SubmitAssignment)
	// router.Get("/class/:classId/assignment", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetAllAssignmentByClassId)
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
