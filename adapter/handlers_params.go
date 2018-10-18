package adapter

// copy from tendermint/rpc/lib/server/http_params.go

import (
	"encoding/hex"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var (
	// Parts of regular expressions
	atom    = "[A-Z0-9!#$%&'*+\\-/=?^_`{|}~]+"
	dotAtom = atom + `(?:\.` + atom + `)*`
	domain  = `[A-Z0-9.-]+\.[A-Z]{2,4}`

	// RegexpInt 整数正则
	RegexpInt = regexp.MustCompile(`^-?[0-9]+$`)
	// RegexpHex 16进制正则
	RegexpHex = regexp.MustCompile(`^(?i)[a-f0-9]+$`)
	// RegexpEmail 正则
	RegexpEmail = regexp.MustCompile(`^(?i)(` + dotAtom + `)@(` + dotAtom + `)$`)
	// RegexpAddress 正则
	RegexpAddress = regexp.MustCompile(`^(?i)[a-z0-9]{25,34}$`)
	// RegexpHost 正则
	RegexpHost = regexp.MustCompile(`^(?i)(` + domain + `)$`)

	//RE_ID12       = regexp.MustCompile(`^[a-zA-Z0-9]{12}$`)
)

// GetParam 获取指定参数
func GetParam(r *http.Request, param string) string {
	s := r.URL.Query().Get(param)
	if s == "" {
		s = r.FormValue(param)
	}
	return s
}

// GetParamByteSlice 获取指定字节切片
func GetParamByteSlice(r *http.Request, param string) ([]byte, error) {
	s := GetParam(r, param)
	return hex.DecodeString(s)
}

// GetParamInt64 获取指定参数
func GetParamInt64(r *http.Request, param string) (int64, error) {
	s := GetParam(r, param)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.Errorf(param, err.Error())
	}
	return i, nil
}

// GetParamInt32 获取指定参数
func GetParamInt32(r *http.Request, param string) (int32, error) {
	s := GetParam(r, param)
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, errors.Errorf(param, err.Error())
	}
	return int32(i), nil
}

// GetParamUint64 获取指定参数
func GetParamUint64(r *http.Request, param string) (uint64, error) {
	s := GetParam(r, param)
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, errors.Errorf(param, err.Error())
	}
	return i, nil
}

// GetParamUint 获取指定参数
func GetParamUint(r *http.Request, param string) (uint, error) {
	s := GetParam(r, param)
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, errors.Errorf(param, err.Error())
	}
	return uint(i), nil
}

// GetParamRegexp 获取指定参数
func GetParamRegexp(r *http.Request, param string, re *regexp.Regexp) (string, error) {
	s := GetParam(r, param)
	if !re.MatchString(s) {
		return "", errors.Errorf(param, "Did not match regular expression %v", re.String())
	}
	return s, nil
}

// GetParamFloat64 获取指定参数
func GetParamFloat64(r *http.Request, param string) (float64, error) {
	s := GetParam(r, param)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Errorf(param, err.Error())
	}
	return f, nil
}
