package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Validator interface
type Validator interface {
	Init() Validator
	AddError(field, message string)
	Validate(dataRules map[string]interface{}, data interface{}) map[string]string
	Unique(field string, data interface{}, storage []interface{})
}

type validator struct {
	Errors map[string]string
}

func (v *validator) Init() Validator {
	v.Errors = make(map[string]string)
	return v
}

func (v *validator) Validate(dataRules map[string]interface{}, data interface{}) map[string]string {
	for field, rawRules := range dataRules {
		rawStringRules, ok := rawRules.(string)
		if ok {
			rules := strings.Split(rawStringRules, "|")
			reverseRules := []string{}
			for i := len(rules) - 1; i >= 0; i-- {
				reverseRules = append(reverseRules, rules[i])
			}
			for _, rule := range reverseRules {
				value := reflect.ValueOf(data).Elem().FieldByName(strings.Title(field)).String()
				switch rule {
				case "required":
					v.required(field, value)
				case "email":
					v.email(field, value)
				// case "unique":
				// 	v.unique(field, value)
				case "forbiddenusernames", "forbiddenUsernames":
					v.forbiddenUsernames(field, value)
				}
				if strings.HasPrefix(rule, "len") {
					rawParams := strings.Replace(rule[4:len(rule)-1], " ", "", -1)
					params := strings.Split(rawParams, ",")
					if len(params) == 2 {
						// min, _ := strconv.Atoi(params[0])
						// max, _ := strconv.Atoi(params[1])
						min := params[0]
						max := params[1]
						v.len(field, value, min, max)
					}
				}
			}
		}
	}
	if len(v.Errors) > 0 {
		tmp := make(map[string]string)
		for k, v := range v.Errors {
			tmp[k] = v
		}
		v.Init()
		return tmp
	}
	return nil
}

func (v *validator) required(field string, value string) {
	if value == "" {
		v.Errors[field] = fmt.Sprintf("%s is required", field)
	}
}

func (v *validator) email(field string, value string) {
	if !govalidator.IsEmail(value) {
		v.Errors[field] = fmt.Sprintf("%s is not correct email", field)
	}
}

func (v *validator) forbiddenUsernames(field string, value string) {
	reserveRegexp := regexp.MustCompile("((?i)admin|(?i)админ)")
	if reserveRegexp.MatchString(value) {
		v.Errors[field] = fmt.Sprintf("%s %s is forbidden", field, value)
	}
}

func (v *validator) Unique(field string, data interface{}, storage []interface{}) {
	for _, s := range storage {
		sValue := reflect.ValueOf(s).Elem().FieldByName(strings.Title(field)).String()
		dataValue := reflect.ValueOf(data).Elem().FieldByName(strings.Title(field)).String()
		if sValue == dataValue {
			v.AddError(field, fmt.Sprintf("%s is unique", field))
			break
		}
	}
}

func (v *validator) len(field string, value string, min string, max string) {
	if !govalidator.RuneLength(value, min, max) {
		v.Errors[field] = fmt.Sprintf("%s must min len %s and max %s", field, min, max)
	}
}

func (v *validator) AddError(field, message string) {
	v.Errors[field] = message
}

// NewValidator create new validator
func NewValidator() Validator {
	return new(validator).Init()
}
