package fibbingit

import (
	"fmt"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/core/models"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// FibbingIt type that implements PlayableGames.
type FibbingIt struct {
	CurrentQuestion models.GenericQuestion
}

// AddGame is used to add new games of type `fibbing_it`.
func (f FibbingIt) AddGame(rulesURL string) (bool, error) {
	t := true
	var newGame = &models.Game{
		Name:     "fibbing_it",
		RulesURL: rulesURL,
		Enabled:  &t,
		Questions: &models.FibbingItQuestions{
			Opinion:  map[string]map[string][]models.Question{},
			FreeForm: map[string][]models.Question{},
			Likely:   []models.Question{},
		},
	}

	inserted, err := database.Insert("game", newGame)
	return inserted, err
}

// GetQuestionPath gets the path to get a specific question in MongoDB. Using string concat i.e. "question.likely".
func (f FibbingIt) GetQuestionPath() string {
	question := f.CurrentQuestion
	questionPath := fmt.Sprintf("questions.%s", question.Round)

	if question.Group.Name != "" {
		questionPath += fmt.Sprintf(".%s", question.Group.Name)
	}
	if question.Group.Type != "" {
		questionPath += fmt.Sprintf(".%s", question.Group.Type)
	}
	return questionPath
}

// ValidateQuestionInput is used to validate input for interacting with questions.
func (f FibbingIt) ValidateQuestionInput() error {
	question := f.CurrentQuestion
	validRounds := map[string]bool{"opinion": true, "likely": true, "free_form": true}
	validTypes := map[string]bool{"answers": true, "questions": true}

	if !validRounds[question.Round] {
		return errors.BadRequestf("Invalid round %s", question.Round)
	}

	if question.Group == nil {
		return errors.BadRequestf("Missing group information %s", question.Group)
	} else if question.Group.Type != "" && !validTypes[question.Group.Type] {
		return errors.BadRequestf("Invalid group type %s", question.Group.Type)
	}

	return nil
}

// QuestionPoolToGenericQuestions converts question pool questions into generic questions that can be returned back to
// a client.
func (f FibbingIt) QuestionPoolToGenericQuestions(questions interface{}) ([]models.GenericQuestion, error) {
	fibbingItQuestions, ok := questions.(models.FibbingItQuestionsPool)
	if !ok {
		errorMessage := "Failed to convert type to FibbingItQuestionsPool."
		log.Error(errorMessage)
		return []models.GenericQuestion{}, errors.New(errorMessage)
	}

	var newGenericQuestions []models.GenericQuestion
	likelyQuestions := likelyQuestionsToGenericQuestion(fibbingItQuestions.Likely)
	newGenericQuestions = append(newGenericQuestions, likelyQuestions...)
	freeFormQuestions := freeFormQuestionsToGenericQuestion(fibbingItQuestions.FreeForm)
	newGenericQuestions = append(newGenericQuestions, freeFormQuestions...)
	opinionQuestions := opinionQuestionsToGenericQuestions(fibbingItQuestions.Opinion)
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
						Type: groupType,
					},
				}
				newGenericQuestions = append(newGenericQuestions, question)
			}
		}
	}
	return newGenericQuestions
}
