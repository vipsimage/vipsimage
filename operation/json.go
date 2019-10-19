package operation

import (
	"encoding/json"
	"net/url"
)

func (th Operation) String() string {
	b, err := json.Marshal(th)
	if err != nil {
		return ""
	}
	return string(b)
}

type opJSON struct {
	Operation string     `json:"operation"`
	ParamsStr string     `json:"params_str"`
	Params    url.Values `json:"params"`
}

func kv2json(kvs ...keyValue) (res []opJSON) {
	for _, v := range kvs {
		res = append(res, opJSON{
			Operation: v.key,
			ParamsStr: v.value,
			Params:    v.Values,
		})
	}

	return
}

type operationJSON struct {
	Rule    string   `json:"rule"`
	Storage opJSON   `json:"storage"`
	Format  string   `json:"format"`
	Handler []opJSON `json:"handler"`
	Compute []opJSON `json:"compute"`
}

// MarshalJSON marshal operation
func (th Operation) MarshalJSON() ([]byte, error) {
	var opJSON operationJSON
	storage := kv2json(th.storage)
	if len(storage) > 0 {
		opJSON.Storage = storage[0]
	}

	opJSON.Rule = th.rule
	opJSON.Format = th.format
	opJSON.Handler = kv2json(th.handlerFunc...)
	opJSON.Compute = kv2json(th.computeFunc...)

	return json.Marshal(opJSON)
}
