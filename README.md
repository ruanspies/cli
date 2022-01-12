## What?

----------

The `alis` command-line tool is the primary CLI tool to create and manage resources on **alis.exchange**.  You can use this tool to perform many common platform tasks either from the command line or in scripts and other automations.

For example, you can use the `alis` tool to:

* list organisations and/or products
* clone a product to your local environment 
* create a new product / organisation
* deploy new versions of your product
* manage the build and deploy steps of your services

## Prerequisites

----------

### 1: Google Cloud SDK

The CLI makes use of Google Cloud SDK authentication to seamlessly authenticate your requests to alis.exchange.  Run the following to authenticate your local environment with Google:

    gcloud auth login

### 2:  Go

Install any one of the **three latest major** [releases of Go](https://golang.org/doc/devel/release.html).  For installation instructions, see Go’s [Getting Started](https://golang.org/doc/install) guide.

### 3: Protocol Buffer compiler

1. Install the **[Protocol buffer](https://developers.google.com/protocol-buffers) compiler**, `protoc`, [version 3](https://developers.google.com/protocol-buffers/docs/proto3). For installation instructions, see [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/).  This tool significantly simplifies working with our Protocol Buffers.

2. Install the required **Go plugins** for the protocol compiler:

    1. Install the protocol compiler plugins for Go using the following commands:

           $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
           $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

    2. Update your `PATH` so that the `protoc` compiler can find the plugins:

           $ export PATH="$PATH:$(go env GOPATH)/bin"

## Installation

```bash
# since these are private libraries, we need to set the GOPRIVATE variables to take this into account.  If not set, the go install will try and retrieve the libraries from the public golang.com domain and fail with at 404 not found error.
go env -w GOPRIVATE=go.protobuf.alis.alis.exchange,github.com/alis-exchange/cli/alis
go install github.com/alis-exchange/cli/alis@latest
```

The above will have installed the `alis` binary in your `$GOPATH/bin` folder.

Ensure that you have 

```bash
# since this CLI is a private repository, you may need to set your access token globally to authenticate underlying git requests automatically
# generate a token here: https://github.com/settings/tokens
export USER = "YOUR_GIT_USERNAME"
export TOKEN = "YOUR_GITHUB_PERSONAL ACCESS TOKEN"
git config --global url."https://${USER}:${TOKEN}@github.com".insteadOf "https://github.com"
```

### Try it out

```bash
# Show help 
alis -h

# list available organisations
alis org list

# Setup your local environment for organisation 'myorg'
alis org get myorg

# list available products
alis product list myorg

# Setup a particular product, say 'ab' in your local environment
alis product get myorg.ab

# In order to run methods in your local environment, you need to generate a key file
alis product getkey myorg.ab
```