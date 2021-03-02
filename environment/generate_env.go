package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
)

func main() {

	// read arguments
	organizationPtr := flag.String("org", "", "url of the Azure DevOps organization")
	projectPtr := flag.String("p", "", "Azure DevOps project name")
	groupNamePtr := flag.String("v", "", "the name of the variable group")
	outputPtr := flag.String("output", "variables.env", "name of the output file to generate")
	patPtr := flag.String("pat", "", "a personal access token created in Azure DevOps")
	flag.Parse()
	organizationURL := *organizationPtr
	project := *projectPtr
	personalAccessToken := *patPtr
	groupName := *groupNamePtr
	outputFileName := *outputPtr

	if organizationURL == "" || project == "" || personalAccessToken == "" || groupName == "" || outputFileName == "" {
		flag.PrintDefaults()
		log.Fatal()
	}

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(organizationURL, personalAccessToken)
	ctx := context.Background()

	// Create a client to interact with the TaskAgent area
	client, err := taskagent.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	// Get variable group
	args := taskagent.GetVariableGroupsArgs{Project: &project, GroupName: &groupName}
	responseValue, err := client.GetVariableGroups(ctx, args)
	if err != nil {
		log.Fatal(err)
	}

	if *responseValue == nil || len(*responseValue) == 0 {
		log.Fatalf("Variable Group %s does not exist", groupName)
	}

	variableGroup := (*responseValue)[0]

	generateFile(*variableGroup.Variables, outputFileName)
}

func generateFile(variables map[string]interface{}, outputFileName string) {

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	//create file
	filename := path.Join(currentDir, outputFileName)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// create map with values and uppercase keys and substituting . for _
	vars := make(map[string]string)
	for k, v := range variables {
		vars[strings.Replace(strings.ToUpper(k), ".", "_", -1)] = v.(map[string]interface{})["value"].(string)
	}

	// get list of keys for sorting
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}

	// sort keys
	sort.Strings(keys)

	// dump to file
	for _, k := range keys {
		f.WriteString(fmt.Sprintf("%s=%s\n", k, vars[k]))
	}
	fmt.Println("Generated file", filename)
}
