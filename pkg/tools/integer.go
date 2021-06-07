package tools

import (
	"github.com/shoriwe/gplasma/pkg/errors"
	"strconv"
	"strings"
)

func ParseInteger(s string) (int64, *errors.Error) {
	s = strings.ReplaceAll(s, "_", "")
	var number int64
	var parsingError error
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		number, parsingError = strconv.ParseInt(s[2:], 16, 64)
	} else if strings.HasPrefix(s, "0o") || strings.HasPrefix(s, "0O") {
		number, parsingError = strconv.ParseInt(s[2:], 8, 64)
	} else if strings.HasPrefix(s, "0B") || strings.HasPrefix(s, "0B") {
		number, parsingError = strconv.ParseInt(s[2:], 2, 64)
	} else {
		number, parsingError = strconv.ParseInt(s, 10, 64)
	}
	if parsingError != nil {
		return 0, errors.New(errors.UnknownLine, parsingError.Error(), errors.GoRuntimeError)
	}
	return number, nil
}
