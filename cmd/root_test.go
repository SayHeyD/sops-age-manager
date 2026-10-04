package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/spf13/cobra"
)

func TestMousetrapDisabled(t *testing.T) {
	if cobra.MousetrapHelpText != "" {
		t.Errorf("expected cobra.MousetrapHelpText to be empty so Windows GUI launcher works, got %q", cobra.MousetrapHelpText)
	}
}

func TestExecuteSopsUIModePrintsMessage(t *testing.T) {
	origLaunchUI := launchUIFunc
	defer func() { launchUIFunc = origLaunchUI }()

	var launched bool
	launchUIFunc = func(cfg *config.Config, logo []byte) {
		launched = true
	}

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	executeSops([]string{})

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !launched {
		t.Error("expected launchUI to be called")
	}

	if !strings.Contains(output, "Running UI ...") {
		t.Errorf("expected output to contain 'Running UI ...', got %q", output)
	}

	if !strings.Contains(output, "Ctrl+C to cancel") {
		t.Errorf("expected output to contain 'Ctrl+C to cancel', got %q", output)
	}
}

func TestExecuteSopsVersionFlag(t *testing.T) {
	origLaunchUI := launchUIFunc
	origAppVersion := appVersion
	origShowVersion := showVersion
	defer func() {
		launchUIFunc = origLaunchUI
		appVersion = origAppVersion
		showVersion = origShowVersion
	}()

	var launched bool
	launchUIFunc = func(cfg *config.Config, logo []byte) {
		launched = true
	}

	appVersion = "1.2.3-test"
	showVersion = true

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	executeSops([]string{})

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if launched {
		t.Error("expected launchUI NOT to be called when showVersion is true")
	}

	if !strings.Contains(output, "sam version: 1.2.3-test") {
		t.Errorf("expected output to contain 'sam version: 1.2.3-test', got %q", output)
	}
}

func TestRootCmdVersionFlagExecution(t *testing.T) {
	origLaunchUI := launchUIFunc
	origAppVersion := appVersion
	defer func() {
		launchUIFunc = origLaunchUI
		appVersion = origAppVersion
	}()

	var launched bool
	launchUIFunc = func(cfg *config.Config, logo []byte) {
		launched = true
	}

	appVersion = "2.0.0-test"

	for _, flag := range []string{"--version", "-v"} {
		t.Run(flag, func(t *testing.T) {
			launched = false
			origStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdout = w

			RootCmd.SetArgs([]string{flag})
			err = RootCmd.Execute()

			_ = w.Close()
			os.Stdout = origStdout

			if err != nil {
				t.Fatalf("unexpected error executing RootCmd with %s: %v", flag, err)
			}

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			if launched {
				t.Errorf("expected launchUI NOT to be called with %s", flag)
			}

			if !strings.Contains(output, "sam version: 2.0.0-test") {
				t.Errorf("expected output to contain 'sam version: 2.0.0-test', got %q", output)
			}
		})
	}
}
