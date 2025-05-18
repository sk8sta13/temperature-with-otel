package validators

import "regexp"

func IsValidZipCode(cep string) bool {
	if cep == "" {
		return false
	}
	re := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return re.MatchString(cep)
}
