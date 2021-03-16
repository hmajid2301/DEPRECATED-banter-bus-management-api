package games

import (
	"fmt"
	"strings"

	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// FibbingIt is the concrete type for the the game interface.
type FibbingIt struct{}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.likely".
func (f FibbingIt) GetQuestionPath(question models.GenericQuestion) string {
	questionPath := fmt.Sprintf("questions.%s", question.Round)

	if question.Group.Name != "" {
		questionPath += fmt.Sprintf(".%s", question.Group.Name)
	}
	if question.Group.Type != "" {
		questionPath += fmt.Sprintf(".%ss", question.Group.Type)
	}
	return questionPath
}

// ValidateQuestion is used to validate input for interacting with questions.
func (f FibbingIt) ValidateQuestion(question models.GenericQuestion) error {
	validRounds := map[string]bool{"opinion": true, "likely": true, "free_form": true}
	validTypes := map[string]bool{"answer": true, "question": true}

	if !validRounds[question.Round] {
		return errors.BadRequestf("invalid round %s", question.Round)
	}

	if question.Group == nil {
		return errors.BadRequestf("missing group information %s", question.Group)
	} else if question.Group.Type != "" && !validTypes[question.Group.Type] {
		return errors.BadRequestf("invalid group type %s", question.Group.Type)
	}

	return nil
}

// NewQuestionPool gets the question pool structure for the Fibbing It game.
func (f FibbingIt) NewQuestionPool() models.QuestionPoolType {
	fibbingItQuestionPool := &models.FibbingItQuestionsPool{}
	fibbingItQuestionPool.NewPool()
	return fibbingItQuestionPool
}

// QuestionPoolToGenericQuestions converts question pool questions into generic questions that can be returned back to
// a client.
func (f FibbingIt) QuestionPoolToGenericQuestions(questions models.QuestionPoolType) ([]models.GenericQuestion, error) {
	fibbingIt, ok := questions.(*models.FibbingItQuestionsPool)
	if !ok {
		return nil, errors.Errorf("invalid question for Fibbing It")
	}
	var newGenericQuestions []models.GenericQuestion

	likelyQuestions := likelyQuestionsToGenericQuestion(fibbingIt.Likely)
	newGenericQuestions = append(newGenericQuestions, likelyQuestions...)
	freeFormQuestions := freeFormQuestionsToGenericQuestion(fibbingIt.FreeForm)
	newGenericQuestions = append(newGenericQuestions, freeFormQuestions...)
	opinionQuestions := opinionQuestionsToGenericQuestions(fibbingIt.Opinion)
	newGenericQuestions = append(newGenericQuestions, opinionQuestions...)

	return newGenericQuestions, nil
}

func likelyQuestionsToGenericQuestion(likely []string) []models.GenericQuestion {
	var newGenericQuestions []models.GenericQuestion
	for _, content := range likely {
		question := models.GenericQuestion{
			Content: content,
			Round:   "likely",
		}
		newGenericQuestions = append(newGenericQuestions, question)
	}

	return newGenericQuestions
}

func freeFormQuestionsToGenericQuestion(freeform map[string][]string) []models.GenericQuestion {
	var newGenericQuestions []models.GenericQuestion
	for groupName, questionGroup := range freeform {
		for _, content := range questionGroup {
			question := models.GenericQuestion{
				Content: content,
				Round:   "free_form",
				Group: &models.GenericQuestionGroup{
					Name: groupName,
				},
			}
			newGenericQuestions = append(newGenericQuestions, question)
		}
	}
	return newGenericQuestions
}

func opinionQuestionsToGenericQuestions(opinion map[string]map[string][]string) []models.GenericQuestion {
	var newGenericQuestions []models.GenericQuestion
	for groupName, questionGroup := range opinion {
		for groupType, questions := range questionGroup {
			for _, content := range questions {
				question := models.GenericQuestion{
					Content: content,
					Round:   "opinion",
					Group: &models.GenericQuestionGroup{
						Name: groupName,
						Type: strings.TrimSuffix(groupType, "s"),
					},
				}
				newGenericQuestions = append(newGenericQuestions, question)
			}
		}
	}
	return newGenericQuestions
}
