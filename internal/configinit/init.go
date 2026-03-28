package configinit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
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
	provider, err := askChoice(reader, writer, "Provider", []string{"ollama", "openai"}, providerDefault)
	if err != nil {
		return err
	}

	final := config.FileConfig{Provider: provider}
	if existing.OpenAIAPIKey != "" {
		final.OpenAIAPIKey = existing.OpenAIAPIKey
	}

	switch provider {
	case "ollama":
		hostDefault := firstNonEmpty(existing.Host, os.Getenv("PLS_OLLAMA_HOST"), os.Getenv("OLLAMA_HOST"), defaultOllamaHost)
		host, err := askLine(reader, writer, "Ollama host", hostDefault)
		if err != nil {
			return err
		}
		final.Host = host

		models, modelFetchErr := fetchOllamaModels(host)
		if modelFetchErr == nil && len(models) > 0 {
			fmt.Fprintln(writer)
			fmt.Fprintln(writer, "Detected Ollama models:")
			for _, model := range models {
				fmt.Fprintf(writer, "  - %s\n", model)
			}
			writer.Flush()
		} else if modelFetchErr != nil {
			fmt.Fprintln(writer)
			fmt.Fprintf(writer, "Could not list Ollama models from %s: %v\n", host, modelFetchErr)
			fmt.Fprintln(writer, "You can still enter a model manually.")
			writer.Flush()
		}

		modelDefault := defaultOllamaModel
		if existing.Provider == "ollama" && existing.Model != "" {
			modelDefault = existing.Model
		} else if len(models) > 0 {
			modelDefault = models[0]
		}
		model, err := askLine(reader, writer, "Ollama model", modelDefault)
		if err != nil {
			return err
		}
		final.Model = model

	case "openai":
		baseURLDefault := defaultOpenAIHost
		if existing.Provider == "openai" && existing.Host != "" {
			baseURLDefault = existing.Host
		}
		baseURL, err := askLine(reader, writer, "OpenAI-compatible base URL", baseURLDefault)
		if err != nil {
			return err
		}
		if strings.TrimSpace(baseURL) != "" && strings.TrimSpace(baseURL) != defaultOpenAIHost {
			final.Host = baseURL
		}

		modelDefault := defaultOpenAIModel
		if existing.Provider == "openai" && existing.Model != "" {
			modelDefault = existing.Model
		}
		model, err := askLine(reader, writer, "OpenAI model", modelDefault)
		if err != nil {
			return err
		}
		final.Model = model

		key, err := askAPIKey(reader, writer, existing.OpenAIAPIKey)
		if err != nil {
			return err
		}
		final.OpenAIAPIKey = key
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
	if final.Provider == "openai" {
		if final.OpenAIAPIKey != "" {
			fmt.Fprintln(writer, "  openaiApiKey: [stored in config]")
		} else {
			fmt.Fprintln(writer, "  openaiApiKey: [not stored; use environment if needed]")
		}
	}
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
	if final.Provider == "ollama" {
		fmt.Fprintln(writer, "You can now run: pls doctor")
	} else {
		fmt.Fprintln(writer, "You can now run: pls doctor")
	}
	fmt.Fprintln(writer, "For a project-specific override, create a local pls.json in that project.")
	return nil
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
	if existing.Provider != "" {
		return existing.Provider
	}
	ollamaHost := firstNonEmpty(existing.Host, os.Getenv("PLS_OLLAMA_HOST"), os.Getenv("OLLAMA_HOST"), defaultOllamaHost)
	if _, err := fetchOllamaModels(ollamaHost); err == nil {
		return "ollama"
	}
	if existing.OpenAIAPIKey != "" || os.Getenv("OPENAI_API_KEY") != "" || os.Getenv("PLS_OPENAI_API_KEY") != "" {
		return "openai"
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
