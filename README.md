# The alis.exchange Command Line Interface

The `alis` command-line interface (CLI) is the primary tool to create and manage resources on **alis.exchange**.  You can use this tool to perform many common platform tasks either from the command line or in scripts and other automations.

For example, you can use the `alis` CLI to:

* list organisations and/or products
* clone a product to your local environment 
* create a new product / organisation
* deploy new versions of your product
* manage the build and deploy steps of your services

## Table of Contents

- [Prerequisites](https://github.com/alis-exchange/cli/blob/main/README.md#prerequisites)
- [Installation](https://github.com/alis-exchange/cli/blob/main/README.md#installation)
- [Try it out](https://github.com/alis-exchange/cli/blob/main/README.md#try-it-out)

## Prerequisites

The CLI requires the following to be set up in order to run.

### Google Cloud SDK

The CLI makes use of Google Cloud SDK authentication to seamlessly authenticate your requests to alis.exchange.  

1. Run `gcloud auth login` from your terminal to authenticate your local environment with Google user account via a web-based authorization flow.
2. Run `gcloud auth application-default login` to acquire new user credentials to use for Application Default Credentials ([ADC](https://developers.google.com/identity/protocols/application-default-credentials)). These are used in calling Google APIs.


### Go

Install any one of the **three latest major** [releases of Go](https://golang.org/doc/devel/release.html).  For installation instructions, see Go’s [Getting Started](https://golang.org/doc/install) guide.

_Check_
After installation, running `go version` should reflect one of the three latest major Go versions.

### Protocol Buffer compiler

1. Install the **[Protocol buffer](https://developers.google.com/protocol-buffers) compiler**, `protoc`, [version 3](https://developers.google.com/protocol-buffers/docs/proto3). For installation instructions, see [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/).  This tool significantly simplifies working with our Protocol Buffers.

2. Install the required **Go plugins** for the protocol compiler:

    1. Install the protocol compiler plugins for Go using the following commands:

            go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
            go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

    2. Update your `PATH` so that the `protoc` compiler can find the plugins:

           export PATH="$PATH:$(go env GOPATH)/bin"
           
### Git Configuration variables

Since the CLI is in a private repository, you will need to ensure that:

- Your Git user credentials are consistent with your GitHub account that was granted access to the CLI. Ensure this by running:
    ```
    git config --global user.name "YOUR_GITHUB_USERNAME"
    git config --global user.email "YOUR_GITHUB_EMAIL"
    ```
- You access the private repository with a SSH request, rather than a HTTP request. 
    1. Generate a [new access token](https://github.com/settings/tokens/new) and set:
        - Note: alis.exchange
        - Expiration: No expiration
        - Scopes: Repo (Full control of private repositories)
    2. Run the following in your terminal:
        
            export GIT_USER="YOUR_GITHUB_USERNAME"
            export TOKEN="PASTE_THE_GENERATED_TOKEN_HERE"
            git config --global url."https://${GIT_USER}:${TOKEN}@github.com".insteadOf "https://github.com"
        
_Check_

Check if this was successful by running `git config -l`. The response should include:

    
    user.name="YOUR_GITHUB_USERNAME"
    user.email="YOUR_GITHUB_EMAIL"
    url.https://{YOUR_GITHUB_USERNAME}:{GITHUB_TOKEN}@github.com.insteadof=https://github.com
    

## Installation

1. Since the CLI is in a private repo, the GOPRIVATE variables need to be set.  If not set, the `go install` will try and retrieve the libraries from the public golang.com domain and fail with at 404 not found error.

```
go env -w GOPRIVATE=go.protobuf.alis.alis.exchange,github.com/alis-exchange/cli/alis
```
2. Install the CLI. This will place the CLI binary in your `$GOPATH/bin` folder.

```
go install github.com/alis-exchange/cli/alis@latest
```
3. Ensure that `$GOPATH/bin` has been added to your `$PATH` such that your terminal can access the `alis` CLI. The following command appends the path to the `.zshrc` file.

```
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
```
   **NOTE**
    
   This command should be added to the 

The above will install the `alis` binary in your `$GOPATH/bin` folder.

## Try it out

```bash
# Show help 
alis -h

# list available organisations
alis org list

# Setup your local environment for organisation 'foo'
alis org get foo

# list available products
alis product list foo

# Setup a particular product, say 'ab' in your local environment
alis product get foo.ab

# In order to run methods in your local environment, you need to generate a key file
alis product getkey foo.ab
```
