package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	serverModels "gitlab.com/banter-bus/banter-bus-management-api/internal/server/models"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/service/models"
)

// AddPool adds a new user pool for an existing user.
func (env *Env) AddPool(
	_ *gin.Context,
	input *serverModels.QuestionPoolInput,
) (struct{}, error) {
	poolLogger := log.WithFields(log.Fields{
		"username": input.UserParams.Username,
	})
	poolLogger.Debug("Trying to add pool.")
	p := service.PoolService{DB: env.DB, Username: input.UserParams.Username, PoolName: input.PoolName}

	languageCode := input.Pool.LanguageCode
	if languageCode == "" {
		languageCode = "en"
	}

	emptyResponse := struct{}{}
	_, err := language.Parse(languageCode)
	if err != nil {
		log.Errorf("failed to parse language code %s, err %s", languageCode, err)
		return emptyResponse, errors.BadRequestf("invalid language code %s", languageCode)
	}

	err = p.AddPool(
		languageCode,
		input.Pool.GameName,
		input.Pool.Privacy,
	)

	if err != nil {
		poolLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Could not add pool.")
	}

	return emptyResponse, err
}

// GetAllPools returns all the user's questions pool.
func (env *Env) GetAllPools(_ *gin.Context, params *serverModels.UserParams) ([]serverModels.Pool, error) {
	userLogger := log.WithFields(log.Fields{
		"username": params.Username,
	})
	userLogger.Debug("Trying to get user pools.")
	p := service.PoolService{DB: env.DB, Username: params.Username}
	pools, err := p.GetPools()

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn(("User does not exist."))
		return []serverModels.Pool{}, errors.NotFoundf("The user %s", params.Username)
	}

	srvPools := env.getPools(pools)
	return srvPools, nil
}

// GetPool returns a single question pool for a specified user.
func (env *Env) GetPool(
	_ *gin.Context,
	params *serverModels.ExistingQuestionPoolParams,
) (serverModels.Pool, error) {
	userLogger := log.WithFields(log.Fields{
		"username":  params.Username,
		"pool_name": params.PoolName,
	})
	userLogger.Debug("Trying to get a single user pool.")
	p := service.PoolService{DB: env.DB, Username: params.Username, PoolName: params.PoolName}
	pool, err := p.GetPool()

	if err != nil {
		userLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Something went wrong, most likely because username or pool name does not exist.")
		return serverModels.Pool{}, err
	}

	singlePool := env.newQuestionPool(pool)
	return singlePool, err
}

func (env *Env) getPools(questionPools models.Pools) []serverModels.Pool {
	var pools []serverModels.Pool

	for _, pool := range questionPools {
		newPool := env.newQuestionPool(pool)
		pools = append(pools, newPool)
	}

	return pools
}

func (env *Env) newQuestionPool(pool models.Pool) serverModels.Pool {
	newPool := serverModels.Pool{
		PoolName:     pool.PoolName,
		GameName:     pool.GameName,
		LanguageCode: pool.LanguageCode,
		Privacy:      pool.Privacy,
	}
	return newPool
}

// RemovePool removes an existing question pool (for a specific user).
func (env *Env) RemovePool(
	_ *gin.Context,
	input *serverModels.ExistingQuestionPoolParams,
) (struct{}, error) {
	poolLogger := log.WithFields(log.Fields{
		"username":  input.UserParams.Username,
		"pool_name": input.PoolParams.PoolName,
	})
	poolLogger.Debug("Trying to remove pool.")
	p := service.PoolService{DB: env.DB, Username: input.UserParams.Username, PoolName: input.PoolParams.PoolName}

	err := p.RemovePool()
	if err != nil {
		poolLogger.WithFields(log.Fields{
			"err": err,
		}).Warn("Could not remove pool.")
	}

	emptyResponse := struct{}{}
	return emptyResponse, err
}

// AddQuestionToPool adds  a question to an existing pool.
func (env *Env) AddQuestionToPool(
	_ *gin.Context,
	input *serverModels.UpdateQuestionPoolInput,
) (struct{}, error) {
	poolLogger := log.WithFields(log.Fields{
		"username":  input.UserParams.Username,
		"pool_name": input.PoolParams.PoolName,
	})
	poolLogger.Debug("Trying to add question to pool.")

	question := env.newGenericQuestion(input.NewQuestion)
	p := service.PoolService{DB: env.DB, Username: input.UserParams.Username, PoolName: input.PoolParams.PoolName}

	err := p.AddPoolQuestion(question)
	if err != nil {
		poolLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Could not add question to pool.")
	}

	emptyResponse := struct{}{}
	return emptyResponse, err
}

// RemoveQuestionFromPool removes a question from an existing pool.
func (env *Env) RemoveQuestionFromPool(
	_ *gin.Context,
	input *serverModels.UpdateQuestionPoolInput,
) (struct{}, error) {
	poolLogger := log.WithFields(log.Fields{
		"username":  input.UserParams.Username,
		"pool_name": input.PoolParams.PoolName,
	})
	poolLogger.Debug("Trying to remove question from pool.")

	question := env.newGenericQuestion(input.NewQuestion)
	p := service.PoolService{DB: env.DB, Username: input.UserParams.Username, PoolName: input.PoolParams.PoolName}

	err := p.RemovePoolQuestion(question)
	if err != nil {
		poolLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Could not remove question from pool.")
	}

	emptyResponse := struct{}{}
	return emptyResponse, err
}
