package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Molnes/Nyhetsjeger/internal/config"
	"github.com/Molnes/Nyhetsjeger/internal/models/users/user_quiz"
	"github.com/Molnes/Nyhetsjeger/internal/models/users/user_quiz_summary"
	utils "github.com/Molnes/Nyhetsjeger/internal/utils"
	"github.com/Molnes/Nyhetsjeger/internal/web_server/web/views/components/quiz_components/play_quiz_components"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type publicApiHandler struct {
	sharedData *config.SharedData
}

// Creates a new PublicApiHandler
func NewPublicApiHandler(sharedData *config.SharedData) *publicApiHandler {
	return &publicApiHandler{sharedData}
}

// Registers the public api handlers to the given echo group
func (h *publicApiHandler) RegisterPublicApiHandlers(g *echo.Group) {
	g.POST("/user-answer", h.postAnswer)
	g.GET("/question", h.getQuestion)
}

func (h *publicApiHandler) postAnswer(c echo.Context) error {
	questionID, err := uuid.Parse(c.QueryParam("question-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or missing question-id")
	}
	pickedAnswerID, err := uuid.Parse(c.FormValue("answer-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or missing answer-id in formdata")
	}

	questionPresentedAt, err := time.Parse(time.RFC3339, c.FormValue("last_question_presented_at"))
	if err != nil {
		return err
	}

	answered, err := user_quiz.AnswerQuestionGuest(h.sharedData.DB, questionID, pickedAnswerID, questionPresentedAt)
	if err != nil {
		return err
	}

	publicQuizId, err := user_quiz.GetOpenQuizId(h.sharedData.DB)
	if err != nil {
		return err
	}
	if publicQuizId != answered.Question.QuizID {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot answer question in non-open quiz without being authenticated.")
	}

	summaryrow := user_quiz_summary.AnsweredQuestion{
		QuestionID:            questionID,
		QuestionText:          answered.Question.Text,
		MaxPoints:             answered.Question.Points,
		ChosenAlternativeID:   answered.ChosenAnswerID,
		ChosenAlternativeText: answered.Question.GetAnswerTextById(answered.ChosenAnswerID),
		IsCorrect:             answered.Question.IsAnswerCorrect(answered.ChosenAnswerID),
		PointsAwarded:         answered.PointsAwarded,
	}

	return utils.Render(c, http.StatusOK, play_quiz_components.FeedbackButtonsWithClientState(answered, &summaryrow))

}

func (h *publicApiHandler) getQuestion(c echo.Context) error {
	openQuizId, err := user_quiz.GetOpenQuizId(h.sharedData.DB)
	if err != nil {
		return err
	}

	quizId, err := uuid.Parse(c.QueryParam("quiz-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Ugyldig eller manglende quiz-id")
	}
	if quizId != openQuizId {
		return echo.NewHTTPError(http.StatusNotFound, "Ingen åpen quiz med den angitte ID-en")
	}

	currentQuestion, err := strconv.ParseUint(c.QueryParam("current-question"), 10, 64)
	if err != nil || currentQuestion < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Ugyldig eller manglende såørsmål nummer")
	}

	data, err := user_quiz.GetQuestionByNumberInQuiz(h.sharedData.DB, quizId, uint(currentQuestion))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Ingen spørsmål med det angitte nummeret")
		}
		return err
	}

	return utils.Render(c, http.StatusOK, play_quiz_components.QuizPlayContent(data))

}
