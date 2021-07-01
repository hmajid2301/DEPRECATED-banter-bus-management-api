package questions

type Drawlosseum struct{}

func (d Drawlosseum) ValidateQuestion(_ QuestionIn) error {
	return nil
}

func (d Drawlosseum) HasGroups(_ string) bool {
	return false
}
