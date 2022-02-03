package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	pbProducts "go.protobuf.alis.alis.exchange/alis/os/resources/products/v1"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// orgCmd represents the product command
var orgCmd = &cobra.Command{
	Use:   "org",
	Short: pterm.Blue("Manages organisations."),
	Long:  pterm.Green("Use this command to manage an organisation."),
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Error.Println("a valid command is missing\nplease run 'alis org -h' for details.")
	},
	//Example: pterm.LightYellow("alis org ali"),
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.SilenceUsage = true
	orgCmd.SilenceErrors = true
	orgCmd.AddCommand(createOrgCmd)
	orgCmd.AddCommand(getOrgCmd)
	orgCmd.AddCommand(listOrgCmd)
	orgCmd.AddCommand(clearOrgCmd)
}

// createOrgCmd represents the create command
var createOrgCmd = &cobra.Command{
	Use:     "create",
	Short:   pterm.Green("Create a new organisation"),
	Long:    pterm.Green(`Creates a new organisation.`),
	Args:    validateOrgArg,
	Example: pterm.LightYellow("alis org create mycompany"),
	Run: func(cmd *cobra.Command, args []string) {
		organisationID = args[0]

		pterm.Info.Println(cmd.Context())

		// request domain
		domain, err := askUserString("Service domain (for example, alis.services, rezco.services): ", `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		// Create a new product resource
		op, err := alisProductsClient.CreateOrganisation(cmd.Context(), &pbProducts.CreateOrganisationRequest{
			Organisation: &pbProducts.Organisation{
				DisplayName: strings.ToTitle(organisationID),
				State:       pbProducts.Organisation_DEV,
				Owner:       "jan@alis.capital",
				Domain:      domain,
			},
			OrganisationId: organisationID,
		})
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		// wait for the long-running operation to complete.
		err = wait(cmd.Context(), op, "Creating "+organisationID, "Created "+organisationID, 300, true)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
	},
}

// getOrgCmd represents the get command
var getOrgCmd = &cobra.Command{
	Use:   "get",
	Short: pterm.Blue("Retrieves a specified organisation"),
	Long: pterm.Green(
		`This method clones or updates the specified organisation to your local environment 

'google' is a special type of organisation you could pull to gain local access to 
its common protocol buffers.  If you are following the Google API design guidelines,
you most likely will have to run the command: "alis org get google"`),
	Run: func(cmd *cobra.Command, args []string) {
		organisationID = args[0]

		// Google is a special organisation for which we need to perform a custom proto pull.
		if organisationID == "google" {
			// update google common protos.
			spinner, _ := pterm.DefaultSpinner.Start("Updating " + homeDir + "/google/proto... ")
			cmds := "git -C $HOME/alis.exchange/google/proto pull --no-rebase || git clone https://github.com/googleapis/googleapis.git $HOME/alis.exchange/google/proto"
			pterm.Debug.Printf("Shell command:\n%s\n", cmds)
			out, err := exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
			if err != nil {
				pterm.Debug.Printf(fmt.Sprintf("%s", out))
				pterm.Error.Println(err)
				return
			}
			pterm.Debug.Printf(fmt.Sprintf("%s", out))
			spinner.Success("Updated " + homeDir + "/alis.exchange/google/proto. ")
			return
		}

		// Retrieve the organisation resource
		res, err := alisProductsClient.GetOrganisation(cmd.Context(),
			&pbProducts.GetOrganisationRequest{Name: "organisations/" + organisationID})
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		// Clone the proto repository
		spinner, _ := pterm.DefaultSpinner.Start("Updating " + homeDir + "/alis.exchange/" + organisationID + "/proto... ")
		cmds := "git -C $HOME/alis.exchange/" + organisationID + "/proto pull --no-rebase || gcloud source repos clone proto $HOME/alis.exchange/" + organisationID + "/proto --project=" + res.GetGoogleProjectId()
		pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		out, err := exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		if err != nil {
			pterm.Debug.Printf(fmt.Sprintf("%s", out))
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf(fmt.Sprintf("%s", out))

		spinner.Success("Updated repository " + homeDir + "/alis.exchange/" + organisationID + "/proto. ")

		// Clone the protobuf-go repository
		spinner, _ = pterm.DefaultSpinner.Start("Updating " + homeDir + "/alis.exchange/" + organisationID + "/protobuf/go... ")
		cmds = "git -C $HOME/alis.exchange/" + organisationID + "/protobuf/go pull --no-rebase || gcloud source repos clone protobuf-go $HOME/alis.exchange/" + organisationID + "/protobuf/go --project=" + res.GetGoogleProjectId()
		pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		out, err = exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		if err != nil {
			pterm.Debug.Printf(fmt.Sprintf("%s", out))
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf(fmt.Sprintf("%s", out))

		spinner.Success("Updated repository " + homeDir + "/alis.exchange/" + organisationID + "/protobuf/go. ")

		// Clone the api-go repository
		spinner, _ = pterm.DefaultSpinner.Start("Updating " + homeDir + "/alis.exchange/" + organisationID + "/api/go... ")
		cmds = "git -C $HOME/alis.exchange/" + organisationID + "/api/go pull --no-rebase || gcloud source repos clone api-go $HOME/alis.exchange/" + organisationID + "/api/go --project=" + res.GetGoogleProjectId()
		pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		out, err = exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		if err != nil {
			pterm.Debug.Printf(fmt.Sprintf("%s", out))
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf(fmt.Sprintf("%s", out))

		spinner.Success("Updated repository " + homeDir + "/alis.exchange/" + organisationID + "/api/go. ")

		// Clone the protobuf-python repository
		spinner, _ = pterm.DefaultSpinner.Start("Updating " + homeDir + "/alis.exchange/" + organisationID + "/protobuf/python... ")
		cmds = "git -C $HOME/alis.exchange/" + organisationID + "/protobuf/python pull --no-rebase || gcloud source repos clone protobuf-python $HOME/alis.exchange/" + organisationID + "/protobuf/python --project=" + res.GetGoogleProjectId()
		pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		out, err = exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		if err != nil {
			pterm.Debug.Printf(fmt.Sprintf("%s", out))
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf(fmt.Sprintf("%s", out))

		spinner.Success("Updated repository " + homeDir + "/alis.exchange/" + organisationID + "/protobuf/python. ")

		pterm.Debug.Printf("Retrieved Organisation:\n%s\n", res)
		ptermTip.Println("Are you making use of Google protocol buffers?\nRun `alis org get google` to download a local copy\nof of their common protocol buffers as well.")
	},
	Args:    validateOrgArg,
	Example: pterm.LightYellow("alis org get {organisationID}"),
}

// clearOrgCmd represents the clear command
var clearOrgCmd = &cobra.Command{
	Use:   "clear",
	Short: pterm.Blue("Clears the org from your local environment"),
	Long: pterm.Green(
		`This method removes the specified organisation from your local environment.

This will remove all underlying products, protos, protobufs etc..

Please clear organisations not actively working on - its not great to leave 
these lying around in your local development environment.`),
	Run: func(cmd *cobra.Command, args []string) {
		organisationID = strings.Split(args[0], ".")[0]

		// Retrieve the organisation resource
		organisation, err := alisProductsClient.GetOrganisation(cmd.Context(),
			&pbProducts.GetOrganisationRequest{Name: "organisations/" + organisationID})
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf("Organisation:\n%s\n", organisation)

		orgPath := homeDir + "/alis.exchange/" + organisationID
		pterm.Warning.Printf("Removing product '%s' from your local environment.\nFolder location: %s\n", organisationID, orgPath)
		pterm.Warning.Printf("Please also ensure you close any IDEs (pointing to the \norganisation resources (protos, etc) or any underlying \nproducts) you may have open.\n")
		userInput, err := askUserString("Are you sure? (y/n): ", `^[y|n]$`)
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		if userInput == "y" {
			cmds := "rm -rf " + orgPath
			pterm.Debug.Printf("Shell command:\n%s\n", cmds)
			out, err := exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
			if err != nil {
				pterm.Error.Printf(fmt.Sprintf("%s", out))
				return
			}
			pterm.Debug.Printf(fmt.Sprintf("%s", out))
			pterm.Success.Printf("Removed product `%s` from your local environment.\nFolder removed: %s\n", organisationID, orgPath)
		} else {
			pterm.Warning.Printf("Aborted operation.\n Did not remove %s\n", orgPath)
		}

	},
	Args:    validateOrgArg,
	Example: pterm.LightYellow("alis org clear {orgID}"),
}

// listOrgCmd represents the list command
var listOrgCmd = &cobra.Command{
	Use:   "list",
	Short: pterm.Blue("Lists all organisations"),
	//Long: pterm.Green(
	//	`This method lists all the products for a given organisation`),
	Run: func(cmd *cobra.Command, args []string) {

		// Retrieve the organisation resource
		organisations, err := alisProductsClient.ListOrganisations(cmd.Context(),
			&pbProducts.ListOrganisationsRequest{})
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf("Organisation:\n%s\n", organisations.GetOrganisations())

		table := pterm.TableData{{"Index", "OrganisationID", "Display Name", "Owner", "Google Project", "Resource Name", "State", "Updated"}}
		for i, organisation := range organisations.GetOrganisations() {
			resourceID := strings.Split(organisation.GetName(), "/")[1]
			table = append(table, []string{
				strconv.Itoa(i), resourceID, organisation.GetDisplayName(),
				organisation.GetOwner(), organisation.GetGoogleProjectId(),
				organisation.GetName(), organisation.GetState().String(),
				organisation.GetUpdateTime().AsTime().Format(time.RFC3339)})
		}

		err = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(table).Render()
		if err != nil {
			return
		}

	},
	//Args: validateOrgArg,
	Example: pterm.LightYellow("alis org list"),
}
