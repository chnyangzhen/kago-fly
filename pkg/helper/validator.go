package helper

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
)

func init() {
	validate = validator.New()
	en_trans()
	//zh_trans()
}

func en_trans() {
	en := en.New()
	zh := zh.New()
	uni = ut.New(en, en, zh)
	transEn, _ := uni.GetTranslator("en")
	transZh, _ := uni.GetTranslator("zh")
	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, transEn)
	zh_translations.RegisterDefaultTranslations(validate, transZh)
}

func GetValidator() *validator.Validate {
	return validate
}

func ValidateObj(o interface{}) error {
	return validate.Struct(o)
}

func FindTranslator(locale string) (ut.Translator, bool) {
	return uni.FindTranslator(locale)
}

func ValidateString(str string, rule string) error {
	return validate.Var(str, rule)
}
