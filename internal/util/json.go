package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func MustJSON(value any) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func ExtractJSONObject(text string) (string, error) {
	trimmed := strings.TrimSpace(text)
	if json.Valid([]byte(trimmed)) {
		return trimmed, nil
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start == -1 || end == -1 || end <= start {
		return "", fmt.Errorf("response did not contain a JSON object")
	}

	candidate := trimmed[start : end+1]
	if !json.Valid([]byte(candidate)) {
		return "", fmt.Errorf("response did not contain a valid JSON object")
	}

	return candidate, nil
}
