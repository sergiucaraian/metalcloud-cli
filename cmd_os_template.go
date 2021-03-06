package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var osTemplatesCmds = []Command{

	{
		Description:  "Lists available Templates",
		Subject:      "os-template",
		AltSubject:   "template",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list templates", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Template's usage"),
			}
		},
		ExecuteFunc: templatesListCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Create template",
		Subject:      "os-template",
		AltSubject:   "template",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"label":                        c.FlagSet.String("label", _nilDefaultStr, "(Required) Template's label"),
				"display_name":                 c.FlagSet.String("display-name", _nilDefaultStr, "(Required) Template's display name"),
				"size":                         c.FlagSet.Int("size", _nilDefaultInt, "Template's size (bytes)"),
				"local_disk_supported":         c.FlagSet.Bool("local-disk-supported", false, "Template supports local disk install. Default false"),
				"boot_methods_supported":       c.FlagSet.String("boot-methods-supported", _nilDefaultStr, "(Required) Template boot methods supported. Defaults to pxe_iscsi."),
				"boot_type":                    c.FlagSet.String("boot-type", _nilDefaultStr, "(Required) Template boot type. Possible values: 'uefi_only','legacy_only','hybrid' "),
				"description":                  c.FlagSet.String("description", _nilDefaultStr, "Template description"),
				"os_type":                      c.FlagSet.String("os-type", _nilDefaultStr, "(Required) Template operating system type. For example, Ubuntu or CentOS."),
				"os_version":                   c.FlagSet.String("os-version", _nilDefaultStr, "(Required) Template operating system version."),
				"os_architecture":              c.FlagSet.String("os-architecture", _nilDefaultStr, "(Required) Template operating system architecture.Possible values: none, unknown, x86, x86_64."),
				"initial_user":                 c.FlagSet.String("initial-user", _nilDefaultStr, "(Required) Template's initial username, used to verify install."),
				"initial_password":             c.FlagSet.String("initial-password", _nilDefaultStr, "(Required) Template's initial password, used to verify install."),
				"initial_ssh_port":             c.FlagSet.Int("initial-ssh-port", _nilDefaultInt, "(Required) Template's initial ssh port, used to verify install."),
				"change_password_after_deploy": c.FlagSet.Bool("change-password-after-deploy", false, "Option to change the initial_user password on the installed OS after deploy."),
				"repo_url":                     c.FlagSet.String("repo-url", _nilDefaultStr, "Template's location the repository"),
				"os_asset_id_bootloader_local_install_id_or_name": c.FlagSet.String("install-bootloader-asset", _nilDefaultStr, "Template's bootloader asset id during install"),
				"os_asset_id_bootloader_os_boot_id_or_name":       c.FlagSet.String("os-boot-bootloader-asset", _nilDefaultStr, "Template's bootloader asset id during regular server boot"),

				"return_id": c.FlagSet.Bool("return-id", false, "(Flag) If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: templateCreateCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Edit template",
		Subject:      "os-template",
		AltSubject:   "template",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("update template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name":    c.FlagSet.String("id", _nilDefaultStr, "(Required) Template's id or label"),
				"label":                  c.FlagSet.String("label", _nilDefaultStr, "(Required) Template's label"),
				"display_name":           c.FlagSet.String("display-name", _nilDefaultStr, "(Required) Template's display name"),
				"size":                   c.FlagSet.Int("size", _nilDefaultInt, "Template's size (bytes)"),
				"local_disk_supported":   c.FlagSet.Bool("local-disk-supported", false, "Template supports local disk install. Default false"),
				"boot_methods_supported": c.FlagSet.String("boot-methods-supported", _nilDefaultStr, "Template boot methods supported. Defaults to pxe_iscsi."),
				"boot_type":              c.FlagSet.String("boot-type", _nilDefaultStr, "(Required) Template boot type. Possible values: 'uefi_only','legacy_only','hybrid' "),
				"description":            c.FlagSet.String("description", _nilDefaultStr, "Template description"),
				"os_type":                c.FlagSet.String("os-type", _nilDefaultStr, "(Required) Template operating system type. For example, Ubuntu or CentOS."),
				"os_version":             c.FlagSet.String("os-version", _nilDefaultStr, "(Required) Template operating system version."),
				"os_architecture":        c.FlagSet.String("os-architecture", _nilDefaultStr, "(Required) Template operating system architecture.Possible values: none, unknown, x86, x86_64."),
				"os_asset_id_bootloader_local_install_id_or_name": c.FlagSet.String("install-bootloader-asset", _nilDefaultStr, "Template's bootloader asset id during install"),
				"os_asset_id_bootloader_os_boot_id_or_name":       c.FlagSet.String("os-boot-bootloader-asset", _nilDefaultStr, "Template's bootloader asset id during regular server boot"),
				"initial_user":                 c.FlagSet.String("initial-user", _nilDefaultStr, "(Required) Template's initial username, used to verify install."),
				"initial_password":             c.FlagSet.String("initial-password", _nilDefaultStr, "(Required) Template's initial password, used to verify install."),
				"initial_ssh_port":             c.FlagSet.Int("initial-ssh-port", _nilDefaultInt, "(Required) Template's initial ssh port, used to verify install."),
				"change_password_after_deploy": c.FlagSet.Bool("change-password-after-deploy", false, "Option to change the initial_user password on the installed OS after deploy."),
				"repo_url":                     c.FlagSet.String("repo-url", _nilDefaultStr, "Template description"),
			}
		},
		ExecuteFunc: templateEditCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Get template",
		Subject:      "os-template",
		AltSubject:   "template",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Asset's id or name"),
				"format":              c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"show_credentials":    c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the templates initial ssh credentials"),
			}
		},
		ExecuteFunc: templateGetCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Delete template",
		Subject:      "os-template",
		AltSubject:   "template",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Asset's id or name"),
				"autoconfirm":         c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: templateDeleteCmd,
		Endpoint:    ExtendedEndpoint,
	},
}

func templatesListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	list, err := client.OSTemplates()

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 2,
		},
		{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DESCRIPTION",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SIZE_MBYTES",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		{
			FieldName: "BOOT_METHODS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "OS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "INSTALL_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "OS_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "USER_ID",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "UPDATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		installBootloader := ""
		if s.OSAssetBootloaderLocalInstall != 0 {
			asset, err := client.OSAssetGet(s.OSAssetBootloaderLocalInstall)
			if err != nil {
				return "", err
			}
			installBootloader = asset.OSAssetFileName
		}
		osBootloader := ""
		if s.OSAssetBootloaderOSBoot != 0 {
			asset, err := client.OSAssetGet(s.OSAssetBootloaderOSBoot)
			if err != nil {
				return "", err
			}
			osBootloader = asset.OSAssetFileName
		}

		osData := ""

		if s.VolumeTemplateOperatingSystem != nil {
			os := *s.VolumeTemplateOperatingSystem
			osData = fmt.Sprintf("%s %s %s",
				os.OperatingSystemType,
				os.OperatingSystemVersion,
				os.OperatingSystemArchitecture)
		}

		data = append(data, []interface{}{
			s.VolumeTemplateID,
			s.VolumeTemplateLabel,
			s.VolumeTemplateDisplayName,
			s.VolumeTemplateDescription,
			s.VolumeTemplateSizeMBytes,
			s.VolumeTemplateBootMethodsSupported,
			osData,
			installBootloader,
			osBootloader,
			s.UserID,
			s.VolumeTemplateCreatedTimestamp,
			s.VolumeTemplateUpdatedTimestamp,
		})

	}

	TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	return renderTable("Templates", "", getStringParam(c.Arguments["format"]), data, schema)
}

func updateTemplateFromCommand(obj metalcloud.OSTemplate, c *Command, client interfaces.MetalCloudClient, checkRequired bool) (*metalcloud.OSTemplate, error) {

	if v := c.Arguments["label"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateLabel = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("label is required")
		}
	}

	if v := c.Arguments["display_name"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateDisplayName = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("display-name is required")
		}
	}

	if v := c.Arguments["size"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.VolumeTemplateSizeMBytes = *v.(*int)
	}

	if v := c.Arguments["local_disk_supported"]; v != nil && *v.(*bool) {
		obj.VolumeTemplateLocalDiskSupported = true
	}

	obj.VolumeTemplateIsOSTemplate = true

	if v := c.Arguments["boot_methods_supported"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateBootMethodsSupported = *v.(*string)
	}

	if v := c.Arguments["boot_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateBootType = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("boot-type is required")
		}
	}

	if v := c.Arguments["description"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateDescription = *v.(*string)
	}

	//OS Data
	if v := c.Arguments["os_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		vt := metalcloud.OperatingSystem{}
		obj.VolumeTemplateOperatingSystem = &vt
		obj.VolumeTemplateOperatingSystem.OperatingSystemType = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("os-type is required")
		}
	}

	if v := c.Arguments["os_version"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateOperatingSystem.OperatingSystemVersion = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("os-version is required")
		}
	}

	if v := c.Arguments["os_architecture"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateOperatingSystem.OperatingSystemArchitecture = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("os-architecture is required")
		}
	}

	//Boot options

	if v := c.Arguments["os_asset_id_bootloader_local_install_id_or_name"]; v != nil && *v.(*string) != _nilDefaultStr {
		localInstallAsset, err := getOSAssetFromCommand("install_bootloader_asset", "os_asset_id_bootloader_local_install_id_or_name", c, client)
		if err != nil {
			return nil, err
		}
		obj.OSAssetBootloaderLocalInstall = localInstallAsset.OSAssetID
	}

	if v := c.Arguments["os_asset_id_bootloader_os_boot_id_or_name"]; v != nil && *v.(*string) != _nilDefaultStr {
		osBootBootloaderAsset, err := getOSAssetFromCommand("os_boot_bootloader_asset", "os_asset_id_bootloader_os_boot_id_or_name", c, client)
		if err != nil {
			return nil, err
		}
		obj.OSAssetBootloaderOSBoot = osBootBootloaderAsset.OSAssetID
	}

	//Credentials
	if v := c.Arguments["initial_user"]; v != nil && *v.(*string) != _nilDefaultStr {
		creds := metalcloud.OSTemplateCredentials{}
		obj.OSTemplateCredentials = &creds
		obj.OSTemplateCredentials.OSTemplateInitialUser = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("initial-user is required")
		}
	}

	if v := c.Arguments["initial_password"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.OSTemplateCredentials.OSTemplateInitialPassword = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("initial-password is required")
		}
	}

	if v := c.Arguments["initial_ssh_port"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.OSTemplateCredentials.OSTemplateInitialSSHPort = *v.(*int)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("initial-ssh-port is required")
		}
	}

	if v := c.Arguments["change_password_after_deploy"]; v != nil && *v.(*bool) {
		obj.OSTemplateCredentials.OSTemplateChangePasswordAfterDeploy = true
	}

	if v := c.Arguments["repo_url"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateRepoURL = *v.(*string)
	}

	return &obj, nil
}

func templateCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	obj := metalcloud.OSTemplate{}
	updatedObj, err := updateTemplateFromCommand(obj, c, client, true)
	if err != nil {
		return "", err
	}

	ret, err := client.OSTemplateCreate(*updatedObj)
	if err != nil {
		return "", err
	}
	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", ret.VolumeTemplateID), nil
	}

	return "", err
}

func templateEditCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	obj, err := getOSTemplateFromCommand("id", c, client, false)
	newobj := metalcloud.OSTemplate{}
	updatedObj, err := updateTemplateFromCommand(newobj, c, client, false)
	if err != nil {
		return "", err
	}
	_, err = client.OSTemplateUpdate(obj.VolumeTemplateID, *updatedObj)
	return "", err
}

func templateDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retS, err := getOSTemplateFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}
	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting template %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.VolumeTemplateDisplayName,
			retS.VolumeTemplateID)

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

	err = client.OSTemplateDelete(retS.VolumeTemplateID)

	return "", err
}

func getOSTemplateFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient, decryptPasswd bool) (*metalcloud.OSTemplate, error) {

	v, err := getParam(c, "template_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

	if isID {
		return client.OSTemplateGet(id, decryptPasswd)
	}

	list, err := client.OSTemplates()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.VolumeTemplateLabel == label {
			return &s, nil
		}
	}

	if isID {
		return nil, fmt.Errorf("template %d not found", id)
	}

	return nil, fmt.Errorf("template %s not found", label)
}

func templateGetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	showCredentials := false
	if c.Arguments["show_credentials"] != nil && *c.Arguments["show_credentials"].(*bool) {
		showCredentials = true
	}

	template, err := getOSTemplateFromCommand("id", c, client, showCredentials)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 2,
		},
		{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DESCRIPTION",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SIZE_MBYTES",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		{
			FieldName: "BOOT_METHODS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "OS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "USER_ID",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		{
			FieldName: "INSTALL_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "OS_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "UPDATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	credentials := ""

	if showCredentials {

		schema = append(schema, SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: TypeString,
			FieldSize: 5,
		})

		credentials = fmt.Sprintf("user:%s (port %d) passwd:%s (change_password_after_install:%v)",
			template.OSTemplateCredentials.OSTemplateInitialUser,
			template.OSTemplateCredentials.OSTemplateInitialSSHPort,
			template.OSTemplateCredentials.OSTemplateInitialPassword,
			template.OSTemplateCredentials.OSTemplateChangePasswordAfterDeploy)

	}
	osDetails := ""

	if template.VolumeTemplateOperatingSystem != nil {
		os := *template.VolumeTemplateOperatingSystem
		osDetails = fmt.Sprintf("%s %s %s",
			os.OperatingSystemType,
			os.OperatingSystemVersion,
			os.OperatingSystemArchitecture)
	}

	installBootloader := ""
	if template.OSAssetBootloaderLocalInstall != 0 {
		asset, err := client.OSAssetGet(template.OSAssetBootloaderLocalInstall)
		if err != nil {
			return "", err
		}
		installBootloader = asset.OSAssetFileName
	}
	osBootloader := ""
	if template.OSAssetBootloaderOSBoot != 0 {
		asset, err := client.OSAssetGet(template.OSAssetBootloaderOSBoot)
		if err != nil {
			return "", err
		}
		osBootloader = asset.OSAssetFileName
	}

	data = append(data, []interface{}{
		template.VolumeTemplateID,
		template.VolumeTemplateLabel,
		template.VolumeTemplateDisplayName,
		template.VolumeTemplateDescription,
		template.VolumeTemplateSizeMBytes,
		template.VolumeTemplateBootMethodsSupported,
		osDetails,
		template.UserID,
		installBootloader,
		osBootloader,
		template.VolumeTemplateCreatedTimestamp,
		template.VolumeTemplateUpdatedTimestamp,
		credentials,
	})

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
		sb.WriteString(fmt.Sprintf("Template %s (%d)\n", template.VolumeTemplateLabel, template.VolumeTemplateID))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

	}

	return sb.String(), nil
}
