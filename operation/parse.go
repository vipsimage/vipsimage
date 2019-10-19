package operation

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Parse operation rule
func Parse(operationRule string) (op Operation, err error) {
	op.baseOperation.rule = operationRule

	operationRule = strings.Trim(operationRule, ";")
	rules := strings.Split(operationRule, ";")

	for _, v := range rules {
		rule := strings.Split(v, ":")
		if len(rule) == 0 || len(rule) > 2 {
			err = fmt.Errorf("incorrect number of parameters, rule: %s", v)
			return
		}

		name := rule[0]
		paramsStr := ""
		params := url.Values{}
		if len(rule) > 1 {
			paramsStr = rule[1]

			params, err = url.ParseQuery(paramsStr)
			if err != nil {
				return
			}
		}

		// set image target format
		if name == "format" || name == "f" {
			op.baseOperation.format = paramsStr

			continue
		}

		// set store setting
		if name == "store" {
			op.baseOperation.storage, err = parseStorage(params)
			if err != nil {
				err = errors.WithMessage(err, "storage parse error")
				return
			}
			continue
		}

		if strings.HasPrefix(name, "compute-") {
			op.computeFunc = append(op.computeFunc, keyValue{
				key:    strings.ToLower(strings.TrimPrefix(name, "compute-")),
				value:  paramsStr,
				Values: params,
			})
			continue
		}

		// default append to image handler
		op.handlerFunc = append(op.handlerFunc, keyValue{
			key:    name,
			value:  paramsStr,
			Values: params,
		})
	}

	return
}
