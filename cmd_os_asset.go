package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var osAssetsCmds = []Command{

	Command{
		Description:  "Lists available Assets",
		Subject:      "assets",
		AltSubject:   "assets",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list secrets", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Asset's usage"),
			}
		},
		ExecuteFunc: assetsListCmd,
		Endpoint:    ExtendedEndpoint,
	},
	Command{
		Description:  "Create asset",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create asset", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"filename":               c.FlagSet.String("filename", _nilDefaultStr, "Asset's filename"),
				"usage":                  c.FlagSet.String("usage", _nilDefaultStr, "Asset's usage. Possible values: \"bootloader\", \"ipxe_config_local_install\",\"ipxe_config_os_boot\",\"onie_installer\""),
				"mime":                   c.FlagSet.String("mime", _nilDefaultStr, "Required. Asset's mime type. Possible values: \"text/plain\",\"application/octet-stream\""),
				"url":                    c.FlagSet.String("url", _nilDefaultStr, "Asset's source url. If present it will not read content anymore"),
				"read_content_from_pipe": c.FlagSet.Bool("pipe", false, "Read secret's content read from pipe instead of terminal input"),
				"return_id":              c.FlagSet.Bool("return-id", false, "(Flag) If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: assetCreateCmd,
		Endpoint:    ExtendedEndpoint,
	},
	Command{
		Description:  "Delete asset",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete asset", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Asset's id or name"),
				"autoconfirm":      c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: assetDeleteCmd,
		Endpoint:    ExtendedEndpoint,
	},
	Command{
		Description:  "Add (associate) asset to template",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "associate",
		AltPredicate: "assign",
		FlagSet:      flag.NewFlagSet("associate template to asset", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name":    c.FlagSet.String("id", _nilDefaultStr, "Asset's id or filename"),
				"template_id_or_name": c.FlagSet.String("template-id", _nilDefaultStr, "Template's id or name"),
				"path":                c.FlagSet.String("path", _nilDefaultStr, "Path to associate asset to"),
			}
		},
		ExecuteFunc: associateAssetCmd,
		Endpoint:    ExtendedEndpoint,
	},
	Command{
		Description:  "List associated assets",
		Subject:      "assets",
		AltSubject:   "assets",
		Predicate:    "associated",
		AltPredicate: "template",
		FlagSet:      flag.NewFlagSet("associated assets", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Template's id or name"),
			}
		},
		ExecuteFunc: templateListAssociatedAssetsCmd,
		Endpoint:    ExtendedEndpoint,
	},
}

func assetsListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	list, err := client.OSAssets()

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 2,
		},
		SchemaField{
			FieldName: "FILENAME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "FILE_SIZE_BYTES",
			FieldType: TypeInt,
			FieldSize: 4,
		},
		SchemaField{
			FieldName: "FILE_MIME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "USAGE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "SOURCE_URL",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "CHECKSUM_SHA256",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		data = append(data, []interface{}{
			s.OSAssetID,
			s.OSAssetFileName,
			s.OSAssetFileSizeBytes,
			s.OSAssetFileMime,
			s.OSAssetUsage,
			s.OSAssetSourceURL,
			s.OSAssetContentsSHA256Hex,
		})

	}

	return renderTable("Assets", "", getStringParam(c.Arguments["format"]), data, schema)
}

func assetCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	obj := metalcloud.OSAsset{
		OSAssetFileName: getStringParam(c.Arguments["filename"]),
		OSAssetUsage:    getStringParam(c.Arguments["usage"]),
		OSAssetFileMime: getStringParam(c.Arguments["mime"]),
	}

	content := []byte{}

	if v, ok := getStringParamOk(c.Arguments["url"]); ok {
		obj.OSAssetSourceURL = v
	} else {

		if getBoolParam(c.Arguments["read_content_from_pipe"]) {
			_content, err := readInputFromPipe()
			if err != nil {
				return "", err
			}
			content = _content
		} else {
			_content, err := requestInputSilent("Asset content:")
			if err != nil {
				return "", err
			}
			content = _content
		}

		obj.OSAssetContentsBase64 = base64.StdEncoding.EncodeToString([]byte(content))
	}

	ret, err := client.OSAssetCreate(obj)

	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.OSAssetID), nil
	}

	return "", err
}

func assetDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retS, err := getOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Deleting asset  %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.OSAssetFileName,
			retS.OSAssetID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.OSAssetDelete(retS.OSAssetID)
	}

	return "", err
}

//asset_id_or_name
func getOSAssetFromCommand(paramName string, internalParamName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.OSAsset, error) {

	v, err := getParam(c, internalParamName, paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

	if isID {
		return client.OSAssetGet(id)
	}

	list, err := client.OSAssets()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.OSAssetFileName == label {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate secret with id/name %v", *v.(*interface{}))
}

func associateAssetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	asset, err := getOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	template, err := getOSTemplateFromCommand("template_id", c, client, false)
	if err != nil {
		return "", err
	}

	path := ""
	if v := c.Arguments["path"]; v != nil && *v.(*string) != _nilDefaultStr {
		path = *v.(*string)
	} else {
		return "", fmt.Errorf("path is required")
	}

	return "", client.OSTemplateAddOSAsset(template.VolumeTemplateID, asset.OSAssetID, path)
}

func templateListAssociatedAssetsCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	ret, err := getOSTemplateFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}

	list, err := client.OSTemplateOSAssets(ret.VolumeTemplateID)

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "PATH",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 2,
		},
		SchemaField{
			FieldName: "FILENAME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "FILE_SIZE_BYTES",
			FieldType: TypeInt,
			FieldSize: 4,
		},
		SchemaField{
			FieldName: "FILE_MIME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "USAGE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "SOURCE_URL",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}
	for path, s := range *list {

		data = append(data, []interface{}{
			path,
			s.OSAssetID,
			s.OSAssetFileName,
			s.OSAssetFileSizeBytes,
			s.OSAssetFileMime,
			s.OSAssetUsage,
			s.OSAssetSourceURL,
		})

	}

	var sb strings.Builder

	format := c.Arguments["format"]
	if format == nil {
		var f string
		f = ""
		format = &f
	}

	switch *format.(*string) {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:
		sb.WriteString(fmt.Sprintf("Assets associated to template (%s #%d)\n", ret.VolumeTemplateLabel, ret.VolumeTemplateID))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d assets\n\n", len(*list)))
	}

	return sb.String(), nil
}