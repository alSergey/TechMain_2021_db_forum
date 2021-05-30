package null

import "database/sql"

func NewNullInt(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewIntFromNull(i sql.NullInt64) int64 {
	if i.Valid {
		return i.Int64
	}

	return 0
}
