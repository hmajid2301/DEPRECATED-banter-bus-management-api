package internal

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

type GameParams struct {
	Name string `description:"The name of the game." example:"quibly" path:"name" validate:"required"`
}

type RoundParams struct {
	Round string `description:"Name of the round for a game." validate:"required" example:"free_form" query:"round"`
}

// unmarshalBSONToStruct is a custom unmarshal function, that will unmarshal BSON to structs.
// This function is to be used when a subfield can take multiple types i.e. storing question
// data differently for different games. It will unmarshal that field (polymorphic one)
// i.e. `Questions`, into BSON raw data. This can then be cast into the correct struct for
// the polymorphic field.
//
// The first unmarshal gets the Raw BSON data. The Raw BSON data allows us to unmarshal sub-objects like `Questions`
// field to a specific struct.
//
// The second unmarshal converts the raw BSON data into a struct i.e. `QuestionPool`, note in this example `Questions`
//  field is type `interface{}`.
//
// Next, we unmarshal the subField into raw BSON data, in the example above this would be the `Questions` field.
// This way we only have raw BSON data related to that field and can be cast appropriate.
func UnmarshalBSONToStruct(data []byte, structType interface{}, subField interface{}) error {
	err := bson.Unmarshal(data, structType)
	if err != nil {
		return err
	}

	err = bson.Unmarshal(data, subField)
	if err != nil {
		return err
	}

	return nil
}

// Refer to the BSON function above, this function works almost exactly the same as the one above.
func UnmarshalJSONToStruct(data []byte, structType interface{}, subField interface{}) error {
	err := json.Unmarshal(data, structType)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, subField)
	if err != nil {
		return err
	}

	return nil
}

func GetEnabledBool(enabledStr string) *bool {
	var (
		t = true
		f = false
	)

	var n *bool
	filters := map[string]*bool{
		"enabled":  &t,
		"disabled": &f,
		"all":      n,
	}

	enabled := filters[enabledStr]
	return enabled
}
