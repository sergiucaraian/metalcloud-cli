package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var variablesCmds = []Command{

	{
		Description:  "Lists available variables",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list variables", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Variable's usage"),
			}
		},
		ExecuteFunc: variablesListCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Create variable",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create variable", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"name":                   c.FlagSet.String("name", _nilDefaultStr, "Variable's name"),
				"usage":                  c.FlagSet.String("usage", _nilDefaultStr, "Variable's usage"),
				"read_content_from_pipe": c.FlagSet.Bool("pipe", false, "Read variable's content read from pipe instead of terminal input"),
				"return_id":              c.FlagSet.Bool("return-id", false, "(Flag) If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: variableCreateCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Delete variable",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete variable", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"variable_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Variable's id or name"),
				"autoconfirm":         c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: variableDeleteCmd,
		Endpoint:    ExtendedEndpoint,
	},
}

func variablesListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	usage := *c.Arguments["usage"].(*string)
	if usage == _nilDefaultStr {
		usage = ""
	}

	list, err := client.Variables(usage)

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "USAGE",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "UPDATED",
			FieldType: TypeString,
			FieldSize: 20,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		data = append(data, []interface{}{
			s.VariableID,
			s.VariableName,
			s.VariableUsage,
			s.VariableCreatedTimestamp,
			s.VariableUpdatedTimestamp,
		})

	}

	TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	return renderTable("Variables", "", getStringParam(c.Arguments["format"]), data, schema)
}

func variableCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	variable := metalcloud.Variable{}

	if v := c.Arguments["name"]; v != nil && *v.(*string) != _nilDefaultStr {
		variable.VariableName = *v.(*string)
	} else {
		return "", fmt.Errorf("name is required")
	}

	if v := c.Arguments["usage"]; v != nil && *v.(*string) != _nilDefaultStr {
		variable.VariableUsage = *v.(*string)
	}

	var err error
	content := []byte{}
	if v := c.Arguments["read_content_from_pipe"]; *v.(*bool) {
		content, err = readInputFromPipe()
	} else {
		content, err = requestInput("Variable content:")
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	b, err := json.Marshal(content)
	variable.VariableJSON = string(b)

	ret, err := client.VariableCreate(variable)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", ret.VariableID), nil
	}

	return "", err
}

func variableDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retS, err := getVariableFromCommand("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting variable %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.VariableName,
			retS.VariableID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm, err = requestConfirmation(confirmationMessage)
		if err != nil {
			return "", err
		}

	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.VariableDelete(retS.VariableID)

	return "", err
}

func getVariableFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.Variable, error) {

	v, err := getParam(c, "variable_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

	if isID {
		return client.VariableGet(id)
	}

	variables, err := client.Variables("")
	if err != nil {
		return nil, err
	}

	for _, s := range *variables {
		if s.VariableName == label {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate variable with id/name %v", *v.(*interface{}))
}
