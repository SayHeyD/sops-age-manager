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
