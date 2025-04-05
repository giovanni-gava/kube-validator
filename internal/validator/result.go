package validator

type ValidationResult struct {
	Message string `json:"message"`
	Level   string `json:"level"` // "error", "warning", "suggestion"
}
