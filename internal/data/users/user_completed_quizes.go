package users

import (
	"github.com/Molnes/Nyhetsjeger/internal/data/quizzes"
	"github.com/google/uuid"
)

type UserCompletedQuiz struct {
	Quiz              quizzes.Quiz
	AnsweredQuestions []UserAnsweredQUestion
}

func GetCompletedQuizzes(userID uuid.UUID) ([]UserCompletedQuiz, error) {
	return SampleUserCompletedQuizzes, nil
}

var SampleUserCompletedQuizzes []UserCompletedQuiz = []UserCompletedQuiz{
	{
		Quiz: quizzes.SampleQuiz,
		AnsweredQuestions: []UserAnsweredQUestion{
			{
				UserID:            uuid.New(),
				QuestionID:        uuid.New(),
				ChosenAlternative: quizzes.SampleQuiz.Questions[0].Alternatives[0],
			},
		},
	},
}
