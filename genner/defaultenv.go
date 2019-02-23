package genner

import (
	"errors"
)

type GennerEnvironmentDefault struct {
	buffer      []string
	asInterface []interface{}
}

func (env *GennerEnvironmentDefault) SetInputBuffer(buffer []string) {
	env.buffer = buffer
	env.asInterface = nil
}

func (env *GennerEnvironmentDefault) OpReadInput([]interface{}) ([]interface{}, error) {
	if env.asInterface != nil {
		return env.asInterface, nil
	}

	if env.buffer == nil {
		return nil, errors.New("no input data to read")
	}

	env.asInterface = []interface{}{}
	for _, value := range env.buffer {
		env.asInterface = append(env.asInterface, value)
	}

	return env.asInterface, nil
}
