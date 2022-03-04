package apperrors

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var errNoargs = errors.New("apperrors: no args")

func E(ctx context.Context, args ...interface{}) error {
	if len(args) == 0 {
		return errNoargs
	}
	var st *status.Status
	var code codes.Code
	var msg string

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			msg = arg
		case codes.Code:
			code = arg
		case error:
			code = parseError(arg)
			msg = arg.Error()
		}
	}
	st = status.New(code, msg)
	return st.Err()
}

// parse errors check for some specific.
// nolint:govet
func parseError(err error) codes.Code {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return codes.NotFound
	}
	err1, ok := err.(*mysql.MySQLError)
	if !ok {
		return codes.Unknown
	}

	switch err1.Number {
	case 1062:
		return codes.AlreadyExists
	case 1452:
		return codes.AlreadyExists
	default:
		return codes.Unknown
	}

}
