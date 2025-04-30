package sqlcutils

import (
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

// Converts pgtype.Numeric to float64, returns 0.00 if Balance is NULL in DB
func NumericToFloat64(entity pgtype.Numeric) float64 {
	var res float64
	if !entity.Valid { // if field is NULL in DB
		return res
	}
	res, _ = entity.Int.Float64()

	return res / math.Pow(float64(10), -float64(entity.Exp))
}
