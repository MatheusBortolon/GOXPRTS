package validator

// IsValidCEP checks if the CEP is an 8-digit string.
func IsValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for _, r := range cep {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
