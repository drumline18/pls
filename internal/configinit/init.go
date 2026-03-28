package configinit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"pls/internal/config"
	"pls/internal/types"
)

const (
	defaultOpenAIModel = "gpt-4.1-mini"
	defaultOpenAIHost  = "https://api.openai.com"
	defaultOllamaModel = "qwen2.5-coder:7b-instruct-q4_K_M"
	defaultOllamaHost  = "http://127.0.0.1:11434"
)

type providerOption struct {
	Name           string
	Aliases        []string
	Description    string
	ModelLabel     string
	HostLabel      string
	DefaultModel   string
	DefaultHost    string
	CredentialEnvs []string
}

func Run(flags types.Flags) error {
	if !isInteractiveTerminal() {
		return fmt.Errorf("pls config init requires an interactive terminal")
	}
	return run(flags, os.Stdin, os.Stdout)
}

func run(flags types.Flags, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	path, err := config.ResolvePath(flags.ConfigPath)
	if err != nil {
		return err
	}

	existing, err := config.ReadFile(path, false)
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, "pls config init")
	fmt.Fprintln(writer)
	fmt.Fprintf(writer, "Global config path: %s\n", path)
	if hasConfig(existing) {
		fmt.Fprintln(writer, "Existing global config found. Press Enter to keep current values where available.")
	} else {
		fmt.Fprintln(writer, "No global config found yet. Let's set one up.")
	}
	fmt.Fprintln(writer)
	writer.Flush()

	providerDefault := defaultProviderChoice(existing)
	provider, err := askProviderChoice(reader, writer, providerDefault)
	if err != nil {
		return err
	}
	spec := providerInfo(provider)

	final := config.FileConfig{Provider: provider}
	if existing.OpenAIAPIKey != "" {
		final.OpenAIAPIKey = existing.OpenAIAPIKey
	}

	host, err := askProviderHost(reader, writer, spec, existing, "")
	if err != nil {
		return err
	}
	final.Host = host

	model, err := askProviderModel(reader, writer, spec, existing, host, "")
	if err != nil {
		return err
	}
	final.Model = model

	switch provider {
	case "openai":
		key, err := askAPIKey(reader, writer, existing.OpenAIAPIKey)
		if err != nil {
			return err
		}
		final.OpenAIAPIKey = key
	default:
		printCredentialGuidance(writer, spec, existing.OpenAIAPIKey)
	}

	yoloDefault := false
	if existing.YoloMode != nil {
		yoloDefault = *existing.YoloMode
	}
	yolo, err := askYesNo(reader, writer, "Enable yolo mode for safe commands", yoloDefault)
	if err != nil {
		return err
	}
	final.YoloMode = boolPtr(yolo)

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Summary:")
	fmt.Fprintf(writer, "  provider: %s\n", final.Provider)
	fmt.Fprintf(writer, "  model: %s\n", final.Model)
	if final.Host != "" {
		fmt.Fprintf(writer, "  host: %s\n", final.Host)
	}
	fmt.Fprintf(writer, "  yoloMode: %t\n", yolo)
	printCredentialSummary(writer, spec, final.OpenAIAPIKey)
	writer.Flush()

	confirm, err := askYesNo(reader, writer, fmt.Sprintf("Write global config to %s", path), true)
	if err != nil {
		return err
	}
	if !confirm {
		fmt.Fprintln(writer, "Cancelled.")
		return nil
	}

	if err := config.WriteFile(path, final); err != nil {
		return err
	}

	fmt.Fprintln(writer)
	fmt.Fprintf(writer, "Wrote config to %s\n", path)
	fmt.Fprintln(writer, "You can now run: pls doctor")
	fmt.Fprintln(writer, "For a project-specific override, create a local pls.json in that project.")
	return nil
}

func askProviderChoice(reader *bufio.Reader, writer *bufio.Writer, defaultValue string) (string, error) {
	options := supportedProviders()
	fmt.Fprintln(writer, "Available providers:")
	for index, option := range options {
		fmt.Fprintf(writer, "  %d) %s", index+1, option.Name)
		if option.Description != "" {
			fmt.Fprintf(writer, " — %s", option.Description)
		}
		fmt.Fprintln(writer)
	}
	writer.Flush()

	for {
		fmt.Fprintf(writer, "Provider (number or name, default: %s): ", defaultValue)
		writer.Flush()
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}
		value := strings.TrimSpace(line)
		if value == "" {
			return defaultValue, nil
		}
		if index, convErr := strconv.Atoi(value); convErr == nil {
			if index >= 1 && index <= len(options) {
				return options[index-1].Name, nil
			}
		}
		normalized := normalizeProvider(value)
		for _, option := range options {
			if option.Name == normalized {
				return option.Name, nil
			}
		}
		fmt.Fprintln(writer, "Please choose a valid provider number or name.")
	}
}

func askChoice(reader *bufio.Reader, writer *bufio.Writer, label string, choices []string, defaultValue string) (string, error) {
	choiceSet := map[string]struct{}{}
	for _, choice := range choices {
		choiceSet[choice] = struct{}{}
	}

	for {
		fmt.Fprintf(writer, "%s [%s] (default: %s): ", label, strings.Join(choices, "/"), defaultValue)
		writer.Flush()
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}
		value := strings.ToLower(strings.TrimSpace(line))
		if value == "" {
			return defaultValue, nil
		}
		if _, ok := choiceSet[value]; ok {
			return value, nil
		}
		fmt.Fprintf(writer, "Please choose one of: %s\n", strings.Join(choices, ", "))
	}
}

func askLine(reader *bufio.Reader, writer *bufio.Writer, label, defaultValue string) (string, error) {
	fmt.Fprintf(writer, "%s (default: %s): ", label, defaultValue)
	writer.Flush()
	line, err := readLine(reader)
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(line)
	if value == "" {
		return defaultValue, nil
	}
	return value, nil
}

func askRequiredLine(reader *bufio.Reader, writer *bufio.Writer, label string) (string, error) {
	for {
		fmt.Fprintf(writer, "%s: ", label)
		writer.Flush()
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}
		value := strings.TrimSpace(line)
		if value != "" {
			return value, nil
		}
		fmt.Fprintf(writer, "%s is required.\n", label)
	}
}

func askYesNo(reader *bufio.Reader, writer *bufio.Writer, label string, defaultValue bool) (bool, error) {
	defaultLabel := "y/N"
	if defaultValue {
		defaultLabel = "Y/n"
	}

	for {
		fmt.Fprintf(writer, "%s [%s]: ", label, defaultLabel)
		writer.Flush()
		line, err := readLine(reader)
		if err != nil {
			return false, err
		}
		value := strings.ToLower(strings.TrimSpace(line))
		if value == "" {
			return defaultValue, nil
		}
		switch value {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Fprintln(writer, "Please answer y or n.")
		}
	}
}

func askAPIKey(reader *bufio.Reader, writer *bufio.Writer, existing string) (string, error) {
	if existing != "" {
		fmt.Fprintln(writer, "OpenAI API key: press Enter to keep the existing key, type '-' to clear it, or paste a new key.")
	} else {
		fmt.Fprintln(writer, "OpenAI API key: leave blank to use OPENAI_API_KEY / PLS_OPENAI_API_KEY from the environment.")
	}
	writer.Flush()
	line, err := readLine(reader)
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(line)
	switch {
	case value == "" && existing != "":
		return existing, nil
	case value == "-":
		return "", nil
	default:
		return value, nil
	}
}

func askProviderHost(reader *bufio.Reader, writer *bufio.Writer, spec providerOption, existing config.FileConfig, currentHost string) (string, error) {
	hostDefault := firstNonEmpty(existing.Host, currentHost, spec.DefaultHost)
	if hostDefault == "" {
		return "", nil
	}
	label := firstNonEmpty(spec.HostLabel, "Base URL")
	return askLine(reader, writer, label, hostDefault)
}

func askProviderModel(reader *bufio.Reader, writer *bufio.Writer, spec providerOption, existing config.FileConfig, host, currentModel string) (string, error) {
	modelDefault := firstNonEmpty(existing.Model, currentModel, spec.DefaultModel)
	if spec.Name == "ollama" {
		models, modelFetchErr := fetchOllamaModels(firstNonEmpty(host, spec.DefaultHost))
		if modelFetchErr == nil && len(models) > 0 {
			fmt.Fprintln(writer)
			fmt.Fprintln(writer, "Detected Ollama models:")
			for _, model := range models {
				fmt.Fprintf(writer, "  - %s\n", model)
			}
			writer.Flush()
			if existing.Model == "" && strings.TrimSpace(currentModel) == "" {
				modelDefault = models[0]
			}
		} else if modelFetchErr != nil {
			fmt.Fprintln(writer)
			fmt.Fprintf(writer, "Could not list Ollama models from %s: %v\n", firstNonEmpty(host, spec.DefaultHost), modelFetchErr)
			fmt.Fprintln(writer, "You can still enter a model manually.")
			writer.Flush()
		}
	}

	label := firstNonEmpty(spec.ModelLabel, "Model")
	if strings.TrimSpace(modelDefault) == "" {
		return askRequiredLine(reader, writer, label)
	}
	return askLine(reader, writer, label, modelDefault)
}

func fetchOllamaModels(host string) ([]string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(strings.TrimRight(host, "/") + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d", response.StatusCode)
	}

	var payload struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	models := make([]string, 0, len(payload.Models))
	for _, model := range payload.Models {
		if strings.TrimSpace(model.Name) != "" {
			models = append(models, model.Name)
		}
	}
	sort.Strings(models)
	return models, nil
}

func defaultProviderChoice(existing config.FileConfig) string {
	provider := normalizeProvider(existing.Provider)
	if provider != "" {
		return provider
	}

	ollamaHost := firstNonEmpty(existing.Host, os.Getenv("PLS_OLLAMA_HOST"), os.Getenv("OLLAMA_HOST"), defaultOllamaHost)
	if _, err := fetchOllamaModels(ollamaHost); err == nil {
		return "ollama"
	}

	for _, option := range supportedProviders() {
		if option.Name == "ollama" {
			continue
		}
		if providerCredentialConfigured(option.Name, existing.OpenAIAPIKey) {
			return option.Name
		}
	}
	return "ollama"
}

func hasConfig(cfg config.FileConfig) bool {
	return cfg.Provider != "" || cfg.Model != "" || cfg.Host != "" || cfg.OpenAIAPIKey != "" || cfg.YoloMode != nil
}

func isInteractiveTerminal() bool {
	stdinInfo, err := os.Stdin.Stat()
	if err != nil || (stdinInfo.Mode()&os.ModeCharDevice) == 0 {
		return false
	}
	stdoutInfo, err := os.Stdout.Stat()
	if err != nil || (stdoutInfo.Mode()&os.ModeCharDevice) == 0 {
		return false
	}
	return true
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	if err == io.EOF && line == "" {
		return "", io.EOF
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func boolPtr(value bool) *bool {
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func supportedProviders() []providerOption {
	return []providerOption{
		{
			Name:         "ollama",
			Description:  "local or LAN Ollama server",
			ModelLabel:   "Ollama model",
			HostLabel:    "Ollama host",
			DefaultModel: defaultOllamaModel,
			DefaultHost:  defaultOllamaHost,
		},
		{
			Name:           "openai",
			Description:    "OpenAI or OpenAI-compatible API",
			ModelLabel:     "OpenAI model",
			HostLabel:      "OpenAI-compatible base URL",
			DefaultModel:   defaultOpenAIModel,
			DefaultHost:    defaultOpenAIHost,
			CredentialEnvs: []string{"OPENAI_API_KEY", "PLS_OPENAI_API_KEY"},
		},
		{
			Name:           "anthropic",
			Description:    "Anthropic Claude API",
			ModelLabel:     "Anthropic model",
			DefaultModel:   "claude-3-5-haiku-latest",
			CredentialEnvs: []string{"ANTHROPIC_API_KEY"},
		},
		{
			Name:           "gemini",
			Description:    "Google Gemini API",
			ModelLabel:     "Gemini model",
			DefaultModel:   "gemini-2.5-flash",
			CredentialEnvs: []string{"GEMINI_API_KEY"},
		},
		{
			Name:           "groq",
			Description:    "Groq OpenAI-compatible API",
			ModelLabel:     "Groq model",
			HostLabel:      "Groq base URL",
			DefaultModel:   "llama-3.3-70b-versatile",
			DefaultHost:    "https://api.groq.com/openai/v1",
			CredentialEnvs: []string{"GROQ_API_KEY"},
		},
		{
			Name:           "deepseek",
			Description:    "DeepSeek API",
			ModelLabel:     "DeepSeek model",
			HostLabel:      "DeepSeek base URL",
			DefaultModel:   "deepseek-chat",
			DefaultHost:    "https://api.deepseek.com",
			CredentialEnvs: []string{"DEEPSEEK_API_KEY"},
		},
		{
			Name:           "mistral",
			Description:    "Mistral API",
			ModelLabel:     "Mistral model",
			HostLabel:      "Mistral base URL",
			DefaultModel:   "mistral-small-latest",
			DefaultHost:    "https://api.mistral.ai/v1",
			CredentialEnvs: []string{"MISTRAL_API_KEY"},
		},
		{
			Name:           "zai",
			Aliases:        []string{"z.ai"},
			Description:    "Z.ai API",
			ModelLabel:     "Z.ai model",
			HostLabel:      "Z.ai base URL",
			DefaultModel:   "glm-4.5-air",
			DefaultHost:    "https://api.z.ai/api/paas/v4",
			CredentialEnvs: []string{"ZAI_API_KEY"},
		},
		{
			Name:           "llamacpp",
			Aliases:        []string{"llama.cpp", "llama-cpp"},
			Description:    "llama.cpp OpenAI-compatible server",
			ModelLabel:     "llama.cpp model",
			HostLabel:      "llama.cpp base URL",
			DefaultHost:    "http://127.0.0.1:8080/v1",
		},
		{
			Name:        "llamafile",
			Description: "llamafile OpenAI-compatible server",
			ModelLabel:  "llamafile model",
			HostLabel:   "llamafile base URL",
			DefaultHost: "http://127.0.0.1:8080/v1",
		},
	}
}

func providerInfo(provider string) providerOption {
	normalized := normalizeProvider(provider)
	for _, option := range supportedProviders() {
		if option.Name == normalized {
			return option
		}
	}
	return providerOption{Name: normalized}
}

func normalizeProvider(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	for _, option := range supportedProviders() {
		if normalized == option.Name {
			return option.Name
		}
		for _, alias := range option.Aliases {
			if normalized == strings.ToLower(strings.TrimSpace(alias)) {
				return option.Name
			}
		}
	}
	return normalized
}

func providerCredentialConfigured(provider string, storedOpenAIKey string) bool {
	spec := providerInfo(provider)
	if spec.Name == "openai" && strings.TrimSpace(storedOpenAIKey) != "" {
		return true
	}
	for _, envName := range spec.CredentialEnvs {
		if strings.TrimSpace(os.Getenv(envName)) != "" {
			return true
		}
	}
	return false
}

func printCredentialGuidance(writer *bufio.Writer, spec providerOption, storedOpenAIKey string) {
	if len(spec.CredentialEnvs) == 0 {
		return
	}
	fmt.Fprintln(writer)
	if providerCredentialConfigured(spec.Name, storedOpenAIKey) {
		fmt.Fprintf(writer, "Credentials detected from environment for %s (%s).\n", spec.Name, strings.Join(spec.CredentialEnvs, " / "))
	} else {
		fmt.Fprintf(writer, "Credentials are not stored in pls config for %s. Set %s in your environment.\n", spec.Name, strings.Join(spec.CredentialEnvs, " or "))
	}
	writer.Flush()
}

func printCredentialSummary(writer *bufio.Writer, spec providerOption, storedOpenAIKey string) {
	switch spec.Name {
	case "openai":
		if strings.TrimSpace(storedOpenAIKey) != "" {
			fmt.Fprintln(writer, "  openaiApiKey: [stored in config]")
		} else {
			fmt.Fprintln(writer, "  openaiApiKey: [not stored; use environment if needed]")
		}
	default:
		if len(spec.CredentialEnvs) > 0 {
			if providerCredentialConfigured(spec.Name, storedOpenAIKey) {
				fmt.Fprintf(writer, "  credentials: environment (%s)\n", strings.Join(spec.CredentialEnvs, " / "))
			} else {
				fmt.Fprintf(writer, "  credentials: set %s in your environment\n", strings.Join(spec.CredentialEnvs, " or "))
			}
		}
	}
}
