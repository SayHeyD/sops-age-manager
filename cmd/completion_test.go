package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/test"
	"github.com/spf13/cobra"
)

func TestCompleteKeyNamesWithAvailableKeys(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	key1Path := filepath.Join(testDir.Path, "key1.txt")
	key2Path := filepath.Join(testDir.Path, "key2.txt")
	nestedDir := filepath.Join(testDir.Path, "nested")
	if err := os.MkdirAll(nestedDir, 0700); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}
	key3Path := filepath.Join(nestedDir, "key3.txt")

	if err := os.WriteFile(key1Path, []byte(testKeyContent), 0600); err != nil {
		t.Fatalf("failed to write key1: %v", err)
	}
	if err := os.WriteFile(key2Path, []byte(testKeyContent), 0600); err != nil {
		t.Fatalf("failed to write key2: %v", err)
	}
	if err := os.WriteFile(key3Path, []byte(testKeyContent), 0600); err != nil {
		t.Fatalf("failed to write key3: %v", err)
	}

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	names, directive := completeKeyNames(useKeyCommand, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected directive %v, got %v", cobra.ShellCompDirectiveNoFileComp, directive)
	}

	sort.Strings(names)
	expected := []string{"key1", "key2", "nested/key3"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d keys, got %d: %v", len(expected), len(names), names)
	}
	for i := range expected {
		if names[i] != expected[i] {
			t.Errorf("expected key %q, got %q", expected[i], names[i])
		}
	}
}

func TestCompleteKeyNamesPositionalBoundary(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	key1Path := filepath.Join(testDir.Path, "key1.txt")
	if err := os.WriteFile(key1Path, []byte(testKeyContent), 0600); err != nil {
		t.Fatalf("failed to write key1: %v", err)
	}

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	names, directive := completeKeyNames(useKeyCommand, []string{"key1"}, "")
	if names != nil {
		t.Errorf("expected nil names when positional arg already provided, got %v", names)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected directive %v, got %v", cobra.ShellCompDirectiveNoFileComp, directive)
	}
}

func TestCompleteKeyNamesEmptyOrMissingKeyDir(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	emptyDir := filepath.Join(testDir.Path, "empty-keys")
	if err := os.MkdirAll(emptyDir, 0700); err != nil {
		t.Fatalf("failed to create empty dir: %v", err)
	}

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", emptyDir)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	names, directive := completeKeyNames(useKeyCommand, []string{}, "")
	if names != nil {
		t.Errorf("expected nil names for empty key dir, got %v", names)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected directive %v, got %v", cobra.ShellCompDirectiveNoFileComp, directive)
	}
}

func TestValidArgsFunctionWired(t *testing.T) {
	if useKeyCommand.ValidArgsFunction == nil {
		t.Error("expected useKeyCommand to have ValidArgsFunction set")
	}
	if copyKeyCommand.ValidArgsFunction == nil {
		t.Error("expected copyKeyCommand to have ValidArgsFunction set")
	}
}

func TestCompletionScriptGeneration(t *testing.T) {
	tests := []struct {
		shell       string
		expectedStr string
	}{
		{"bash", "bash completion"},
		{"BASH", "bash completion"},
		{"zsh", "_sam"},
		{"ZSH", "_sam"},
		{"fish", "fish completion"},
		{"powershell", "Register-ArgumentCompleter"},
		{"PowerShell", "Register-ArgumentCompleter"},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := &cobra.Command{Use: "sam"}
			completionCmd := *completionCommand
			cmd.AddCommand(&completionCmd)
			completionCmd.SetOut(&buf)

			err := completionCmd.RunE(&completionCmd, []string{tt.shell})
			if err != nil {
				t.Fatalf("unexpected error generating %s completion: %v", tt.shell, err)
			}

			output := buf.String()
			if len(output) == 0 {
				t.Fatalf("expected non-empty completion script for shell %s", tt.shell)
			}
			if !strings.Contains(output, tt.expectedStr) {
				t.Errorf("expected completion script for %s to contain %q", tt.shell, tt.expectedStr)
			}
		})
	}
}

func TestCompletionUnsupportedShell(t *testing.T) {
	var buf bytes.Buffer
	cmd := &cobra.Command{Use: "sam"}
	completionCmd := *completionCommand
	cmd.AddCommand(&completionCmd)
	completionCmd.SetOut(&buf)

	err := completionCmd.RunE(&completionCmd, []string{"unsupported-shell"})
	if err == nil {
		t.Fatal("expected error for unsupported shell, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported shell type \"unsupported-shell\"") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCompletionCommandExecutionViaRoot(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	RootCmd.SetArgs([]string{"completion", "bash"})
	err = RootCmd.Execute()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("unexpected error executing RootCmd completion: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "bash completion") && !strings.Contains(output, "__complete") {
		t.Errorf("expected bash completion script from RootCmd, got: %q", output)
	}
}
