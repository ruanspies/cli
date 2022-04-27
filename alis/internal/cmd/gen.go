package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	pbProducts "go.protobuf.alis.alis.exchange/alis/os/resources/products/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// neuronCmd represents the neuron command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: pterm.Blue("Generates code"),
	Long:  pterm.Green(`Use this command to generate code.`),
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Error.Println("a valid command is missing\nplease run 'alis gen -h' for details.")
	},
}

var (
	pushPublicProtocolBuffers bool
)

// protobufGenCmd represents the genproto command
var protobufGenCmd = &cobra.Command{
	Use:   "protobuf",
	Short: pterm.Blue("Generates the protocol buffers files for a specified language.  Golang is the default."),
	Long: pterm.Green(
		`This method uses the 'genproto-go' command line to generate protocol buffers
for the specified neuron.  These are generated locally only and should be used for local development of your gRPC services.

These are used in combination with the 'Replace ...' command in your go.mod file to 'point' to the local, instead of the
official protobufs.'

Once local development is done, use the '--push' flag to push the generated protocol buffers to the go.protobuf repository.
Once the protobufs are pushed, you should remove the 'Replace... ' command in your go.mod file and run 'go mod tidy' to pull
the latest protobufs from the repo into your gRPC service.`),
	Example: pterm.LightYellow("alis gen protobuf {orgID}.{productID}.{neuronID}"),
	Args:    validateNeuronArg,
	Run: func(cmd *cobra.Command, args []string) {
		organisationID = strings.Split(args[0], ".")[0]
		productID = strings.Split(args[0], ".")[1]
		neuronID = strings.Split(args[0], ".")[2]

		// Retrieve the organisation resource
		organisation, err := alisProductsClient.GetOrganisation(cmd.Context(),
			&pbProducts.GetOrganisationRequest{Name: "organisations/" + organisationID})
		if err != nil {
			// TODO: handle not found by listing available organisations.
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf("Get Organisation:\n%s\n", organisation)

		// Retrieve the neuron resource
		neuron, err := alisProductsClient.GetNeuron(cmd.Context(),
			&pbProducts.GetNeuronRequest{
				Name: "organisations/" + organisationID + "/products/" + productID + "/neurons/" + neuronID})
		if err != nil {
			// TODO: handle not found by listing available products.
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf("Get Neuron:\n%s\n", neuron)

		// Generate the protocol buffers for Golang
		if genprotoGo {

			var (
				neuronProtobufFullPath string
				neuronProtoFullPath    string
				protobufGoRepoPath     string

				cmds string
			)

			// set required path variables
			if pushPublicProtocolBuffers {
				neuronProtobufFullPath = homeDir + "/alis.exchange/" + organisationID + "/public/protobuf/go/" + organisationID + "/" + productID + "/" + strings.ReplaceAll(neuronID, "-", "/")
				neuronProtoFullPath = homeDir + "/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + "/" + strings.ReplaceAll(neuronID, "-", "/")
				protobufGoRepoPath = homeDir + "/alis.exchange/" + organisationID + "/public/protobuf/go"
				relativeProtoPath := organisationID + "/" + productID + "/" + strings.ReplaceAll(neuronID, "-", "/")

				if pushProtocolBuffers {
					err := clearUncommittedRepoChanges(protobufGoRepoPath)
					if err != nil {
						pterm.Error.Println(err)
						return
					}
				}

				cmds = "rm -rf " + neuronProtobufFullPath + " && " +
					"mkdir -p " + neuronProtobufFullPath
				pterm.Debug.Printf("Shell command:\n%s\n", cmds)
				out, err := exec.CommandContext(context.Background(), "bash", "-c", cmds).CombinedOutput()
				if strings.Contains(fmt.Sprintf("%s", out), "warning") {
					pterm.Warning.Print(fmt.Sprintf("removing existing local directory content...\n%s", out))
				}
				pterm.Debug.Printf("Cleared the local files in directory: %s\n", neuronProtobufFullPath)

				var relativeProtoPaths string
				files, err := ioutil.ReadDir(neuronProtoFullPath)
				for _, f := range files {
					if strings.HasSuffix(f.Name(), ".proto") {
						relativeProtoPaths += relativeProtoPath + "/" + f.Name() + " "
					}
				}
				pterm.Debug.Printf("Relative protopaths: %s\n", relativeProtoPaths)

				descriptorPath, err := generatePublicLocalDescriptorFileFromNeuron(cmd.Context(), neuron.GetName(), neuronProtobufFullPath)
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				pterm.Debug.Println("Successfully created public scoped descriptor.pb. Destination: %s", *descriptorPath)

				// Use the public scoped descriptor.pb to generate the Go files
				cmds = "go env -w GOPRIVATE=go.lib." + organisationID + ".alis.exchange,go.protobuf." + organisationID + ".alis.exchange,proto." + organisationID + ".alis.exchange,cli.alis.dev && " +
					"protoc --go_out=" + protobufGoRepoPath + " --go_opt=paths=source_relative --go-grpc_out=" + protobufGoRepoPath + " --go-grpc_opt=paths=source_relative -I=$HOME/alis.exchange/google/proto --descriptor_set_in=" + *descriptorPath + " " + relativeProtoPaths + " && " +
					"rm -f " + *descriptorPath // remove the public scoped descriptor.pb file

			} else {
				neuronProtobufFullPath = homeDir + "/alis.exchange/" + organisationID + "/protobuf/go/" + organisationID + "/" + productID + "/" + strings.ReplaceAll(neuronID, "-", "/")
				neuronProtoFullPath = homeDir + "/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + "/" + strings.ReplaceAll(neuronID, "-", "/")
				protobufGoRepoPath = homeDir + "/alis.exchange/" + organisationID + "/protobuf/go"

				if pushProtocolBuffers {
					err := clearUncommittedRepoChanges(protobufGoRepoPath)
					if err != nil {
						pterm.Error.Println(err)
						return
					}
				}

				cmds = "rm -rf " + neuronProtobufFullPath + " && " +
					"mkdir -p " + neuronProtobufFullPath + " && " +
					"go env -w GOPRIVATE=go.lib." + organisationID + ".alis.exchange,go.protobuf." + organisationID + ".alis.exchange,proto." + organisationID + ".alis.exchange,cli.alis.dev && " +
					"protoc --go_out=$HOME/alis.exchange/" + organisationID + "/protobuf/go --go_opt=paths=source_relative --go-grpc_out=$HOME/alis.exchange/" + organisationID + "/protobuf/go --go-grpc_opt=paths=source_relative -I=$HOME/alis.exchange/google/proto -I=$HOME/alis.exchange/" + organisationID + "/proto " + neuronProtoFullPath + "/*.proto"
			}

			pterm.Debug.Printf("Shell command:\n%s\n", cmds)
			out, err := exec.CommandContext(context.Background(), "bash", "-c", cmds).CombinedOutput()
			if err != nil {
				pterm.Error.Printf(fmt.Sprintf("%s", out))
				pterm.Error.Println(err)
				return
			}
			if strings.Contains(fmt.Sprintf("%s", out), "warning") {
				pterm.Warning.Print(fmt.Sprintf("Generating protocol buffers for go...\n%s", out))
			}
			pterm.Success.Printf("Generated protocol buffers for Go.\nProto source: %s\n", neuronProtoFullPath)

			// generate ProductDescriptorFile at product level.
			err = genProductDescriptorFile("organisations/" + organisationID + "/products/" + productID)
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Success.Println("Generated Product Descriptor File")

			// Publish to protobuf repository if not in local mode.
			if pushProtocolBuffers {
				// commit protocol buffers in go
				message := fmt.Sprintf("chore(%s): updated by alis_ CLI", neuronID)
				_, err := commitTagAndPush(cmd.Context(), protobufGoRepoPath, neuronProtobufFullPath,
					message, "", true, true)
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				pterm.Success.Println("Published protocol buffers for Go")

				ptermTip.Printf("Now that your protobuf if updated, please ensure that you update your \n" +
					"go.mod file to reflect this new version of your protobuf.\n")
			} else {
				ptermTip.Printf("The protobufs were generated for local development use only. To formally\n" +
					"publish them use the `--push` flag to publish them to the \n" +
					"protobuf libraries.\n")
			}
		}

		// generate protocol buffers for Python
		if genprotoPython {
			neuronProtobufFullPath := homeDir + "/alis.exchange/" + organisationID + "/protobuf/python/" + organisationID + "/" + productID + "/" + strings.ReplaceAll(neuronID, "-", "/")
			neuronProtoFullPath := homeDir + "/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + "/" + strings.ReplaceAll(neuronID, "-", "/")

			// TODO: create __init__.py using golang file io.
			cmds := "rm -rf " + neuronProtobufFullPath + " && " +
				"mkdir -p " + neuronProtobufFullPath + " && " +
				`echo "__import__('pkg_resources').declare_namespace(__name__)" > $HOME/alis.exchange/` + organisationID + `/protobuf/python/` + organisationID + `/` + productID + `/__init__.py` + " && " +
				`echo "__import__('pkg_resources').declare_namespace(__name__)" > $HOME/alis.exchange/` + organisationID + `/protobuf/python/` + organisationID + `/` + productID + `/` + strings.Split(neuronID, "-")[0] + `/__init__.py` + " && " +
				`echo "__import__('pkg_resources').declare_namespace(__name__)" > $HOME/alis.exchange/` + organisationID + `/protobuf/python/` + organisationID + `/` + productID + `/` + strings.Split(neuronID, "-")[0] + `/` + strings.Split(neuronID, "-")[1] + `/__init__.py` + " && " +
				`echo "__import__('pkg_resources').declare_namespace(__name__)" > $HOME/alis.exchange/` + organisationID + `/protobuf/python/` + organisationID + `/` + productID + `/` + strings.Split(neuronID, "-")[0] + `/` + strings.Split(neuronID, "-")[1] + `/` + strings.Split(neuronID, "-")[2] + `/__init__.py` + " && " +
				"python3 -m grpc_tools.protoc --python_out=$HOME/alis.exchange/" + organisationID + "/protobuf/python --grpc_python_out=$HOME/alis.exchange/" + organisationID + "/protobuf/python -I=$HOME/alis.exchange/google/proto -I=$HOME/alis.exchange/" + organisationID + "/proto " + neuronProtoFullPath + "/*.proto"
			pterm.Debug.Printf("Shell command:\n%s\n", cmds)
			out, err := exec.CommandContext(context.Background(), "bash", "-c", cmds).CombinedOutput()
			if err != nil {
				pterm.Error.Printf(fmt.Sprintf("%s", out))
				pterm.Error.Println(err)
				return
			}
			if strings.Contains(fmt.Sprintf("%s", out), "warning") {
				pterm.Warning.Print(fmt.Sprintf("Generating protocol buffers for python...\n%s", out))
			}
			pterm.Success.Printf("Generated protocol buffers for Python.\nProto source: %s\n", neuronProtoFullPath)

			// Publish to protobuf repository if not in local mode.
			if pushProtocolBuffers {

				protobufPythonRepo := fmt.Sprintf("%s/alis.exchange/%s/protobuf/python", homeDir, organisationID)

				// bump setup.py version
				var tpl bytes.Buffer
				setupPyTemplate, err := template.ParseFiles(protobufPythonRepo + "/setup.py")
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				if err = setupPyTemplate.Execute(&tpl, struct{}{}); err != nil {
					pterm.Error.Println(err)
					return
				}
				rel := regexp.MustCompile("version=\"(.*)\",")
				versionComponents := strings.Split(rel.FindStringSubmatch(tpl.String())[1], ".")
				minorVersion, err := strconv.Atoi(versionComponents[2])
				versionComponents[2] = strconv.Itoa(minorVersion + 1)
				newVersion := rel.ReplaceAllString(tpl.String(), "version=\""+strings.Join(versionComponents, ".")+"\",")
				out, err = exec.CommandContext(context.Background(), "bash", "-c", "echo "+"'"+newVersion+"'"+" > "+protobufPythonRepo+"/setup.py").CombinedOutput()
				if err != nil {
					pterm.Error.Printf(fmt.Sprintf("%s", out))
					pterm.Error.Println(err)
					return
				}

				// publish Python package to artifact registry
				tpl = bytes.Buffer{}
				publishTemplate, err := TemplateFs.ReadFile("internal/cmd/neuron/python/publishPython.sh")
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				t, err := template.New("file").Parse(string(publishTemplate))
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				if err := t.Execute(&tpl, struct {
					OrgProjectID   string
					OrganisationID string
				}{
					OrgProjectID:   organisation.GetGoogleProjectId(),
					OrganisationID: organisationID,
				}); err != nil {
					pterm.Error.Println(err)
					return
				}

				out, err = exec.CommandContext(context.Background(), "bash", "-c", tpl.String()).CombinedOutput()
				if err != nil {
					pterm.Error.Printf(fmt.Sprintf("%s", out))
					pterm.Error.Println(err)
					return
				}

				if strings.Contains(fmt.Sprintf("%s", out), "warning") {
					pterm.Warning.Print(fmt.Sprintf("Publishing protocol buffers for python...\n%s", out))
				}

				message := fmt.Sprintf("chore(%s): updated by alis_ CLI", neuronID)
				commitPath := neuronProtobufFullPath +
					" " + protobufPythonRepo + "/setup.py" +
					" " + protobufPythonRepo + "/alis/" + productID + "/__init__.py" +
					" " + protobufPythonRepo + "/alis/" + productID + "/" + strings.Split(neuronID, "-")[0] + "/__init__.py" +
					" " + protobufPythonRepo + "/alis/" + productID + "/" + strings.Split(neuronID, "-")[0] + "/" + strings.Split(neuronID, "-")[1] + "/__init__.py"

				_, err = commitTagAndPush(cmd.Context(), protobufPythonRepo, commitPath, message,
					"", true, true)
				if err != nil {
					return
				}

				pterm.Success.Println("Published protocol buffers for Python")
			} else {
				ptermTip.Printf("The protobufs were generated for local development use only. To formally\n" +
					"publish them use the `--push` flag to publish them to the \n" +
					"protobuf libraries.\n")
			}
		}

		return
	},
}

// descriptorGenCmd represents the descriptor command
var descriptorGenCmd = &cobra.Command{
	Use:   "descriptor",
	Short: pterm.Blue("Generates a file descriptor for the specified organisation, product or neuron."),
	Long: pterm.Green(
		`This method uses the 'protoc --descriptor_set_out=...' command line to generate a local 'descriptor.pb' file.
This file is a serialised google.protobuf.FileDescriptorSet object representing all the relevant proto files.`),
	Example: pterm.LightYellow("alis gen descriptor {orgID}.{productID}.{neuronID}"),
	Args:    validateOrgOrProductOrNeuron,
	Run: func(cmd *cobra.Command, args []string) {

		var name string
		argParts := strings.Split(args[0], ".")

		// The length of the argument parts determine whether the request is at organisation, product or neuron level.
		switch len(argParts) {
		case 1:
			// Retrieve the organisation resource
			name = "organisations/" + argParts[0]
			organisation, err := alisProductsClient.GetOrganisation(cmd.Context(),
				&pbProducts.GetOrganisationRequest{Name: name})
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Debug.Printf("Get Organisation:\n%s\n", organisation)
		case 2:
			// Retrieve the product resource
			name = "organisations/" + argParts[0] + "/products/" + argParts[1]
			product, err := alisProductsClient.GetProduct(cmd.Context(), &pbProducts.GetProductRequest{Name: name})
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Debug.Printf("Get Product:\n%s\n", product)
		case 3:
			// Retrieve the neuron resource
			name = "organisations/" + argParts[0] + "/products/" + argParts[1] + "/neurons/" + argParts[2]
			neuron, err := alisProductsClient.GetNeuron(cmd.Context(),
				&pbProducts.GetNeuronRequest{
					Name: name})
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Debug.Printf("Get Neuron:\n%s\n", neuron)
		}

		// generate ProductDescriptorFile at the relevant resource level.
		descriptorPath, err := genDescriptorFile(name)
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		pterm.Info.Printf("Generated %s\n", descriptorPath)

		return
	},
}

var (
	productsDocsGenCustomFlag bool
	productsDocsGenPublicFlag bool
)

// productDocsGenCmd represents the gendocs command
var productDocsGenCmd = &cobra.Command{
	Use:   "docs",
	Short: pterm.Blue("Generates documentation files for all proto files in the specified product."),
	Long: pterm.Green(
		`This method uses the 'gendocs' protoc plugin to generate documentation for the
specified product and publishes it at a URL specified.`),
	Example: pterm.LightYellow("alis gen docs {orgID}.{productID}"),
	Args:    validateProductArg,
	Run: func(cmd *cobra.Command, args []string) {
		organisationID = strings.Split(args[0], ".")[0]
		productID = strings.Split(args[0], ".")[1]

		// Retrieve the organisation resource
		organisation, err := alisProductsClient.GetOrganisation(cmd.Context(),
			&pbProducts.GetOrganisationRequest{Name: "organisations/" + organisationID})
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf("GetOrganisation:\n%s\n", organisation)

		// Retrieve the product resource
		product, err := alisProductsClient.GetProduct(cmd.Context(),
			&pbProducts.GetProductRequest{Name: "organisations/" + organisationID + "/products/" + productID})
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf("GetProduct:\n%s\n", product)

		var envs []*pbProducts.Neuron_Env
		var op *longrunning.Operation

		pterm.Println("")
		pterm.Info.Println("1. First select a product deployment for which to generate documentation.")
		ptermTip.Println("The documentation will be generated for all the neuron versions in the deployment.")
		deployment, err := selectProductDeployment(cmd.Context(), product.GetName())
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Printf("GetDeployment:\n%s\n", deployment)
		envs = append(envs, &pbProducts.Neuron_Env{
			Name:  "ALIS_OS_DOCS_PRODUCT_DEPLOYMENT",
			Value: deployment.GetName(),
		})

		pterm.Println("")
		pterm.Info.Println("2. Specification of the documentation visibility scope.")
		var apiVisibility string
		if !productsDocsGenPublicFlag && !productsDocsGenCustomFlag {

			scope, err := askUserString("Specify either 'PUBLIC' or 'CUSTOM': ", "^(PUBLIC|CUSTOM)$")
			if err != nil {
				pterm.Error.Println(err)
				return
			}

			if scope == "PUBLIC" {
				productsDocsGenPublicFlag = true
			} else if scope == "CUSTOM" {
				productsDocsGenCustomFlag = true
			} else {
				pterm.Error.Println("Invalid scope specified.")
				return
			}
		}

		if productsDocsGenPublicFlag {
			ptermTip.Println("PUBLIC restriction scope will only contain documentation where there are no visibility restrictions specified. \n" +
				"If you want to restrict the visibility of the documentation, ensure that you have added visibility\n" +
				"restrictions to the proto.")
			pterm.Println("")
			pterm.Warning.Println("Specifying a PUBLIC restriction scope will make the generated documentation publicly available and " +
				"generate content on all the proto content that do not contain google.api.visibility options.")

			confirmPublic, err := askUserString("Are you sure you want to generate public documentation? (y/n): ", `^[y|n]$`)
			if err != nil {
				pterm.Error.Println(err)
				return
			}

			if confirmPublic == "n" {
				productsDocsGenPublicFlag = false
				productsDocsGenCustomFlag = true
				pterm.Println("")
				pterm.Info.Println("Alright. We will restricting access to the documentation and allow you to specify custom scopes.")
				pterm.Println("")
			}
		}

		if productsDocsGenCustomFlag {
			ptermTip.Println("CUSTOM will only contain documentation where the exact restriction scope is met and access to the" +
				"documentation will be regulated by a IAM group.")
			apiVisibility, err = askUserString("Specify the exact custom visibility scopes that should be matched."+
				"Multiple values may be seperated by a comma (Example: INTERNAL, PREVIEW): ", `^[A-Za-z0-9-, ]+$`)
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Debug.Println("Custom API visibility scope: ", apiVisibility)

			envs = append(envs, &pbProducts.Neuron_Env{
				Name:  "ALIS_OS_DOCS_RESTRICTION",
				Value: apiVisibility,
			})
		} else {
			pterm.Error.Println("Invalid scope specified.")
			return
		}

		////TODO: Requires definition of resource prior to implementation
		//pterm.Info.Println("3. WELCOME TEXT\n The documentation contains a generic welcome text that orientates users" +
		//	"on how to use the documentation.")
		//welcomeText, err := askUserString("Press enter to keep the generic text or type text to specify your own.", `.*`)
		//if err != nil {
		//	pterm.Error.Println(err)
		//	return
		//}
		//pterm.Debug.Println("Welcome text ", welcomeText)

		////TODO: Requires definition of resource prior to implementation
		//pterm.Info.Println("4. CUSTOM URL\n The documentation allows you to redirect readers to a custom URL such as your" +
		//	"product's marketing page")
		//customURL, err := askUserString("Full URL:", `.*`)
		//if err != nil {
		//	pterm.Error.Println(err)
		//	return
		//}
		//pterm.Debug.Println("Welcome text ", customURL)
		pterm.Println("")
		pterm.Info.Println("3. Specify where the base URL to host the documentation")
		dnsConfig, err := selectDnsConfig(cmd.Context())
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Println("DNS Config ", dnsConfig)
		envs = append(envs, &pbProducts.Neuron_Env{
			Name:  "ALIS_OS_DNS_PROJECT",
			Value: dnsConfig.project,
		})
		envs = append(envs, &pbProducts.Neuron_Env{
			Name:  "ALIS_OS_DNS_ZONE",
			Value: dnsConfig.zoneName,
		})
		pterm.Println("")
		pterm.Info.Println("4. Specify the custom URL that the documentation should be hosted on.")
		ptermTip.Println("Has to end with '" + dnsConfig.baseURL + "' (Example: myproduct." + dnsConfig.baseURL + ")")
		docsCustomURL, err := askUserString("Specify the custom URL: ", `[a-z.-]\.`+dnsConfig.baseURL)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Println("Custom URL: ", docsCustomURL)
		envs = append(envs, &pbProducts.Neuron_Env{
			Name:  "ALIS_OS_DNS_RECORD",
			Value: docsCustomURL,
		})

		// Create new product deployment
		// TODO: This should be made way more elegant for outside users.
		prodDeployment, err := createProductDeployment(cmd.Context(), "organisations/alis/products/ex")
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Debug.Println("New product deployment created: ", prodDeployment)

		// Retrieve the latest version
		neuronID := "resources-docs-v1"
		neuronName := "organisations/alis/products/ex/neurons/" + neuronID
		res, err := alisProductsClient.ListNeuronVersions(cmd.Context(), &pbProducts.ListNeuronVersionsRequest{
			Parent:   neuronName,
			ReadMask: &fieldmaskpb.FieldMask{Paths: []string{"version"}},
		})
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		if len(res.GetNeuronVersions()) == 0 {
			pterm.Error.Println("there are no versions available, please run `alis neron build ...` to create a version")
			return
		}

		latestVersion := res.GetNeuronVersions()[0].GetVersion()
		pterm.Debug.Println("Latest version: ", latestVersion)
		pterm.Debug.Println("Envs: ", envs)

		op, err = alisProductsClient.CreateNeuronDeployment(cmd.Context(), &pbProducts.CreateNeuronDeploymentRequest{
			Parent: prodDeployment.GetName(),
			NeuronDeployment: &pbProducts.NeuronDeployment{
				Name:    neuronName + "/versions/" + latestVersion,
				Version: latestVersion,
				Envs:    envs,
			},
			NeuronDeploymentId: neuronID,
		})
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		// check if we need to wait for operation to complete.
		if asyncFlag {
			pterm.Debug.Printf("GetOperation:\n%s\n", op)
			pterm.Success.Printf("Launched service in async mode.\n see long-running operation " + op.GetName() + " to monitor state\n")
		} else {

			successMessage := "Documentation has been deployed to: https://" + docsCustomURL + ".\n\n" +
				"NOTE that the DNS propagation and SSL certificate issuing can take up to 24 hours before the " +
				"site is available."

			if apiVisibility != "" {
				successMessage = successMessage + "\n" +
					"Manage access to documentation at: " + "https://groups.google.com/a/alis.exchange/g/" + prodDeployment.GetGoogleProjectId() + "/members"
			}

			successMessage = successMessage + "\n\n 🚩 For the documentation content to be generated, access to the documentation needs to be granted to: " +
				"alis-exchange@" + prodDeployment.GetGoogleProjectId() + ".iam.gserviceaccount.com in the following group:\n" +
				"https://console.cloud.google.com/iam-admin/groups/01fob9te30m28ba?organizationId=666464741197 \n" +
				"Contact any of the group managers for assistance."
			// wait for the long-running operation to complete.
			err := wait(cmd.Context(), op, "Creating documentation for "+product.GetName(), successMessage, 300, true)
			if err != nil {
				pterm.Error.Println(err)
				return
			}
		}

		///////////////////////////////
		///////////////////////////////
		//// Generate the index.html
		//cmds := "go env -w GOPRIVATE=go.lib." + organisationID + ".alis.exchange,go.protobuf." + organisationID + ".alis.exchange,proto." + organisationID + ".alis.exchange,cli.alis.dev &&" +
		//	"protoc --plugin=protoc-gen-doc=$HOME/go/bin/protoc-gen-doc --doc_out=$HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " --doc_opt=html,docs.html -I=$HOME/alis.exchange/google/proto -I=$HOME/alis.exchange/" + organisationID + "/proto $(find $HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " -iname \"*.proto\")"
		//pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		//out, err := exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		//if err != nil {
		//	pterm.Error.Printf(fmt.Sprintf("%s", out))
		//	pterm.Error.Println(err)
		//	return
		//}
		//
		//// Generate markdown
		//cmds = "go env -w GOPRIVATE=go.lib." + organisationID + ".alis.exchange,go.protobuf." + organisationID + ".alis.exchange,proto." + organisationID + ".alis.exchange,cli.alis.dev &&" +
		//	"protoc --plugin=protoc-gen-doc=$HOME/go/bin/protoc-gen-doc --doc_out=$HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " --doc_opt=markdown,docs.md -I=$HOME/alis.exchange/google/proto -I=$HOME/alis.exchange/" + organisationID + "/proto $(find $HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " -iname \"*.proto\")"
		//pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		//out, err = exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		//if err != nil {
		//	pterm.Error.Printf(fmt.Sprintf("%s", out))
		//	pterm.Error.Println(err)
		//	return
		//}
		//// Generate json
		//cmds = "go env -w GOPRIVATE=go.lib." + organisationID + ".alis.exchange,go.protobuf." + organisationID + ".alis.exchange,proto." + organisationID + ".alis.exchange,cli.alis.dev &&" +
		//	"protoc --plugin=protoc-gen-doc=$HOME/go/bin/protoc-gen-doc --doc_out=$HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " --doc_opt=json,docs.json -I=$HOME/alis.exchange/google/proto -I=$HOME/alis.exchange/" + organisationID + "/proto $(find $HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " -iname \"*.proto\")"
		//pterm.Debug.Printf("Shell command:\n%s\n", cmds)
		//out, err = exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		//if err != nil {
		//	pterm.Error.Printf(fmt.Sprintf("%s", out))
		//	pterm.Error.Println(err)
		//	return
		//}
		//
		////// Generate openapi description
		////cmds = "go env -w GOPRIVATE=go.lib." + organisationID + ".alis.exchange,go.protobuf." + organisationID + ".alis.exchange,proto." + organisationID + ".alis.exchange,cli.alis.dev &&" +
		////	"protoc --openapi_out=$HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " -I=$HOME/alis.exchange/google/proto -I=$HOME/alis.exchange/" + organisationID + "/proto $(find $HOME/alis.exchange/" + organisationID + "/proto/" + organisationID + "/" + productID + " -iname \"*.proto\")"
		////out, err = exec.CommandContext(cmd.Context(), "bash", "-c", cmds).CombinedOutput()
		////if err != nil {
		////	pterm.Error.Printf(fmt.Sprintf("%s", out))
		////	pterm.Error.Println(err)
		////	return
		////}
		//
		//if strings.Contains(fmt.Sprintf("%s", out), "warning") {
		//	pterm.Warning.Print(fmt.Sprintf("Generating documentation from protos...\n%s", out))
		//} else {
		//	pterm.Debug.Print(fmt.Sprintf("%s\n", out))
		//}
		//
		//pterm.Success.Printf("Generated documentation at %s\n", homeDir+"/alis.exchange/"+organisationID+"/products/"+productID)

		return

	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(protobufGenCmd)
	genCmd.AddCommand(descriptorGenCmd)
	genCmd.AddCommand(productDocsGenCmd)
	neuronCmd.SilenceUsage = true
	neuronCmd.SilenceErrors = true

	protobufGenCmd.Flags().BoolVarP(&pushProtocolBuffers, "publish", "p", false, pterm.Green("Generate the protocol buffers and push them to the protobuf repository"))
	protobufGenCmd.Flags().BoolVar(&pushPublicProtocolBuffers, "public", false, pterm.Green("Generate public protocol buffers and push them to the public protobuf repository."))

	productDocsGenCmd.Flags().BoolVar(&productsDocsGenCustomFlag, "custom", false, pterm.Green("Set custom visibility scopes for the documentation being generated."))
	productDocsGenCmd.Flags().BoolVar(&productsDocsGenPublicFlag, "public", false, pterm.Green("Set documentation visibility scope as public."))
	rootCmd.AddCommand(&cobra.Command{
		Use:    "iloveprotos",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			generateParsedText()
		},
	})

	// protobuf flags
	protobufGenCmd.Flags().BoolVar(&genprotoGo, "go", true, pterm.Green("Generate the protocol buffers for Golang"))
	protobufGenCmd.Flags().BoolVar(&genprotoPython, "python", false, pterm.Green("Generate the protocol buffers for Python"))
}
