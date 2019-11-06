package main

import (
	"encoding/json"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/mock"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestInfrastructureRevertCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia-edited",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	autoconfirm := true
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id": &infra.InfrastructureID,
			"autoconfirm":       &autoconfirm,
		},
	}

	client.EXPECT().
		InfrastructureOperationCancel(infra.InfrastructureID).
		Return(nil).
		Times(1)

	ret, err := infrastructureRevertCmd(&cmd, client)

	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil())

}

func TestInfrastructureDeployCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia-edited",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	bFalse := false
	bTrue := true
	timeout := 256
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id":             &infra.InfrastructureID,
			"allow_data_loss":               &bFalse,
			"hard_shutdown_after_timeout":   &bTrue,
			"attempt_soft_shutdown":         &bFalse,
			"soft_shutdown_timeout_seconds": &timeout,
		},
	}

	expectedShutdownOptions := metalcloud.ShutdownOptions{
		HardShutdownAfterTimeout:   true,
		AttemptSoftShutdown:        false,
		SoftShutdownTimeoutSeconds: timeout,
	}

	client.EXPECT().
		InfrastructureDeploy(infra.InfrastructureID, expectedShutdownOptions, false, false).
		Return(nil).
		Times(1)

	//test first without confirmation
	ret, err := infrastructureDeleteCmd(&cmd, client)
	Expect(ret).To(Equal(""))
	Expect(err).NotTo(BeNil()) //should throw error indicating confirmation not given
	Expect(err.Error()).To(Equal("Operation not confirmed. Aborting"))

	cmd.Arguments["autoconfirm"] = &bTrue

	ret, err = infrastructureDeployCmd(&cmd, client)
	Expect(ret).To(Equal(""))
	Expect(err).To(BeNil()) //should be nil

}

func TestInfrastructureListCmd(t *testing.T) {
	RegisterTestingT(t)

	responseBody = `{"result": ` + _infrastructuresFixture1 + `,"jsonrpc": "2.0","id": 0}`

	client, err := metalcloud.GetMetalcloudClient("user", "APIKey", httpServer.URL, false)
	Expect(err).To(BeNil())

	format := "json"
	cmd := Command{
		Arguments: map[string]interface{}{
			"format": &format,
		},
	}

	ret, err1 := infrastructureListCmd(&cmd, client)

	Expect(err1).To(BeNil())

	reqBody := (<-requestChan).body
	Expect(reqBody).NotTo(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())
	Expect(m).NotTo(BeNil())

}

func TestInfrastructureGetCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	infra2 := metalcloud.Infrastructure{
		InfrastructureID:    10003,
		InfrastructureLabel: "testinfra2",
	}

	iao := metalcloud.InstanceArrayOperation{
		InstanceArrayID:           11,
		InstanceArrayLabel:        "testia-edited",
		InstanceArrayDeployType:   "edit",
		InstanceArrayDeployStatus: "not_started",
	}

	ia := metalcloud.InstanceArray{
		InstanceArrayID:            11,
		InstanceArrayLabel:         "testia",
		InfrastructureID:           infra.InfrastructureID,
		InstanceArrayOperation:     &iao,
		InstanceArrayServiceStatus: "active",
	}

	dao := metalcloud.DriveArrayOperation{
		DriveArrayID:           10,
		DriveArrayLabel:        "test-edited",
		InstanceArrayID:        ia.InstanceArrayID,
		InfrastructureID:       infra.InfrastructureID,
		DriveArrayCount:        101,
		DriveArrayDeployType:   "edit",
		DriveArrayDeployStatus: "not_started",
	}

	da := metalcloud.DriveArray{
		DriveArrayID:            10,
		DriveArrayLabel:         "test",
		InstanceArrayID:         ia.InstanceArrayID,
		InfrastructureID:        infra.InfrastructureID,
		DriveArrayCount:         101,
		DriveArrayOperation:     &dao,
		DriveArrayServiceStatus: "active",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(infra2.InfrastructureID).
		Return(&infra2, nil).
		AnyTimes()

	client.EXPECT().
		InstanceArrayGet(ia.InstanceArrayID).
		Return(&ia, nil).
		AnyTimes()

	format := "json"
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id": &infra.InfrastructureID,
			"format":            &format,
		},
	}

	iaList := map[string]metalcloud.InstanceArray{
		ia.InstanceArrayLabel + ".vanilla": ia,
	}

	client.EXPECT().
		InstanceArrays(gomock.Any()).
		Return(&iaList, nil).
		AnyTimes()

	daList := map[string]metalcloud.DriveArray{
		da.DriveArrayLabel + ".vanilla": da,
	}
	client.EXPECT().
		DriveArrays(gomock.Any()).
		Return(&daList, nil).
		AnyTimes()

	ret, err := infrastructureGetCmd(&cmd, client)

	Expect(ret).To(Not(Equal("")))
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(r["STATUS"].(string)).To(Equal("edited"))
	Expect(r["LABEL"].(string)).To(Equal(iao.InstanceArrayLabel))

	//test with label instead of id

	infraList := map[string]metalcloud.Infrastructure{
		infra.InfrastructureLabel:  infra,
		infra2.InfrastructureLabel: infra2,
	}

	client.EXPECT().
		Infrastructures().
		Return(&infraList, nil).
		AnyTimes()

	cmd = Command{
		Arguments: map[string]interface{}{
			"infrastructure_label": &infra.InfrastructureLabel,
			"format":               &format,
		},
	}

	ret, err = infrastructureGetCmd(&cmd, client)
	Expect(err).To(BeNil())

	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r = m[0].(map[string]interface{})
	Expect(r["STATUS"].(string)).To(Equal("edited"))
	Expect(r["LABEL"].(string)).To(Equal(iao.InstanceArrayLabel))

}

func TestGetInfrastructureIDFromCommand(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	infra := metalcloud.Infrastructure{
		InfrastructureID:    10002,
		InfrastructureLabel: "testinfra",
	}

	infra2 := metalcloud.Infrastructure{
		InfrastructureID:    10003,
		InfrastructureLabel: "testinfra2",
	}

	infra3 := metalcloud.Infrastructure{
		InfrastructureID:    10004,
		InfrastructureLabel: "testinfra",
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		InfrastructureGet(infra.InfrastructureID).
		Return(&infra, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(infra2.InfrastructureID).
		Return(&infra2, nil).
		AnyTimes()

	client.EXPECT().
		InfrastructureGet(infra3.InfrastructureID).
		Return(&infra2, nil).
		AnyTimes()

	infraListAmbigous := map[string]metalcloud.Infrastructure{
		infra.InfrastructureLabel:        infra,
		infra2.InfrastructureLabel:       infra2,
		infra3.InfrastructureLabel + "1": infra3,
	}

	client.EXPECT().
		Infrastructures().
		Return(&infraListAmbigous, nil).
		AnyTimes()

	//check with id
	cmd := Command{
		Arguments: map[string]interface{}{
			"infrastructure_id": &infra.InfrastructureID,
		},
	}

	infrastructureID, err := getInfrastructureIDFromCommand(&cmd, client)
	Expect(err).To(BeNil())
	Expect(infrastructureID).To(Equal(infra.InfrastructureID))

	//check with ambiguous label
	cmd = Command{
		Arguments: map[string]interface{}{
			"infrastructure_label": &infra.InfrastructureLabel,
		},
	}

	infrastructureID, err = getInfrastructureIDFromCommand(&cmd, client)
	Expect(err).NotTo(BeNil())

	//check with wrong label
	blablah := "asdasdasdasd"
	cmd.Arguments["infrastructure_label"] = &blablah

	infrastructureID, err = getInfrastructureIDFromCommand(&cmd, client)
	Expect(err).NotTo(BeNil())

	//check with correct label
	cmd.Arguments["infrastructure_label"] = &infra2.InfrastructureLabel

	infrastructureID, err = getInfrastructureIDFromCommand(&cmd, client)
	Expect(err).To(BeNil())
	Expect(infrastructureID).To(Equal(infra2.InfrastructureID))

}

const _InstanceArraysFixture1 = "{\"workers.vanilla\":{\"instance_array_id\":35516,\"instance_array_instance_count\":2,\"instance_array_ipv4_subnet_create_auto\":true,\"instance_array_ip_allocate_auto\":true,\"instance_array_ram_gbytes\":1,\"instance_array_processor_count\":1,\"instance_array_processor_core_mhz\":1000,\"instance_array_processor_core_count\":1,\"infrastructure_id\":25524,\"instance_array_service_status\":\"active\",\"instance_array_change_id\":215807,\"instance_array_label\":\"workers\",\"instance_array_subdomain\":\"workers.vanilla.complex-demo.7.bigstep.io\",\"drive_array_id_boot\":45928,\"instance_array_gui_settings_json\":\"{\\\"nRowIndex\\\":0,\\\"nColumnIndex\\\":3,\\\"bShowWidgetChildren\\\":true,\\\"randomInstanceID\\\":\\\"rand:0.6337124950169671\\\",\\\"userAgent\\\":\\\"Mozilla\\\\/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit\\\\/605.1.15 (KHTML, like Gecko) Version\\\\/13.0.2 Safari\\\\/605.1.15\\\"}\",\"instance_array_updated_timestamp\":\"2019-10-11T07:18:57Z\",\"instance_array_created_timestamp\":\"2019-03-28T15:23:18Z\",\"cluster_id\":40559,\"cluster_role_group\":\"none\",\"instance_array_disk_count\":0,\"instance_array_disk_size_mbytes\":0,\"instance_array_firewall_managed\":true,\"volume_template_id\":null,\"instance_array_boot_method\":\"pxe_iscsi\",\"instance_array_virtual_interfaces_enabled\":false,\"instance_array_operation\":{\"instance_array_change_id\":215807,\"instance_array_id\":35516,\"instance_array_instance_count\":2,\"instance_array_ipv4_subnet_create_auto\":true,\"instance_array_ip_allocate_auto\":true,\"instance_array_ram_gbytes\":1,\"instance_array_processor_count\":1,\"instance_array_processor_core_mhz\":1000,\"instance_array_processor_core_count\":1,\"instance_array_deploy_type\":\"edit\",\"instance_array_deploy_status\":\"finished\",\"instance_array_label\":\"workers\",\"instance_array_subdomain\":\"workers.vanilla.complex-demo.7.bigstep.io\",\"drive_array_id_boot\":45928,\"instance_array_gui_settings_json\":\"{\\\"nRowIndex\\\":0,\\\"nColumnIndex\\\":3,\\\"bShowWidgetChildren\\\":true,\\\"randomInstanceID\\\":\\\"rand:0.6337124950169671\\\",\\\"userAgent\\\":\\\"Mozilla\\\\/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit\\\\/605.1.15 (KHTML, like Gecko) Version\\\\/13.0.2 Safari\\\\/605.1.15\\\"}\",\"instance_array_updated_timestamp\":\"2019-10-11T07:18:57Z\",\"instance_array_disk_count\":0,\"instance_array_disk_size_mbytes\":0,\"instance_array_firewall_managed\":true,\"volume_template_id\":null,\"instance_array_boot_method\":\"pxe_iscsi\",\"instance_array_virtual_interfaces_enabled\":false,\"type\":\"InstanceArrayOperation\",\"instance_array_disk_types\":[],\"instance_array_firewall_rules\":[{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow SSH traffic.\",\"firewall_rule_port_range_end\":22,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":22,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 10.0.0.0/8.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"10.255.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"10.0.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 172.16.0.0/12.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"172.31.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"172.16.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 192.168.0.0/16.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"192.168.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"192.168.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv4 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv6 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv6\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow SNMP traffic.\",\"firewall_rule_port_range_end\":161,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":161,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"89.36.24.2\",\"firewall_rule_source_ip_address_range_start\":\"89.36.24.2\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"5.13.86.205\",\"firewall_rule_source_ip_address_range_start\":\"5.13.86.205\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null}],\"instance_array_interfaces\":[{\"instance_array_interface_change_id\":633555,\"instance_array_interface_id\":139705,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:19Z\",\"instance_array_interface_subdomain\":\"if0.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if0\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633554,\"instance_array_interface_id\":139706,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if1.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if1\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633552,\"instance_array_interface_id\":139707,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if2.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if2\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633562,\"instance_array_interface_id\":139708,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_subdomain\":\"if3.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if3\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]}]},\"type\":\"InstanceArray\",\"instance_array_disk_types\":[],\"instance_array_firewall_rules\":[{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow SSH traffic.\",\"firewall_rule_port_range_end\":22,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":22,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 10.0.0.0/8.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"10.255.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"10.0.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 172.16.0.0/12.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"172.31.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"172.16.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 192.168.0.0/16.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"192.168.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"192.168.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv4 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv6 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv6\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow SNMP traffic.\",\"firewall_rule_port_range_end\":161,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":161,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"89.36.24.2\",\"firewall_rule_source_ip_address_range_start\":\"89.36.24.2\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"5.13.86.205\",\"firewall_rule_source_ip_address_range_start\":\"5.13.86.205\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null}],\"instance_array_interfaces\":[{\"instance_array_interface_id\":139705,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633555,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:19Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if0\",\"instance_array_interface_subdomain\":\"if0.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633555,\"instance_array_interface_id\":139705,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:19Z\",\"instance_array_interface_subdomain\":\"if0.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if0\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139706,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633554,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if1\",\"instance_array_interface_subdomain\":\"if1.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633554,\"instance_array_interface_id\":139706,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if1.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if1\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139707,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633552,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if2\",\"instance_array_interface_subdomain\":\"if2.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633552,\"instance_array_interface_id\":139707,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if2.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if2\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139708,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633562,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if3\",\"instance_array_interface_subdomain\":\"if3.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633562,\"instance_array_interface_id\":139708,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_subdomain\":\"if3.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if3\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]}]},\"master.vanilla\":{\"instance_array_id\":35517,\"instance_array_instance_count\":1,\"instance_array_ipv4_subnet_create_auto\":true,\"instance_array_ip_allocate_auto\":true,\"instance_array_ram_gbytes\":1,\"instance_array_processor_count\":1,\"instance_array_processor_core_mhz\":1000,\"instance_array_processor_core_count\":1,\"infrastructure_id\":25524,\"instance_array_service_status\":\"active\",\"instance_array_change_id\":215806,\"instance_array_label\":\"master\",\"instance_array_subdomain\":\"master.vanilla.complex-demo.7.bigstep.io\",\"drive_array_id_boot\":45929,\"instance_array_gui_settings_json\":\"{\\\"nRowIndex\\\":0,\\\"nColumnIndex\\\":2,\\\"bShowWidgetChildren\\\":true,\\\"randomInstanceID\\\":\\\"rand:0.6337124950169671\\\",\\\"userAgent\\\":\\\"Mozilla\\\\/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit\\\\/605.1.15 (KHTML, like Gecko) Version\\\\/13.0.2 Safari\\\\/605.1.15\\\"}\",\"instance_array_updated_timestamp\":\"2019-10-11T07:18:20Z\",\"instance_array_created_timestamp\":\"2019-03-28T15:24:42Z\",\"cluster_id\":40559,\"cluster_role_group\":\"none\",\"instance_array_disk_count\":0,\"instance_array_disk_size_mbytes\":0,\"instance_array_firewall_managed\":true,\"volume_template_id\":null,\"instance_array_boot_method\":\"pxe_iscsi\",\"instance_array_virtual_interfaces_enabled\":false,\"instance_array_operation\":{\"instance_array_change_id\":215806,\"instance_array_id\":35517,\"instance_array_instance_count\":1,\"instance_array_ipv4_subnet_create_auto\":true,\"instance_array_ip_allocate_auto\":true,\"instance_array_ram_gbytes\":1,\"instance_array_processor_count\":1,\"instance_array_processor_core_mhz\":1000,\"instance_array_processor_core_count\":1,\"instance_array_deploy_type\":\"edit\",\"instance_array_deploy_status\":\"finished\",\"instance_array_label\":\"master\",\"instance_array_subdomain\":\"master.vanilla.complex-demo.7.bigstep.io\",\"drive_array_id_boot\":45929,\"instance_array_gui_settings_json\":\"{\\\"nRowIndex\\\":0,\\\"nColumnIndex\\\":2,\\\"bShowWidgetChildren\\\":true,\\\"randomInstanceID\\\":\\\"rand:0.6337124950169671\\\",\\\"userAgent\\\":\\\"Mozilla\\\\/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit\\\\/605.1.15 (KHTML, like Gecko) Version\\\\/13.0.2 Safari\\\\/605.1.15\\\"}\",\"instance_array_updated_timestamp\":\"2019-10-11T07:18:20Z\",\"instance_array_disk_count\":0,\"instance_array_disk_size_mbytes\":0,\"instance_array_firewall_managed\":true,\"volume_template_id\":null,\"instance_array_boot_method\":\"pxe_iscsi\",\"instance_array_virtual_interfaces_enabled\":false,\"type\":\"InstanceArrayOperation\",\"instance_array_disk_types\":[],\"instance_array_firewall_rules\":[{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow SSH traffic.\",\"firewall_rule_port_range_end\":22,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":22,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 10.0.0.0/8.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"10.255.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"10.0.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 172.16.0.0/12.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"172.31.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"172.16.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 192.168.0.0/16.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"192.168.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"192.168.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv4 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv6 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv6\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow SNMP traffic.\",\"firewall_rule_port_range_end\":161,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":161,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"89.36.24.2\",\"firewall_rule_source_ip_address_range_start\":\"89.36.24.2\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null}],\"instance_array_interfaces\":[{\"instance_array_interface_change_id\":633561,\"instance_array_interface_id\":139709,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:44Z\",\"instance_array_interface_subdomain\":\"if0.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if0\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633560,\"instance_array_interface_id\":139710,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:43Z\",\"instance_array_interface_subdomain\":\"if1.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if1\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633558,\"instance_array_interface_id\":139711,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:42Z\",\"instance_array_interface_subdomain\":\"if2.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if2\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633563,\"instance_array_interface_id\":139712,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_subdomain\":\"if3.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if3\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]}]},\"type\":\"InstanceArray\",\"instance_array_disk_types\":[],\"instance_array_firewall_rules\":[{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow SSH traffic.\",\"firewall_rule_port_range_end\":22,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":22,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 10.0.0.0/8.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"10.255.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"10.0.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 172.16.0.0/12.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"172.31.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"172.16.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 192.168.0.0/16.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"192.168.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"192.168.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv4 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv6 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv6\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow SNMP traffic.\",\"firewall_rule_port_range_end\":161,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":161,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"89.36.24.2\",\"firewall_rule_source_ip_address_range_start\":\"89.36.24.2\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null}],\"instance_array_interfaces\":[{\"instance_array_interface_id\":139709,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35517,\"instance_array_interface_change_id\":633561,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:44Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:24:42Z\",\"instance_array_interface_label\":\"if0\",\"instance_array_interface_subdomain\":\"if0.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633561,\"instance_array_interface_id\":139709,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:44Z\",\"instance_array_interface_subdomain\":\"if0.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if0\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139710,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35517,\"instance_array_interface_change_id\":633560,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:43Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:24:42Z\",\"instance_array_interface_label\":\"if1\",\"instance_array_interface_subdomain\":\"if1.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633560,\"instance_array_interface_id\":139710,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:43Z\",\"instance_array_interface_subdomain\":\"if1.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if1\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139711,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35517,\"instance_array_interface_change_id\":633558,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:42Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:24:42Z\",\"instance_array_interface_label\":\"if2\",\"instance_array_interface_subdomain\":\"if2.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633558,\"instance_array_interface_id\":139711,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:24:42Z\",\"instance_array_interface_subdomain\":\"if2.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if2\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139712,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35517,\"instance_array_interface_change_id\":633563,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:24:42Z\",\"instance_array_interface_label\":\"if3\",\"instance_array_interface_subdomain\":\"if3.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633563,\"instance_array_interface_id\":139712,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35517,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_subdomain\":\"if3.db.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if3\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]}]}}"

const _InstanceArrayGetfixture1 = "{\"instance_array_id\":35516,\"instance_array_instance_count\":2,\"instance_array_ipv4_subnet_create_auto\":true,\"instance_array_ip_allocate_auto\":true,\"instance_array_ram_gbytes\":1,\"instance_array_processor_count\":1,\"instance_array_processor_core_mhz\":1000,\"instance_array_processor_core_count\":1,\"infrastructure_id\":25524,\"instance_array_service_status\":\"active\",\"instance_array_change_id\":215807,\"instance_array_label\":\"workers\",\"instance_array_subdomain\":\"workers.vanilla.complex-demo.7.bigstep.io\",\"drive_array_id_boot\":45928,\"instance_array_gui_settings_json\":\"{\\\"nRowIndex\\\":0,\\\"nColumnIndex\\\":3,\\\"bShowWidgetChildren\\\":true,\\\"randomInstanceID\\\":\\\"rand:0.6337124950169671\\\",\\\"userAgent\\\":\\\"Mozilla\\\\/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit\\\\/605.1.15 (KHTML, like Gecko) Version\\\\/13.0.2 Safari\\\\/605.1.15\\\"}\",\"instance_array_updated_timestamp\":\"2019-10-11T07:18:57Z\",\"instance_array_created_timestamp\":\"2019-03-28T15:23:18Z\",\"cluster_id\":40559,\"cluster_role_group\":\"none\",\"instance_array_disk_count\":0,\"instance_array_disk_size_mbytes\":0,\"instance_array_firewall_managed\":true,\"volume_template_id\":null,\"instance_array_boot_method\":\"pxe_iscsi\",\"instance_array_virtual_interfaces_enabled\":false,\"instance_array_operation\":{\"instance_array_change_id\":215807,\"instance_array_id\":35516,\"instance_array_instance_count\":2,\"instance_array_ipv4_subnet_create_auto\":true,\"instance_array_ip_allocate_auto\":true,\"instance_array_ram_gbytes\":1,\"instance_array_processor_count\":1,\"instance_array_processor_core_mhz\":1000,\"instance_array_processor_core_count\":1,\"instance_array_deploy_type\":\"edit\",\"instance_array_deploy_status\":\"finished\",\"instance_array_label\":\"workers\",\"instance_array_subdomain\":\"workers.vanilla.complex-demo.7.bigstep.io\",\"drive_array_id_boot\":45928,\"instance_array_gui_settings_json\":\"{\\\"nRowIndex\\\":0,\\\"nColumnIndex\\\":3,\\\"bShowWidgetChildren\\\":true,\\\"randomInstanceID\\\":\\\"rand:0.6337124950169671\\\",\\\"userAgent\\\":\\\"Mozilla\\\\/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit\\\\/605.1.15 (KHTML, like Gecko) Version\\\\/13.0.2 Safari\\\\/605.1.15\\\"}\",\"instance_array_updated_timestamp\":\"2019-10-11T07:18:57Z\",\"instance_array_disk_count\":0,\"instance_array_disk_size_mbytes\":0,\"instance_array_firewall_managed\":true,\"volume_template_id\":null,\"instance_array_boot_method\":\"pxe_iscsi\",\"instance_array_virtual_interfaces_enabled\":false,\"type\":\"InstanceArrayOperation\",\"instance_array_disk_types\":[],\"instance_array_firewall_rules\":[{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow SSH traffic.\",\"firewall_rule_port_range_end\":22,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":22,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 10.0.0.0/8.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"10.255.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"10.0.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 172.16.0.0/12.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"172.31.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"172.16.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 192.168.0.0/16.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"192.168.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"192.168.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv4 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv6 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv6\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow SNMP traffic.\",\"firewall_rule_port_range_end\":161,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":161,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"89.36.24.2\",\"firewall_rule_source_ip_address_range_start\":\"89.36.24.2\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"5.13.86.205\",\"firewall_rule_source_ip_address_range_start\":\"5.13.86.205\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null}],\"instance_array_interfaces\":[{\"instance_array_interface_change_id\":633555,\"instance_array_interface_id\":139705,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:19Z\",\"instance_array_interface_subdomain\":\"if0.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if0\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633554,\"instance_array_interface_id\":139706,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if1.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if1\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633552,\"instance_array_interface_id\":139707,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if2.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if2\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_change_id\":633562,\"instance_array_interface_id\":139708,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_subdomain\":\"if3.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if3\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]}]},\"type\":\"InstanceArray\",\"instance_array_disk_types\":[],\"instance_array_firewall_rules\":[{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow SSH traffic.\",\"firewall_rule_port_range_end\":22,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":22,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow RDP traffic.\",\"firewall_rule_port_range_end\":3389,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":3389,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 10.0.0.0/8.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"10.255.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"10.0.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 172.16.0.0/12.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"172.31.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"172.16.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"all\",\"firewall_rule_description\":\"Allow traffic on 192.168.0.0/16.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":\"192.168.255.255\",\"firewall_rule_destination_ip_address_range_start\":\"192.168.0.0\"},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv4 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"icmp\",\"firewall_rule_description\":\"Allow IPv6 ICMP traffic.\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv6\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":false,\"firewall_rule_protocol\":\"udp\",\"firewall_rule_description\":\"Allow SNMP traffic.\",\"firewall_rule_port_range_end\":161,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":161,\"firewall_rule_source_ip_address_range_end\":null,\"firewall_rule_source_ip_address_range_start\":null,\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"89.36.24.2\",\"firewall_rule_source_ip_address_range_start\":\"89.36.24.2\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null},{\"type\":\"FirewallRule\",\"firewall_rule_enabled\":true,\"firewall_rule_protocol\":\"tcp\",\"firewall_rule_description\":\"Rule description\",\"firewall_rule_port_range_end\":null,\"firewall_rule_ip_address_type\":\"ipv4\",\"firewall_rule_port_range_start\":null,\"firewall_rule_source_ip_address_range_end\":\"5.13.86.205\",\"firewall_rule_source_ip_address_range_start\":\"5.13.86.205\",\"firewall_rule_destination_ip_address_range_end\":null,\"firewall_rule_destination_ip_address_range_start\":null}],\"instance_array_interfaces\":[{\"instance_array_interface_id\":139705,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633555,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:19Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if0\",\"instance_array_interface_subdomain\":\"if0.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633555,\"instance_array_interface_id\":139705,\"network_id\":58438,\"instance_array_interface_index\":0,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:19Z\",\"instance_array_interface_subdomain\":\"if0.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if0\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139706,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633554,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if1\",\"instance_array_interface_subdomain\":\"if1.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633554,\"instance_array_interface_id\":139706,\"network_id\":58437,\"instance_array_interface_index\":1,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if1.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if1\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139707,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633552,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if2\",\"instance_array_interface_subdomain\":\"if2.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633552,\"instance_array_interface_id\":139707,\"network_id\":null,\"instance_array_interface_index\":2,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_subdomain\":\"if2.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if2\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]},{\"instance_array_interface_id\":139708,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35516,\"instance_array_interface_change_id\":633562,\"instance_array_interface_service_status\":\"active\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_created_timestamp\":\"2019-03-28T15:23:18Z\",\"instance_array_interface_label\":\"if3\",\"instance_array_interface_subdomain\":\"if3.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_operation\":{\"instance_array_interface_change_id\":633562,\"instance_array_interface_id\":139708,\"network_id\":58439,\"instance_array_interface_index\":3,\"instance_array_id\":35516,\"instance_array_interface_deploy_type\":\"create\",\"instance_array_interface_deploy_status\":\"finished\",\"instance_array_interface_updated_timestamp\":\"2019-03-28T15:25:10Z\",\"instance_array_interface_subdomain\":\"if3.frontend.vanilla.complex-demo.7.bigstep.io\",\"instance_array_interface_label\":\"if3\",\"type\":\"InstanceArrayInterfaceOperation\",\"instance_array_interface_lagg_indexes\":[]},\"type\":\"InstanceArrayInterface\",\"instance_array_interface_lagg_indexes\":[]}]}"

const _infrastructuresFixture1 = "{\"my-terraform-infra4\":{\"infrastructure_id\":3553,\"datacenter_name\":\"us-santaclara\",\"user_id_owner\":2,\"infrastructure_label\":\"my-terraform-infra4\",\"infrastructure_created_timestamp\":\"2019-10-16T21:31:18Z\",\"infrastructure_subdomain\":\"my-terraform-infra4.2.poc.metalcloud.io\",\"infrastructure_change_id\":7617,\"infrastructure_service_status\":\"ordered\",\"infrastructure_touch_unixtime\":\"1571460507.3856\",\"infrastructure_updated_timestamp\":\"2019-10-19T04:48:27Z\",\"infrastructure_gui_settings_json\":\"\",\"infrastructure_private_datacenters_json\":null,\"infrastructure_deploy_id\":null,\"infrastructure_design_is_locked\":false,\"infrastructure_operation\":{\"infrastructure_change_id\":7617,\"infrastructure_id\":3553,\"datacenter_name\":\"us-santaclara\",\"user_id_owner\":2,\"infrastructure_label\":\"my-terraform-infra4\",\"infrastructure_subdomain\":\"my-terraform-infra4.2.poc.metalcloud.io\",\"infrastructure_deploy_type\":\"create\",\"infrastructure_deploy_status\":\"not_started\",\"infrastructure_updated_timestamp\":\"2019-10-19T04:48:27Z\",\"infrastructure_gui_settings_json\":\"\",\"infrastructure_private_datacenters_json\":null,\"infrastructure_deploy_id\":null,\"type\":\"InfrastructureOperation\",\"subnet_pool_lan\":null,\"infrastructure_reserved_lan_ip_ranges\":[]},\"type\":\"Infrastructure\",\"subnet_pool_lan\":null,\"infrastructure_reserved_lan_ip_ranges\":[],\"user_email_owner\":\"alex\"},\"my-terraform-infra5\":{\"infrastructure_id\":3574,\"datacenter_name\":\"us-santaclara\",\"user_id_owner\":2,\"infrastructure_label\":\"my-terraform-infra5\",\"infrastructure_created_timestamp\":\"2019-10-19T15:21:16Z\",\"infrastructure_subdomain\":\"my-terraform-infra5.2.poc.metalcloud.io\",\"infrastructure_change_id\":7618,\"infrastructure_service_status\":\"ordered\",\"infrastructure_touch_unixtime\":\"1571498477.2109\",\"infrastructure_updated_timestamp\":\"2019-10-19T15:21:16Z\",\"infrastructure_gui_settings_json\":\"\",\"infrastructure_private_datacenters_json\":null,\"infrastructure_deploy_id\":null,\"infrastructure_design_is_locked\":false,\"infrastructure_operation\":{\"infrastructure_change_id\":7618,\"infrastructure_id\":3574,\"datacenter_name\":\"us-santaclara\",\"user_id_owner\":2,\"infrastructure_label\":\"my-terraform-infra5\",\"infrastructure_subdomain\":\"my-terraform-infra5.2.poc.metalcloud.io\",\"infrastructure_deploy_type\":\"create\",\"infrastructure_deploy_status\":\"not_started\",\"infrastructure_updated_timestamp\":\"2019-10-19T15:21:16Z\",\"infrastructure_gui_settings_json\":\"\",\"infrastructure_private_datacenters_json\":null,\"infrastructure_deploy_id\":null,\"type\":\"InfrastructureOperation\",\"subnet_pool_lan\":null,\"infrastructure_reserved_lan_ip_ranges\":[]},\"type\":\"Infrastructure\",\"subnet_pool_lan\":null,\"infrastructure_reserved_lan_ip_ranges\":[],\"user_email_owner\":\"alex\"},\"my-terraform-infra6\":{\"infrastructure_id\":3576,\"datacenter_name\":\"us-santaclara\",\"user_id_owner\":2,\"infrastructure_label\":\"my-terraform-infra6\",\"infrastructure_created_timestamp\":\"2019-10-19T15:25:53Z\",\"infrastructure_subdomain\":\"my-terraform-infra6.2.poc.metalcloud.io\",\"infrastructure_change_id\":7619,\"infrastructure_service_status\":\"ordered\",\"infrastructure_touch_unixtime\":\"1571498753.6277\",\"infrastructure_updated_timestamp\":\"2019-10-19T15:25:53Z\",\"infrastructure_gui_settings_json\":\"\",\"infrastructure_private_datacenters_json\":null,\"infrastructure_deploy_id\":null,\"infrastructure_design_is_locked\":false,\"infrastructure_operation\":{\"infrastructure_change_id\":7619,\"infrastructure_id\":3576,\"datacenter_name\":\"us-santaclara\",\"user_id_owner\":2,\"infrastructure_label\":\"my-terraform-infra6\",\"infrastructure_subdomain\":\"my-terraform-infra6.2.poc.metalcloud.io\",\"infrastructure_deploy_type\":\"create\",\"infrastructure_deploy_status\":\"not_started\",\"infrastructure_updated_timestamp\":\"2019-10-19T15:25:53Z\",\"infrastructure_gui_settings_json\":\"\",\"infrastructure_private_datacenters_json\":null,\"infrastructure_deploy_id\":null,\"type\":\"InfrastructureOperation\",\"subnet_pool_lan\":null,\"infrastructure_reserved_lan_ip_ranges\":[]},\"type\":\"Infrastructure\",\"subnet_pool_lan\":null,\"infrastructure_reserved_lan_ip_ranges\":[],\"user_email_owner\":\"alex\"}}"
