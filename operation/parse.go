package operation

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func Parse(operationRule string) (op Operation, err error) {
	op.baseOperation.rule = operationRule

	operationRule = strings.Trim(operationRule, ";")
	rules := strings.Split(operationRule, ";")

	for _, v := range rules {
		rule := strings.Split(v, ":")
		if len(rule) == 0 || len(rule) > 2 {
			err = errors.New(fmt.Sprintf("Incorrect number of parameters, rule: %s", v))
			return
		}

		name := rule[0]
		params := ""
		if len(rule) > 1 {
			params = rule[1]
		}

		// set image target format
		if name == "format" || name == "f" {
			op.baseOperation.format = params

			continue
		}

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
				key:   strings.ToLower(strings.TrimPrefix(name, "compute-")),
				value: params,
			})
			continue
		}

		// default append to image handler
		op.handlerFunc = append(op.handlerFunc, keyValue{
			key:   name,
			value: params,
		})
	}

	return
}
