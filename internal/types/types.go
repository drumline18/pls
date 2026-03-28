package types

import (
	"errors"
	"fmt"
	"strings"
)

type Flags struct {
	Provider        string
	Model           string
	Shell           string
	OS              string
	Host            string
	ConfigPath      string
	JSON            bool
	Help            bool
	PrintConfigPath bool
}

type ParsedArgs struct {
	Flags        Flags
	RequestParts []string
	Help         bool
}

type Config struct {
	Provider     string
	Model        string
	Host         string
	ConfigPath   string
	OutputJSON   bool
	OpenAIAPIKey string
}

type RuntimeContext struct {
	CWD           string `json:"cwd"`
	OS            string `json:"os"`
	Shell         string `json:"shell"`
	HomeDirectory string `json:"homeDirectory"`
	IsTTY         bool   `json:"isTTY"`
}

type Messages struct {
	System string
	User   map[string]any
}

type Suggestion struct {
	Command               string `json:"command"`
	Explanation           string `json:"explanation"`
	Risk                  string `json:"risk"`
	RequiresConfirmation  bool   `json:"requiresConfirmation"`
	NeedsClarification    bool   `json:"needsClarification"`
	ClarificationQuestion string `json:"clarificationQuestion"`
	Notes                 string `json:"notes"`
	Platform              string `json:"platform"`
	Refused               bool   `json:"refused"`
}

func ValidateSuggestion(s Suggestion) (Suggestion, error) {
	s.Command = strings.TrimSpace(s.Command)
	s.Explanation = strings.TrimSpace(s.Explanation)
	s.Risk = strings.ToLower(strings.TrimSpace(s.Risk))
	s.ClarificationQuestion = strings.TrimSpace(s.ClarificationQuestion)
	s.Notes = strings.TrimSpace(s.Notes)
	s.Platform = strings.TrimSpace(s.Platform)

	if s.Explanation == "" {
		return Suggestion{}, errors.New("explanation must be a non-empty string")
	}

	switch s.Risk {
	case "low", "medium", "high", "critical":
	default:
		return Suggestion{}, fmt.Errorf("invalid risk level: %s", s.Risk)
	}

	if s.NeedsClarification {
		if s.ClarificationQuestion == "" {
			return Suggestion{}, errors.New("clarificationQuestion is required when needsClarification is true")
		}
		return s, nil
	}

	if !s.Refused && s.Command == "" {
		return Suggestion{}, errors.New("command must be a non-empty string unless the request is refused or needs clarification")
	}

	return s, nil
}
