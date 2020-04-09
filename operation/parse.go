package operation

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Rule is operation rule describe
type Rule struct {
	Storage    `json:"storage,omitempty"`
	*Format    `json:"format,omitempty"`
	*Operation `json:"operation,omitempty"`
	*Computed  `json:"computed,omitempty"`
}

// Parse operation rule
func Parse(rule string) (r Rule, err error) {
	decode := json.NewDecoder(strings.NewReader(rule))
	decode.UseNumber()

	err = decode.Decode(&r)
	if err != nil {
		return
	}

	err = validate.Struct(r)
	if err != nil {
		return
	}

	return
}
