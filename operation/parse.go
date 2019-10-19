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

		name, paramsStr, params, err := parseQuery(rule)
		if err != nil {
			return op, errors.WithMessage(err, "parse query error")
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
				return op, errors.WithMessage(err, "storage parse error")
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

func parseQuery(rule []string) (name, paramsStr string, params url.Values, err error) {
	name = rule[0]
	if len(rule) > 1 {
		paramsStr = rule[1]

		params, err = url.ParseQuery(paramsStr)
		if err != nil {
			return
		}
	}

	return
}
