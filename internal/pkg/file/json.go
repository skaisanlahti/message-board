package file

import (
	"encoding/json"
	"os"
)

func ReadJSON[T any](path string) (T, error) {
	var data T
	bytes, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
