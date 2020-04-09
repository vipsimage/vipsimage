package operation

import (
	"crypto/md5" // nolint
	"encoding/hex"
)

// Computed operation
type Computed struct {
	MD5 MD5 `json:"md5"`
}

// HasCompute return true if has computed function
// todo
func (th Computed) HasCompute() bool {
	if th.MD5 {
		return true
	}

	return false
}

// Compute image
func (th Computed) Compute(content []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})

	if th.MD5 {
		result, err := th.MD5.Calculate(content)
		if err != nil {
			return res, err
		}
		res["md5"] = result
	}

	return
}

// Calculator Calculate image
type Calculator interface {
	Calculate(content []byte) (interface{}, error)
}

// MD5 option
type MD5 bool

// Calculate the md5 of the picture
func (MD5) Calculate(content []byte) (interface{}, error) {
	contentMd5 := md5.Sum(content) // nolint
	return hex.EncodeToString(contentMd5[:]), nil
}
