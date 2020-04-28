package translate

import (
	"cloud/lib/logger"
	"log"
	"reflect"
	"strings"
	"test/null"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type User struct {
	Email    string    `json:"email" validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Password string    `json:"password" validate:"required"`
	N        int       `json:"n" validate:"required,min=10"`
	Time     time.Time `json:"time" validate:"required"`
	NullTime null.Time `json:"null_time" validate:"required"`
}

func Translate() {
	v := validator.New()

	en := en.New()
	uni := ut.New(en, en)
	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	a := User{
		Email:    "something@gmail.com",
		Name:     "",
		Password: "p",
		N:        1,
	}
	if err := v.Struct(a); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			logger.Debug(e.Translate(trans))
		}

	}
}
