package doctor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"pls/internal/config"
	runtimeinfo "pls/internal/runtimeinfo"
	"pls/internal/types"
)

type Report struct {
	Joke                 string `json:"joke"`
	OverallStatus        string `json:"overallStatus"`
	ConfigPath           string `json:"configPath"`
	ConfigExists         bool   `json:"configExists"`
	ConfigError          string `json:"configError,omitempty"`
	Provider             string `json:"provider,omitempty"`
	Model                string `json:"model,omitempty"`
	Host                 string `json:"host,omitempty"`
	Shell                string `json:"shell,omitempty"`
	OS                   string `json:"os,omitempty"`
	CWD                  string `json:"cwd,omitempty"`
	ExecutablePath       string `json:"executablePath,omitempty"`
	InstalledCommandPath string `json:"installedCommandPath,omitempty"`
	InPath               bool   `json:"inPath"`
	ProviderStatus       string `json:"providerStatus,omitempty"`
	ProviderMessage      string `json:"providerMessage,omitempty"`
}

func Run(ctx context.Context, flags types.Flags) (Report, error) {
	report := Report{
		Joke: "Why did the CLI go to the doctor? It had too many terminal conditions.",
	}

	configPath, err := config.ResolvePath(flags.ConfigPath)
	if err != nil {
		return report, err
	}
	report.ConfigPath = configPath

	if info, statErr := os.Stat(configPath); statErr == nil && !info.IsDir() {
		report.ConfigExists = true
	} else if statErr != nil && !os.IsNotExist(statErr) {
		report.ConfigError = statErr.Error()
	}

	runtimeContext, err := runtimeinfo.Get(flags.Shell, flags.OS)
	if err == nil {
		report.Shell = runtimeContext.Shell
		report.OS = runtimeContext.OS
		report.CWD = runtimeContext.CWD
	}

	if exePath, exeErr := os.Executable(); exeErr == nil {
		report.ExecutablePath = exePath
	}
	if installedPath, lookErr := exec.LookPath("pls"); lookErr == nil {
		report.InstalledCommandPath = installedPath
		report.InPath = true
	}

	cfg, loadErr := config.Load(flags)
	if loadErr != nil {
		report.ConfigError = loadErr.Error()
		report.OverallStatus = "warn"
		return report, nil
	}

	report.Provider = cfg.Provider
	report.Model = cfg.Model
	report.Host = cfg.Host

	status, message := checkProvider(ctx, cfg)
	report.ProviderStatus = status
	report.ProviderMessage = message

	report.OverallStatus = overallStatus(report)
	return report, nil
}

func Human(report Report) string {
	lines := []string{
		report.Joke,
		"",
		"pls doctor",
		fmt.Sprintf("  overall: %s", report.OverallStatus),
		"",
		"Config:",
		fmt.Sprintf("  path: %s", report.ConfigPath),
		fmt.Sprintf("  exists: %s", yesNo(report.ConfigExists)),
	}

	if report.ConfigError != "" {
		lines = append(lines, fmt.Sprintf("  error: %s", report.ConfigError))
	}

	lines = append(lines,
		"",
		"Runtime:",
		fmt.Sprintf("  os: %s", emptyFallback(report.OS, "unknown")),
		fmt.Sprintf("  shell: %s", emptyFallback(report.Shell, "unknown")),
		fmt.Sprintf("  cwd: %s", emptyFallback(report.CWD, "unknown")),
		"",
		"Install:",
		fmt.Sprintf("  executable: %s", emptyFallback(report.ExecutablePath, "unknown")),
		fmt.Sprintf("  command in PATH: %s", yesNo(report.InPath)),
	)

	if report.InstalledCommandPath != "" {
		lines = append(lines, fmt.Sprintf("  command path: %s", report.InstalledCommandPath))
	}

	lines = append(lines,
		"",
		"Provider:",
		fmt.Sprintf("  provider: %s", emptyFallback(report.Provider, "unknown")),
		fmt.Sprintf("  model: %s", emptyFallback(report.Model, "unknown")),
		fmt.Sprintf("  host: %s", emptyFallback(report.Host, "unknown")),
		fmt.Sprintf("  status: %s", emptyFallback(report.ProviderStatus, "unknown")),
	)

	if report.ProviderMessage != "" {
		lines = append(lines, fmt.Sprintf("  message: %s", report.ProviderMessage))
	}

	return strings.Join(lines, "\n")
}

func checkProvider(ctx context.Context, cfg types.Config) (string, string) {
	switch cfg.Provider {
	case "ollama":
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(cfg.Host, "/")+"/api/tags", nil)
		if err != nil {
			return "error", err.Error()
		}
		response, err := httpClient().Do(request)
		if err != nil {
			return "error", err.Error()
		}
		defer response.Body.Close()
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			return "ok", "connected to Ollama successfully"
		}
		return "error", fmt.Sprintf("Ollama returned HTTP %d", response.StatusCode)
	case "openai":
		if cfg.OpenAIAPIKey == "" {
			return "warn", "OpenAI API key is not configured"
		}
		return "ok", "OpenAI API key is configured; network probe skipped"
	default:
		return "warn", "provider health check is not implemented for this provider"
	}
}

func overallStatus(report Report) string {
	if report.ConfigError != "" {
		return "warn"
	}
	if report.ProviderStatus == "error" {
		return "warn"
	}
	return "ok"
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

func httpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

