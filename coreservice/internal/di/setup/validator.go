package setup

import (
	"coreservice/internal/di"

	"github.com/go-playground/validator/v10"
)

func mustValiadtor() di.ValidatorType {
	return validator.New(validator.WithRequiredStructEnabled())
}
