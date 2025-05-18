package dto

import (
	"encoding/json"
)

type ZipCode struct {
	ZipCode string `json:"zipcode"`
}

func (z *ZipCode) UnmarshalJSON(data []byte) error {
	var temp map[string]string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	z.ZipCode = ""
	if val, ok := temp["zipcode"]; ok {
		z.ZipCode = val
	} else if val, ok := temp["cep"]; ok {
		z.ZipCode = val
	}

	return nil
}
