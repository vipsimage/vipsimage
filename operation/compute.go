package operation

import (
	"crypto/md5" // nolint
	"encoding/hex"
)

// ComputeFunc compute function
type ComputeFunc func(content []byte, params keyValue, op baseOperation) (interface{}, error)

// GetComputeFunc return compute function by name.
func GetComputeFunc(name string) (ComputeFunc, bool) {
	switch name {
	case "md5":
		return imageMd5, true
	default:
		return nil, false
	}

}

// imageMd5 hash image
func imageMd5(content []byte, _ keyValue, _ baseOperation) (interface{}, error) {
	contentMd5 := md5.Sum(content) // nolint
	return hex.EncodeToString(contentMd5[:]), nil
}
