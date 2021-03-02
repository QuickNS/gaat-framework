# GART

Golang Azure Resources Testing framework

## What is GART?

GART is a framework written in [Go Programming Language](https://golang.org) with the goal of creating infrastructure deployment tests targeting the Azure platform. With GART we can easily create tests that validate the deployment of Azure resource and their configuration.

GART was built with IaC (infrastructure as code) projects in mind and its main goal is to enable automated testing after deployment of infrastructure happens in Azure (either through Terraform, ARM or any other technology).

- **Integration with Azure DevOps pipelines**:
- **Phased testing through the use of test tags**: allowing certain tests to run at specific stages of the pipeline ensuring we can stop deployment early in case core components are not setup correctly
- **Test data**: tbd
- **Kubernetes testing**: tbd
- **Azure KeyVault**: tbd
- **Test Results in JUnit format**: 

## Installing software

To use GART, it is required you have the following software on your development environment.

- Go (version 1.14.1 or up is recommended): https://golang.org/doc/install
- Git: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
- Azure CLI: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest

Verify your Go installation by running:
```batch
$ go version
go version go1.14.1 linux/amd64
```

If you intend to build tests that target AKS (Azure Kubernetes Service), make sure you add Kubectl (requires Azure CLI):

```
az aks install-c
```

## Add VS Code Go support

You can use any IDE to develop your tests, but here are instructions on how to configure VS Code:

- Install VS Code:  https://code.visualstudio.com/download

- Open VS Code and search for `golang.go` on the Extensions tab.

- After installing the extension you will be automatically prompted to install several `go` tools (linter, language server, etc). Make sure everything is working correctly:

```
Tools environment: GOPATH=/home/nunos/go
Installing 9 tools at /home/nunos/go/bin in module mode.
  gopkgs
  go-outline
  gotests
  gomodifytags
  impl
  goplay
  dlv
  golint
  gopls

  ...

  All tools successfully installed. You are ready to Go :).
```

## Setting up GART

Start by cloning this repo and execute the following command:

```bash
go mod tidy -v
```
This will make sure all the required packages are downloaded and installed.

Note: ocasionally, when using Git, you may run into issues after switching branches. If that happens, clean up the mod cache by doing:

```bash
	go clean -modcache
```
## How GART works

GART is built on top of the Golang test framework


Before you start developing your Azure infrastructure tests using GART, you first need to configure your development environment so you can run the tests locally.

