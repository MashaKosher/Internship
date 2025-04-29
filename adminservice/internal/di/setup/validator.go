package setup

import (
	"adminservice/internal/di"

	"github.com/go-playground/validator/v10"
)

func mustValiadtor() di.ValidatorType {
	return validator.New(validator.WithRequiredStructEnabled())
}
