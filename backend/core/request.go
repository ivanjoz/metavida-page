package core

import (
	"encoding/json"
	"errors"
	"strings"
)

func DecodeJSONBody(body *string, target any) error {
	if body == nil {
		return errors.New("missing request body")
	}
	return json.NewDecoder(strings.NewReader(*body)).Decode(target)
}
