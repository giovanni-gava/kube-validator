package validator

type ValidationResult struct {
	Message string `json:"message"`
	Level   string `json:"level"`
	DocLink string `json:"docLink,omitempty"`
}
