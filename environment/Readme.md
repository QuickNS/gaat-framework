# Environment file generator

This code retrieves a variable group by ID from Azure DevOps and creates an output environment file that is used to run tests locally.

Usage:

```bash
go run generate_env.go -org https://dev.azure.com/<your_org> -p <project_name> -id <variable_group_id> -pat <token> --output variables.env
```

- org: the URL of your organization on Azure DevOps
- p: the project name
- id: the ID of the variable group. It's an integer value and can be obtained by the URL when browsing the Variable Group
- pat: a personal access token generated in Azure DevOps for api access
- output: the name of the generated file
