package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os/exec"
	"regexp"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: pterm.Blue("Update the alis_ CLI to the latest version"),
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Info.Printf("Current version: %s\n", VERSION)
		spinner, _ := pterm.DefaultSpinner.Start("Updating alis_ command line interface...")
		cmds := "go env -w GOPRIVATE=go.protobuf.alis.alis.exchange,github.com/alis-x/cli/alis,go.lib.alis.dev && go install github.com/alis-x/cli/alis@latest"
		pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		out, err := exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		if err != nil {
			pterm.Debug.Println(cmds)
			spinner.Fail(fmt.Sprintf("%s", out))
			return
		}

		// check latest version
		cmds = "alis --version"
		pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		out, err = exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		if err != nil {
			pterm.Debug.Println(cmds)
			spinner.Fail(fmt.Sprintf("%s", out))
			return
		}

		// get new version, if updated.
		v := regexp.MustCompile(`(?m)alis version (\d+.\d+.\d+)`).FindAllStringSubmatch(fmt.Sprintf("%s", out), -1)
		if v[0][1] == VERSION {
			spinner.Success("You already have the latest version installed.")
		} else {
			spinner.Success(fmt.Sprintf("Updated version: %s -> %s\n", VERSION, v[0][1]))
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
