package validator

import (
	"github.com/go-playground/validator/v10"
	"sync"
	"kago.fly/pkg/helper"
	"kago.fly/pkg/response"
)

type DataValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (c *DataValidator) Validate(i interface{}) error {
	c.lazyInit()
	err := c.validate.Struct(i)
	if err != nil {
		translator, _ := helper.FindTranslator("en")
		for _, err := range err.(validator.ValidationErrors) {
			panic(response.NewParamError(err.Translate(translator)))
		}
	}
	return nil
}

func (c *DataValidator) lazyInit() {
	c.once.Do(func() {
		c.validate = helper.GetValidator()
	})
}
