package configshow

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/drumline18/pls/internal/config"
	"github.com/drumline18/pls/internal/types"
)

type Report struct {
	GlobalConfigPath      string `json:"globalConfigPath"`
	GlobalConfigExists    bool   `json:"globalConfigExists"`
	LocalConfigPath       string `json:"localConfigPath,omitempty"`
	LocalConfigExists     bool   `json:"localConfigExists"`
	EffectiveProvider     string `json:"effectiveProvider,omitempty"`
	EffectiveModel        string `json:"effectiveModel,omitempty"`
	EffectiveHost         string `json:"effectiveHost,omitempty"`
	ConfigStoredAPIKeyConfigured bool  `json:"configStoredApiKeyConfigured"`
	YoloMode             bool   `json:"yoloMode"`
	YoloSource           string `json:"yoloSource,omitempty"`
}

func Run(flags types.Flags) (Report, error) {
	globalPath, err := config.ResolvePath(flags.ConfigPath)
	if err != nil {
		return Report{}, err
	}

	report := Report{GlobalConfigPath: globalPath}
	if info, statErr := os.Stat(globalPath); statErr == nil && !info.IsDir() {
		report.GlobalConfigExists = true
	} else if statErr != nil && !os.IsNotExist(statErr) {
		return Report{}, statErr
	}

	cfg, err := config.Load(flags)
	if err != nil {
		return Report{}, err
	}

	report.LocalConfigPath = cfg.LocalConfigPath
	report.LocalConfigExists = strings.TrimSpace(cfg.LocalConfigPath) != ""
	report.EffectiveProvider = cfg.Provider
	report.EffectiveModel = cfg.Model
	report.EffectiveHost = cfg.Host
	report.ConfigStoredAPIKeyConfigured = strings.TrimSpace(cfg.OpenAIAPIKey) != ""
	report.YoloMode = cfg.YoloMode
	report.YoloSource = cfg.YoloSource
	return report, nil
}

func Human(report Report) string {
	lines := []string{
		"pls config show",
		fmt.Sprintf("  global path: %s", report.GlobalConfigPath),
		fmt.Sprintf("  global exists: %s", yesNo(report.GlobalConfigExists)),
		fmt.Sprintf("  local override: %s", emptyFallback(report.LocalConfigPath, "none")),
		fmt.Sprintf("  local exists: %s", yesNo(report.LocalConfigExists)),
		"",
		"Effective:",
		fmt.Sprintf("  provider: %s", emptyFallback(report.EffectiveProvider, "unknown")),
		fmt.Sprintf("  model: %s", emptyFallback(report.EffectiveModel, "unknown")),
		fmt.Sprintf("  host: %s", emptyFallback(report.EffectiveHost, "unknown")),
		fmt.Sprintf("  config-stored api key configured: %s", yesNo(report.ConfigStoredAPIKeyConfigured)),
		fmt.Sprintf("  yolo mode: %s", yesNo(report.YoloMode)),
		fmt.Sprintf("  yolo source: %s", emptyFallback(report.YoloSource, "default")),
	}
	return strings.Join(lines, "\n")
}

func JSON(report Report) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func emptyFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
