package operation

import (
	"crypto/md5"
	"encoding/hex"
)

type ComputeFunc func(content []byte, op baseOperation) (interface{}, error)

func GetComputeFunc(name string) (ComputeFunc, bool) {
	switch name {
	case "md5":
		return imageMd5, true
	default:
		return nil, false
	}
}

func imageMd5(content []byte, _ baseOperation) (interface{}, error) {
	s := md5.Sum(content)
	return hex.EncodeToString(s[:]), nil
}
