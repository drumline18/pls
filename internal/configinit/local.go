package configinit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/drumline18/pls/internal/config"
	"github.com/drumline18/pls/internal/types"
)

func RunLocal(flags types.Flags) error {
	if !isInteractiveTerminal() {
		return fmt.Errorf("pls config local init requires an interactive terminal")
	}
	return runLocal(flags, os.Stdin, os.Stdout)
}

func runLocal(flags types.Flags, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(cwd, "pls.json")

	existing, err := config.ReadFile(path, false)
	if err != nil {
		return err
	}
	current, err := config.Load(flags)
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, "pls config local init")
	fmt.Fprintln(writer)
	fmt.Fprintf(writer, "Local config path: %s\n", path)
	fmt.Fprintln(writer, "This writes project-local overrides only. Unset fields will continue to inherit from environment/global config.")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Current effective values:")
	fmt.Fprintf(writer, "  provider: %s\n", current.Provider)
	fmt.Fprintf(writer, "  model: %s\n", current.Model)
	fmt.Fprintf(writer, "  host: %s\n", current.Host)
	fmt.Fprintf(writer, "  yolo mode: %t (%s)\n", current.YoloMode, current.YoloSource)
	fmt.Fprintln(writer)
	writer.Flush()

	final := config.FileConfig{}

	overrideProvider, err := askYesNo(reader, writer, "Override provider for this project", existing.Provider != "")
	if err != nil {
		return err
	}
	providerForLocal := current.Provider
	providerChanged := false
	if overrideProvider {
		providerDefault := firstNonEmpty(normalizeProvider(existing.Provider), normalizeProvider(current.Provider))
		provider, err := askProviderChoice(reader, writer, providerDefault)
		if err != nil {
			return err
		}
		final.Provider = provider
		providerForLocal = provider
		providerChanged = normalizeProvider(provider) != normalizeProvider(current.Provider)
	}
	spec := providerInfo(providerForLocal)
	promptExisting := scopedConfigForProvider(existing, providerForLocal)
	if providerChanged && normalizeProvider(existing.Provider) == "" {
		promptExisting.Host = ""
		promptExisting.Model = ""
	}
	currentHostForProvider := ""
	currentModelForProvider := ""
	if normalizeProvider(current.Provider) == normalizeProvider(providerForLocal) {
		currentHostForProvider = current.Host
		currentModelForProvider = current.Model
	}

	overrideHostDefault := promptExisting.Host != "" || (providerChanged && spec.DefaultHost != "")
	overrideHost, err := askYesNo(reader, writer, "Override host for this project", overrideHostDefault)
	if err != nil {
		return err
	}
	if overrideHost {
		hostDefault := firstNonEmpty(promptExisting.Host, currentHostForProvider, spec.DefaultHost)
		if hostDefault == "" {
			fmt.Fprintln(writer, "This provider does not need a host override by default.")
		} else {
			host, err := askLine(reader, writer, firstNonEmpty(spec.HostLabel, "Project host"), hostDefault)
			if err != nil {
				return err
			}
			final.Host = host
		}
	}

	overrideModelDefault := promptExisting.Model != "" || providerChanged
	overrideModel, err := askYesNo(reader, writer, "Override model for this project", overrideModelDefault)
	if err != nil {
		return err
	}
	if overrideModel {
		hostForModel := firstNonEmpty(final.Host, promptExisting.Host, currentHostForProvider, spec.DefaultHost)
		model, err := askProviderModel(reader, writer, spec, promptExisting, hostForModel, currentModelForProvider)
		if err != nil {
			return err
		}
		final.Model = model
	}

	yoloChoice, err := askChoice(reader, writer, "Project yolo mode", []string{"inherit", "on", "off"}, triStateDefault(existing.YoloMode))
	if err != nil {
		return err
	}
	switch yoloChoice {
	case "on":
		final.YoloMode = boolPtr(true)
	case "off":
		final.YoloMode = boolPtr(false)
	default:
		final.YoloMode = nil
	}

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Local override summary:")
	fmt.Fprintf(writer, "  provider: %s\n", localValueOrInherit(final.Provider))
	fmt.Fprintf(writer, "  host: %s\n", localValueOrInherit(final.Host))
	fmt.Fprintf(writer, "  model: %s\n", localValueOrInherit(final.Model))
	fmt.Fprintf(writer, "  yoloMode: %s\n", triStateLabel(final.YoloMode))
	fmt.Fprintln(writer, "  credentials: never stored in local config; inherit from environment/global config")
	if len(spec.CredentialEnvs) > 0 {
		fmt.Fprintf(writer, "  credential envs: %s\n", strings.Join(spec.CredentialEnvs, " / "))
	}
	writer.Flush()

	if !hasConfig(final) {
		fmt.Fprintln(writer)
		fmt.Fprintln(writer, "No local overrides selected.")
		if hasConfig(existing) {
			remove, err := askYesNo(reader, writer, fmt.Sprintf("Remove existing local config at %s", path), true)
			if err != nil {
				return err
			}
			if !remove {
				fmt.Fprintln(writer, "Cancelled.")
				return nil
			}
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return err
			}
			fmt.Fprintf(writer, "Removed local config at %s\n", path)
			return nil
		}
		fmt.Fprintln(writer, "Nothing to write.")
		return nil
	}

	confirm, err := askYesNo(reader, writer, fmt.Sprintf("Write local config to %s", path), true)
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
	fmt.Fprintf(writer, "Wrote local config to %s\n", path)
	fmt.Fprintln(writer, "Run 'pls config show' here to verify the effective overrides.")
	return nil
}

func defaultHostForProvider(provider string) string {
	return providerInfo(provider).DefaultHost
}

func defaultModelForProvider(provider string) string {
	return providerInfo(provider).DefaultModel
}

func triStateDefault(value *bool) string {
	if value == nil {
		return "inherit"
	}
	if *value {
		return "on"
	}
	return "off"
}

func triStateLabel(value *bool) string {
	if value == nil {
		return "inherit"
	}
	if *value {
		return "on"
	}
	return "off"
}

func localValueOrInherit(value string) string {
	if strings.TrimSpace(value) == "" {
		return "inherit"
	}
	return value
}
