package resp

import (
	"strconv"
)

var (
	endTerm = []byte{'\r', '\n'}
)

func buildInteger(n int64) []byte {
	retval := []byte{}
	retval = append(retval, ':')
	retval = strconv.AppendInt(retval, n, 10)
	retval = append(retval, endTerm...)
	return retval
}

func buildString(s string) []byte {
	retval := []byte{}
	retval = append(retval, []byte(s)...)
	retval = append(retval, endTerm...)
	return retval
}

func buildError(err error) []byte {
	retval := []byte{}
	retval = append(retval, '-')
	if err == nil {
		retval = append(retval, []byte("error is nil")...)
	} else {
		retval = append(retval, []byte(err.Error())...)
	}
	retval = append(retval, endTerm...)
	return retval
}

func buildBulkString(s string) []byte {
	retval := []byte{}
	retval = append(retval, '$')
	retval = strconv.AppendInt(retval, int64(len(s)), 10)
	retval = append(retval, endTerm...)
	retval = append(retval, buildString(s)...)
	return retval
}

func BuildArray(a []interface{}) []byte {
	retval := []byte{}
	retval = append(retval, '*')
	if a == nil {
		retval = append(retval, []byte("-1")...)
		retval = append(retval, endTerm...)
		return retval
	}

	retval = strconv.AppendInt(retval, int64(len(a)), 10)
	retval = append(retval, endTerm...)
	for i := 0; i < len(a); i++ {
		switch v := a[i].(type) {
		case []interface{}:
			retval = append(retval, BuildArray(a)...)
		case nil:
			retval = append(retval, BuildArray(nil)...)
		case int64:
			retval = append(retval, buildInteger(v)...)
		case string:
			retval = append(retval, buildString(v)...)
		case error:
			retval = append(retval, buildError(v)...)
		}
	}
	return retval
}

func buildBulk(a []byte) []byte {
	retval := []byte{}
	retval = append(retval, '$')
	if a == nil {
		retval = append(retval, []byte("-1")...)
	} else {
		retval = append(retval, a...)
	}

	retval = append(retval, endTerm...)
	return retval
}

func BuildCommand(cmd string, args ...interface{}) []byte {
	retval := []byte{}
	retval = append(retval, '*')
	retval = strconv.AppendInt(retval, int64(1+len(args)), 10)
	retval = append(retval, endTerm...)
	retval = append(retval, buildBulkString(cmd)...)
	for _, arg := range args {
		switch a := arg.(type) {
		case string:
			retval = append(retval, buildBulkString(a)...)
		case []byte:
			retval = append(retval, buildBulk(a)...)
		}
	}

	return retval
}
