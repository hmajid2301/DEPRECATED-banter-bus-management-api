package service

import (
	"github.com/juju/errors"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// PoolService is struct data required by all User question pool service functions.
type PoolService struct {
	DB       database.Database
	Username string
	PoolName string
}

// AddPool add question pool to a user
func (p *PoolService) AddPool(langCode string, gameName string, privacy string) error {
	u := UserService{DB: p.DB, Username: p.Username}
	exists, _ := u.doesUserExist()
	if !exists {
		return errors.NotFoundf("the user %s", p.Username)
	}

	currPool, _ := p.GetPool()
	if (currPool != models.Pool{}) {
		return errors.AlreadyExistsf("pool %s", p.PoolName)
	}

	pool := models.Pool{
		Username:     p.Username,
		PoolName:     p.PoolName,
		GameName:     gameName,
		LanguageCode: langCode,
		Privacy:      privacy,
	}

	inserted, err := pool.Add(p.DB)
	if !inserted && err == nil {
		return errors.Errorf("failed to add new pool %s", p.PoolName)
	}
	return err
}

// RemovePool deletes a  pool from an user
func (p *PoolService) RemovePool() error {
	_, err := p.GetPool()
	if err != nil {
		return err
	}

	filter := map[string]string{"username": p.Username, "pool_name": p.PoolName}
	questions := models.Questions{}
	removed, err := questions.Delete(p.DB, filter)
	if !removed || err != nil {
		return errors.Errorf("failed to remove pool %s", p.PoolName)
	}

	pools := models.Pools{}
	removed, err = pools.Delete(p.DB, filter)

	if !removed || err != nil {
		return errors.Errorf("failed to remove pool %s for user %s", p.PoolName, p.Username)
	}
	return err
}

// GetPool gets a single question pool from a user
func (p *PoolService) GetPool() (models.Pool, error) {
	u := UserService{DB: p.DB, Username: p.Username}
	exists, _ := u.doesUserExist()
	if !exists {
		return models.Pool{}, errors.NotFoundf("the user %s", p.Username)
	}

	var (
		filter = map[string]string{"username": p.Username, "pool_name": p.PoolName}
		pool   models.Pool
	)

	err := pool.Get(p.DB, filter)
	if err != nil {
		return models.Pool{}, errors.NotFoundf("the pool %s", p.PoolName)
	}

	return pool, err
}

// GetPools gets a specific user's (all of them) question pools
func (p *PoolService) GetPools() (models.Pools, error) {
	u := UserService{DB: p.DB, Username: p.Username}
	exists, _ := u.doesUserExist()
	if !exists {
		return models.Pools{}, errors.NotFoundf("the user %s", p.Username)
	}

	var (
		filter = map[string]string{"username": p.Username}
		pools  models.Pools
	)

	err := pools.Get(u.DB, filter)
	return pools, err
}

// AddPoolQuestion adds a question to a pool.
func (p *PoolService) AddPoolQuestion(question models.GenericQuestion) error {
	pool, err := p.GetPool()
	if err != nil {
		return err
	}

	exists := p.questionExist(question, pool.GameName)
	if exists {
		return errors.AlreadyExistsf("question '%s'", question.Content)
	}

	err = validateQuestion(pool.GameName, question)
	if err != nil {
		return err
	}

	t := true
	quest := models.Question{
		Username: p.Username,
		PoolName: p.PoolName,
		GameName: pool.GameName,
		Round:    question.Round,
		Enabled:  &t,
		Content: map[string]string{
			pool.LanguageCode: question.Content,
		},
	}

	if question.Group != nil {
		quest.Group.Name = question.Group.Name
		quest.Group.Type = question.Group.Type
	}

	inserted, err := quest.Add(p.DB)
	if !inserted || err != nil {
		return errors.Errorf("failed to add a new question")
	}

	return err
}

// RemovePoolQuestion removes a question from the pool.
func (p *PoolService) RemovePoolQuestion(question models.GenericQuestion) error {
	pool, err := p.GetPool()
	if err != nil {
		return err
	}

	exists := p.questionExist(question, pool.GameName)
	if !exists {
		return errors.NotFoundf("question '%s'", question.Content)
	}

	filter := p.filter(question, pool.GameName)
	deleted, err := p.DB.Delete("question", filter)
	if !deleted || err != nil {
		return errors.Errorf("failed to remove question")
	}

	return err
}

func (p *PoolService) questionExist(question models.GenericQuestion, gameName string) bool {
	filter := p.filter(question, gameName)
	currQuest := &models.Question{}
	q := QuestionService{DB: p.DB, Question: question, GameName: gameName}
	err := currQuest.Get(q.DB, filter)
	return err == nil
}

func (p *PoolService) filter(question models.GenericQuestion, gameName string) map[string]string {
	q := QuestionService{DB: p.DB, Question: question, GameName: gameName}
	filter := q.filter()
	filter["pool_name"] = p.PoolName
	filter["username"] = p.Username
	return filter
}
