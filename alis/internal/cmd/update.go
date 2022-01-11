package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"os/exec"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: pterm.Blue("Update the alis_ CLI to the latest version"),
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Info.Printf("Current version: %s\n", VERSION)
		spinner, _ := pterm.DefaultSpinner.Start("Updating alis_ command line interface...")
		cmds := "go env -w GOPRIVATE=go.protobuf.alis.alis.exchange,github.com/alis-exchange/cli/alis,go.lib.alis.dev && go install github.com/alis-exchange/cli/alis@latest"
		out, err := exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		if err != nil {
			pterm.Debug.Println(cmds)
			spinner.Fail(fmt.Sprintf("%s", out))
			return
		}
		spinner.Success("Updated alis_ CLI to the latest version.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
