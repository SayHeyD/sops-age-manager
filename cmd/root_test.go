package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestMousetrapDisabled(t *testing.T) {
	if cobra.MousetrapHelpText != "" {
		t.Errorf("expected cobra.MousetrapHelpText to be empty so Windows GUI launcher works, got %q", cobra.MousetrapHelpText)
	}
}
