package helper

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func Float64ToPGNumeric(src float64) (pgtype.Numeric, error) {
	pgNumeric := pgtype.Numeric{}
	err := pgNumeric.Scan(fmt.Sprintf("%f", src))
	return pgNumeric, err
}

func PGNumericToFloat64(n pgtype.Numeric) (float64, error) {
	f, err := n.Float64Value()
	if err != nil {
		return 0, err
	}
	return f.Float64, nil
}

