package sqlcutils

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

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

// Converts Numner to pgtype.Numeric
func NumberToNumeric[T int | float64](number T) (pgtype.Numeric, error) {

	strNum := fmt.Sprint(number)

	arr := strings.Split(strNum, ".")

	var exp int
	if len(arr) == 2 {
		exp = -2
	} else {
		exp = -1
	}
	return pgtype.Numeric{Int: big.NewInt(int64(float64(number) * math.Pow(10, float64(-exp)))), Exp: int32(exp), Valid: true}, nil
}

func Int4ToInt(num pgtype.Int4) int {
	return int(num.Int32)
}

func IntToInt4(num int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(num),
		Valid: true,
	}
}

func StringToDate(date string) (pgtype.Date, error) {
	layout := "2006-01-02"
	taskTime, err := time.Parse(layout, date)
	if err != nil {
		return pgtype.Date{}, err
	}

	return pgtype.Date{Time: taskTime, Valid: true}, nil
}
