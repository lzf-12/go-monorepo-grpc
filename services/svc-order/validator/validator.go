package validator

import (
	"fmt"
	"log"
	"ops-monorepo/services/svc-order/internal/delivery/types"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type IValidator interface {
	ValidateOrderItems(items []types.StockItemRequest) ([]map[string]interface{}, error)
}

type Validator struct {
	instance   *validator.Validate
	translator ut.Translator
}

func NewValidator() IValidator {

	// init instance of 'validate' with sane defaults
	// init default translator
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Panic("translator not found")
	}

	// register custom translation for message
	registerCustomTranslations(validate, trans)

	// register tag e.Field() use json tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("unique", populateUniqueList)

	return &Validator{
		instance:   validate,
		translator: trans,
	}
}

func registerCustomTranslations(validate *validator.Validate, trans ut.Translator) {

	// required
	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	// gt > 0
	_ = validate.RegisterTranslation("gt", trans, func(ut ut.Translator) error {
		return ut.Add("gt", "{0} must be greater than {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gt", fe.Field(), fe.Param())

		// special case for gt=0
		if fe.Param() == "0" {
			return fmt.Sprintf("%s must be greater than zero", fe.Field())
		}
		return t
	})
}

func (m *Validator) ValidateOrderItems(items []types.StockItemRequest) ([]map[string]interface{}, error) {

	errList := make([]map[string]interface{}, 0)
	duplicateList := map[string]bool{}
	uniqueIndexList := map[int]bool{}

	for i, item := range items {
		errObj := map[string]interface{}{}
		err := m.instance.Struct(item)
		if err != nil {
			// not accepted error by validator
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return nil, err
			}

			validationErrors, ok := err.(validator.ValidationErrors)
			if ok {
				for _, ve := range validationErrors {

					// handle tag unique, check duplicate
					if ve.Tag() == "unique" {

						value := ve.Value().(string)
						if duplicateList[value] {
							errObj[ve.Field()] = ErrMsgFieldShouldUnique
							errObj["row"] = i + 1

							if !uniqueIndexList[i] {
								errList = append(errList, errObj)
								uniqueIndexList[i] = true
							}

						}

						duplicateList[value] = true
					}

					if ve.Tag() != "unique" {
						errObj[ve.Field()] = ErrMsgFieldShouldUnique
						errObj["row"] = i + 1

						if !uniqueIndexList[i] {
							errList = append(errList, errObj)
							uniqueIndexList[i] = true
						}
					}
				}
			}
		}

		// return if all row error checked
		if i+1 == len(items) {
			return errList, err
		}
	}

	return errList, nil
}

func populateUniqueList(fl validator.FieldLevel) bool {
	return false
}
