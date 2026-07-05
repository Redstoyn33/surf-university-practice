package repository

import "strings"

func isUniqueConstraintErr(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
