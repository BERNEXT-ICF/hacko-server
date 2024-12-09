package handler

import (
	"hacko-app/internal/adapter"
	"hacko-app/internal/middleware"
	"hacko-app/internal/module/quiz/entity"
	"hacko-app/internal/module/quiz/ports"
	"hacko-app/internal/module/quiz/repository"
	"hacko-app/internal/module/quiz/service"
	"hacko-app/pkg/errmsg"
	"hacko-app/pkg/response"
	"strconv"

	// "strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type quizHandler struct {
	service ports.QuizService
}

func NewQuizHandler() *quizHandler {
	var handler = new(quizHandler)

	repo := repository.NewQuizRepository(adapter.Adapters.HackoPostgres)
	quizService := service.NewQuizService(repo)

	handler.service = quizService
	return handler
}

func (h *quizHandler) Register(router fiber.Router) {
	router.Post("/class/:classId/quiz", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.CreateQuiz)
	router.Post("/class/quiz/:quizId/question-quiz", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.CreateQuestionQuiz)
	router.Get("/class/:classId/quiz", middleware.AuthMiddleware, middleware.AuthRole([]string{"user", "admin", "teacher"}), h.GetAllQuiz)
}

func (h *quizHandler) CreateQuiz(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateQuizRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("classId")

	req.ClassId = id

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateQuiz(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *quizHandler) CreateQuestionQuiz(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateQuestionQuizRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	id := c.Params("quizId")

	quizId, err := strconv.Atoi(id) 
	if err != nil {
		log.Warn().Err(err).Msg("handler::UpdateClass - Failed to parse quizId")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errmsg.NewCustomErrors(500, errmsg.WithMessage("Failed to parse quizId"))))
	}

	req.QuizId = quizId 

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateAssignment - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateQuestionQuiz(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *quizHandler) GetAllQuiz(c *fiber.Ctx) error {
	var (
		req = new(entity.GetAllQuizRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	req.ClassId = c.Params("classId")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAllQuiz - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetAllQuiz(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}