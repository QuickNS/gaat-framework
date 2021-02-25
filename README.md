# GART

Golang Azure Resources Testing framework

## What is GART?

GART is a framework written in [Go Programming Language](https://golang.org) with the goal of creating infrastructure deployment tests targeting the Azure platform. With GART we can easily create tests that validate Azure resource deployments and configuration which can be used to validate a desired state configuration on Azure.

GART was built with IaC (infrastructure as code) projects in mind and its main goal is to enable automated testing after deployment of infrastructure happens in Azure (either through Terraform, ARM or any other technology).

- **Integration with Azure DevOps pipelines**:
- **Phased testing through the use of test tags**: allowing certain tests to run at specific stages of the pipeline ensuring we can stop deployment early in case core components are not setup correctly
- **Test data**: tbd
- **Kubernetes testing**: tbd
- **Azure KeyVault**: tbd
- **Test Results in JUnit format**: 
