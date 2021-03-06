package validation

import (
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func validateString(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	_, ok := value.(string)
	return ok
}

func validateDigits(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	str, ok := value.(string)
	if ok {
		return regexDigits.FindAllString(str, 1) == nil
	}
	return false
}

func validateAlpha(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	parameters = []string{patternAlpha}
	return validateRegex(field, value, parameters, form)
}

func validateAlphaDash(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	parameters = []string{patternAlphaDash}
	return validateRegex(field, value, parameters, form)
}

func validateAlphaNumeric(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	parameters = []string{patternAlphaNumeric}
	return validateRegex(field, value, parameters, form)
}

func validateEmail(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	parameters = []string{patternEmail}
	return validateRegex(field, value, parameters, form)
}

func validateStartsWith(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	RequireParametersCount("starts_with", parameters, 1)
	str, ok := value.(string)
	if ok {
		for _, prefix := range parameters {
			if strings.HasPrefix(str, prefix) {
				return true
			}
		}
	}
	return false
}

func validateEndsWith(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	RequireParametersCount("ends_with", parameters, 1)
	str, ok := value.(string)
	if ok {
		for _, prefix := range parameters {
			if strings.HasSuffix(str, prefix) {
				return true
			}
		}
	}
	return false
}

func validateIP(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	str, ok := value.(string)
	if ok {
		ip := net.ParseIP(str)
		if ip != nil {
			form[field] = ip
			return true
		}
	}

	return false
}

func validateIPv4(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	if validateIP(field, value, parameters, form) {
		return form[field].(net.IP).To4() != nil
	}
	return false
}

func validateIPv6(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	if validateIP(field, value, parameters, form) {
		return form[field].(net.IP).To4() == nil
	}
	return false
}

func validateJSON(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	str, ok := value.(string)
	if ok {
		var data interface{}
		err := json.Unmarshal([]byte(str), &data)
		if err == nil {
			form[field] = data
			return true
		}
	}
	return false
}

func validateRegex(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	RequireParametersCount("regex", parameters, 1)
	str, ok := value.(string)
	if ok {
		return regexp.MustCompile(parameters[0]).MatchString(str)
	}
	return false
}

func validateTimezone(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	tz, ok := value.(string)
	if ok && tz != "Local" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			form[field] = loc
			return true
		}
	}
	return false
}

func validateURL(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	str, ok := value.(string)
	if ok {
		url, err := url.ParseRequestURI(str)
		if err == nil {
			form[field] = url
			return true
		}
	}
	return false
}

func validateUUID(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	str, ok := value.(string)
	if ok {
		id, err := uuid.Parse(str)
		if err == nil {
			if len(parameters) == 1 {
				version, err := strconv.Atoi(parameters[0])
				if err == nil && id.Version() != uuid.Version(version) {
					return false
				}
			}
			form[field] = id
			return true
		}
	}
	return false
}
