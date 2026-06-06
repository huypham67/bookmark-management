package dbutils

import (
	"errors"
	"strings"
)

var errorFilter = []func(error) (bool, error){
	IsDuplicationError,
	IsForeignKeyViolationError,
	IsRecordNotFoundError,
}

var (
	ErrDuplicationType         = errors.New("duplicate type error")
	ErrRecordNotFoundType      = errors.New("record not found")
	ErrForeignKeyViolationType = errors.New("foreign key violation")
)

func ClassifyError(err error) error {
	for _, filter := range errorFilter {
		if isMatch, errType := filter(err); isMatch {
			return errType
		}
	}
	return err
}

func IsDuplicationError(err error) (bool, error) {
	errMsg := err.Error()
	return strings.Contains(errMsg, "duplicate key value violates unique constraint") ||
		strings.Contains(errMsg, "UNIQUE constraint failed"), ErrDuplicationType
}

func IsForeignKeyViolationError(err error) (bool, error) {
	errMsg := err.Error()
	return strings.Contains(errMsg, "violates foreign key constraint") ||
		strings.Contains(errMsg, "FOREIGN KEY constraint failed"), ErrForeignKeyViolationType
}

func IsRecordNotFoundError(err error) (bool, error) {
	return strings.Contains(err.Error(), "record not found"), ErrRecordNotFoundType
}
