/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2022 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 *
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"fmt"

	"hcm/pkg/adaptor/types/core"
	"hcm/pkg/criteria/constant"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/criteria/validator"
)

// -------------------------- Create --------------------------

// SubnetCreateOption defines create subnet options.
type SubnetCreateOption[T SubnetCreateExt] struct {
	Name       string  `json:"name" validate:"required"`
	Memo       *string `json:"memo,omitempty" validate:"omitempty"`
	CloudVpcID string  `json:"cloud_vpc_id" validate:"required"`
	Extension  *T      `json:"extension" validate:"required"`
}

// SubnetCreateExt defines create subnet extensional info.
type SubnetCreateExt interface {
	TCloudSubnetCreateExt | AwsSubnetCreateExt | GcpSubnetCreateExt | AzureSubnetCreateExt | HuaWeiSubnetCreateExt
}

// TCloudSubnetCreateExt defines tencent cloud create subnet extensional info.
type TCloudSubnetCreateExt struct {
	Region   string `json:"region" validate:"required"`
	Zone     string `json:"zone" validate:"required"`
	IPv4Cidr string `json:"ipv4_cidr" validate:"required,cidrv4"`
}

// AwsSubnetCreateExt defines create aws subnet extensional info.
type AwsSubnetCreateExt struct {
	Region   string  `json:"region" validate:"required"`
	Zone     *string `json:"zone" validate:"omitempty"`
	IPv4Cidr *string `json:"ipv4_cidr" validate:"omitempty,cidrv4"`
	IPv6Cidr *string `json:"ipv6_cidr" validate:"omitempty,cidrv6"`
}

// GcpSubnetCreateExt defines create gcp subnet extensional info.
type GcpSubnetCreateExt struct {
	Region                string `json:"region" validate:"required"`
	IPv4Cidr              string `json:"ipv4_cidr" validate:"required,cidrv4"`
	PrivateIpGoogleAccess bool   `json:"private_ip_google_access" validate:"omitempty"`
	EnableFlowLogs        bool   `json:"enable_flow_logs" validate:"omitempty"`
}

// AzureSubnetCreateExt defines create azure subnet extensional info.
type AzureSubnetCreateExt struct {
	ResourceGroup        string   `json:"resource_group" validate:"required"`
	IPv4Cidr             []string `json:"ipv4_cidr" validate:"required,dive,cidrv4"`
	IPv6Cidr             []string `json:"ipv6_cidr" validate:"omitempty,dive,cidrv6"`
	CloudRouteTableID    string   `json:"cloud_route_table_id,omitempty" validate:"omitempty"`
	NatGateway           string   `json:"nat_gateway,omitempty" validate:"omitempty"`
	NetworkSecurityGroup string   `json:"network_security_group,omitempty" validate:"omitempty"`
}

// HuaWeiSubnetCreateExt defines create huawei subnet extensional info.
type HuaWeiSubnetCreateExt struct {
	Region     string  `json:"region" validate:"required"`
	Zone       *string `json:"zone" validate:"omitempty"`
	IPv4Cidr   string  `json:"ipv4_cidr" validate:"required,cidrv4"`
	Ipv6Enable bool    `json:"ipv6_enable" validate:"omitempty"`
	GatewayIp  string  `json:"gateway_ip" validate:"required"`
}

// TCloudSubnetCreateOption defines tcloud create subnet options.
type TCloudSubnetCreateOption SubnetCreateOption[TCloudSubnetCreateExt]

// Validate TCloudSubnetCreateOption.
func (c TCloudSubnetCreateOption) Validate() error {
	return validator.Validate.Struct(c)
}

// TCloudSubnetsCreateOption defines create tencent cloud subnets options.
type TCloudSubnetsCreateOption struct {
	AccountID  string                     `json:"account_id" validate:"required"`
	Region     string                     `json:"region" validate:"required"`
	CloudVpcID string                     `json:"cloud_vpc_id" validate:"required"`
	Subnets    []TCloudOneSubnetCreateOpt `json:"subnets" validate:"min=1,max=100"`
}

// TCloudOneSubnetCreateOpt defines create one tencent cloud subnets options for TCloudSubnetsCreateOption.
type TCloudOneSubnetCreateOpt struct {
	IPv4Cidr          string `json:"ipv4_cidr" validate:"required,cidrv4"`
	Name              string `json:"name" validate:"required"`
	Zone              string `json:"zone" validate:"required"`
	CloudRouteTableID string `json:"cloud_route_table_id" validate:"omitempty"`
}

// Validate TCloudSubnetsCreateOption.
func (c TCloudSubnetsCreateOption) Validate() error {
	return validator.Validate.Struct(c)
}

// AwsSubnetCreateOption defines aws create subnet options.
type AwsSubnetCreateOption SubnetCreateOption[AwsSubnetCreateExt]

// Validate AwsSubnetCreateOption.
func (c AwsSubnetCreateOption) Validate() error {
	return validator.Validate.Struct(c)
}

// AwsDefaultSubnetCreateOption defines create default aws subnet extensional info.
type AwsDefaultSubnetCreateOption struct {
	Region string `json:"region" validate:"required"`
	Zone   string `json:"zone" validate:"required"`
}

// Validate AwsDefaultSubnetCreateOption.
func (c AwsDefaultSubnetCreateOption) Validate() error {
	return validator.Validate.Struct(c)
}

// GcpSubnetCreateOption defines gcp create subnet options.
type GcpSubnetCreateOption SubnetCreateOption[GcpSubnetCreateExt]

// Validate GcpSubnetCreateOption.
func (c GcpSubnetCreateOption) Validate() error {
	return validator.Validate.Struct(c)
}

// AzureSubnetCreateOption defines azure create subnet options.
type AzureSubnetCreateOption SubnetCreateOption[AzureSubnetCreateExt]

// Validate AzureSubnetCreateOption.
func (c AzureSubnetCreateOption) Validate() error {
	return validator.Validate.Struct(c)
}

// HuaWeiSubnetCreateOption defines HuaWei create subnet options.
type HuaWeiSubnetCreateOption SubnetCreateOption[HuaWeiSubnetCreateExt]

// Validate HuaWeiSubnetCreateOption.
func (c HuaWeiSubnetCreateOption) Validate() error {
	return validator.Validate.Struct(c)
}

// -------------------------- Update --------------------------

// SubnetUpdateOption defines update subnet options.
type SubnetUpdateOption struct {
	ResourceID string                `json:"resource_id"`
	Data       *BaseSubnetUpdateData `json:"data"`
}

// BaseSubnetUpdateData defines the basic update subnet instance data.
type BaseSubnetUpdateData struct {
	Memo *string `json:"memo"`
}

// Validate BaseSubnetUpdateData.
func (s BaseSubnetUpdateData) Validate() error {
	if s.Memo == nil {
		return errf.New(errf.InvalidParameter, "memo is required")
	}
	return nil
}

// TCloudSubnetUpdateOption defines tencent cloud update subnet options.
type TCloudSubnetUpdateOption struct{}

// Validate TCloudSubnetUpdateOption.
func (s TCloudSubnetUpdateOption) Validate() error {
	return nil
}

// AwsSubnetUpdateOption defines aws update subnet options.
type AwsSubnetUpdateOption struct{}

// Validate AwsSubnetUpdateOption.
func (s AwsSubnetUpdateOption) Validate() error {
	return nil
}

// GcpSubnetUpdateOption defines gcp update subnet options.
type GcpSubnetUpdateOption struct {
	SubnetUpdateOption `json:",inline"`
	Region             string `json:"region"`
}

// Validate GcpSubnetUpdateOption.
func (s GcpSubnetUpdateOption) Validate() error {
	if len(s.ResourceID) == 0 {
		return errf.New(errf.InvalidParameter, "resource id is required")
	}

	if s.Data == nil {
		return errf.New(errf.InvalidParameter, "update data is required")
	}

	if err := s.Data.Validate(); err != nil {
		return err
	}

	return nil
}

// AzureSubnetUpdateOption defines azure update subnet options.
type AzureSubnetUpdateOption struct{}

// Validate AzureSubnetUpdateOption.
func (s AzureSubnetUpdateOption) Validate() error {
	return nil
}

// HuaWeiSubnetUpdateOption defines huawei update subnet options.
type HuaWeiSubnetUpdateOption struct {
	SubnetUpdateOption `json:",inline"`
	Region             string `json:"region"`
	Name               string `json:"name"`
	VpcID              string `json:"vpc_id"`
}

// Validate HuaWeiSubnetUpdateOption.
func (s HuaWeiSubnetUpdateOption) Validate() error {
	if err := s.Data.Validate(); err != nil {
		return err
	}

	if len(s.Region) == 0 {
		return errf.New(errf.InvalidParameter, "region is required")
	}

	if len(s.VpcID) == 0 {
		return errf.New(errf.InvalidParameter, "vpc id is required")
	}
	return nil
}

// ------------------------- Delete -------------------------

// AzureSubnetDeleteOption defines azure delete subnet options.
type AzureSubnetDeleteOption struct {
	core.AzureDeleteOption `json:",inline"`
	VpcID                  string `json:"vpc_id"`
}

// Validate AzureSubnetDeleteOption.
func (a AzureSubnetDeleteOption) Validate() error {
	if err := a.AzureDeleteOption.Validate(); err != nil {
		return err
	}

	if len(a.VpcID) == 0 {
		return errf.New(errf.InvalidParameter, "vpc id must be set")
	}

	return nil
}

// HuaWeiSubnetDeleteOption defines huawei delete subnet options.
type HuaWeiSubnetDeleteOption struct {
	core.BaseRegionalDeleteOption `json:",inline"`
	VpcID                         string `json:"vpc_id"`
}

// Validate HuaWeiSubnetDeleteOption.
func (s HuaWeiSubnetDeleteOption) Validate() error {
	if err := s.BaseRegionalDeleteOption.Validate(); err != nil {
		return err
	}

	if len(s.VpcID) == 0 {
		return errf.New(errf.InvalidParameter, "vpc id is required")
	}
	return nil
}

// -------------------------- List --------------------------

// TCloudSubnetListResult defines tencent cloud list subnet result.
type TCloudSubnetListResult struct {
	Count   *uint64        `json:"count,omitempty"`
	Details []TCloudSubnet `json:"details"`
}

// AwsSubnetListResult defines aws list subnet result.
type AwsSubnetListResult struct {
	NextToken *string     `json:"next_token,omitempty"`
	Details   []AwsSubnet `json:"details"`
}

// GcpSubnetListOption basic gcp list subnet options.
type GcpSubnetListOption struct {
	core.GcpListOption `json:",inline"`
	Region             string `json:"region"`
}

// Validate gcp list subnet option.
func (g GcpSubnetListOption) Validate() error {
	if err := g.GcpListOption.Validate(); err != nil {
		return err
	}

	if len(g.Region) == 0 {
		return errf.New(errf.InvalidParameter, "region can be empty")
	}

	return nil
}

// GcpSubnetListResult defines gcp list subnet result.
type GcpSubnetListResult struct {
	NextPageToken string      `json:"next_page_token,omitempty"`
	Details       []GcpSubnet `json:"details"`
}

// AzureSubnetListOption defines azure list subnet options.
type AzureSubnetListOption struct {
	core.AzureListOption `json:",inline"`
	CloudVpcID           string `json:"cloud_vpc_id"`
}

// Validate AzureSubnetListOption.
func (a AzureSubnetListOption) Validate() error {
	if err := a.AzureListOption.Validate(); err != nil {
		return err
	}

	if len(a.CloudVpcID) == 0 {
		return errf.New(errf.InvalidParameter, "vpc id must be set")
	}

	return nil
}

// AzureSubnetListByIDOption defines azure list subnet options.
type AzureSubnetListByIDOption struct {
	core.AzureListByIDOption `json:",inline"`
	CloudVpcID               string `json:"cloud_vpc_id"`
}

// Validate AzureSubnetListOption.
func (a AzureSubnetListByIDOption) Validate() error {
	if err := a.AzureListByIDOption.Validate(); err != nil {
		return err
	}

	if len(a.CloudVpcID) == 0 {
		return errf.New(errf.InvalidParameter, "cloud vpc id must be set")
	}

	return nil
}

// AzureSubnetListResult defines azure list subnet result.
type AzureSubnetListResult struct {
	Details []AzureSubnet `json:"details"`
}

// HuaWeiSubnetListOption defines huawei list subnet options.
type HuaWeiSubnetListOption struct {
	Region     string           `json:"region"`
	Page       *core.HuaWeiPage `json:"page,omitempty"`
	CloudVpcID string           `json:"vpc_id,omitempty"`
}

// Validate huawei list option.
func (s HuaWeiSubnetListOption) Validate() error {
	if len(s.Region) == 0 {
		return errf.New(errf.InvalidParameter, "region is required")
	}

	if s.Page != nil {
		if err := s.Page.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// HuaWeiSubnetListByIDOption ...
type HuaWeiSubnetListByIDOption struct {
	Region     string   `json:"region" validate:"required"`
	CloudVpcID string   `json:"cloud_vpc_id" validate:"required"`
	CloudIDs   []string `json:"cloud_ids" validate:"required"`
}

// Validate HuaWeiSubnetListByIDOption.
func (opt HuaWeiSubnetListByIDOption) Validate() error {
	if err := validator.Validate.Struct(opt); err != nil {
		return err
	}

	if len(opt.CloudIDs) > constant.BatchOperationMaxLimit {
		return fmt.Errorf("cloudIDs should <= %d", constant.BatchOperationMaxLimit)
	}

	return nil
}

// HuaWeiSubnetListResult defines huawei list subnet result.
type HuaWeiSubnetListResult struct {
	Details []HuaWeiSubnet `json:"details"`
}

// Subnet defines subnet struct.
type Subnet[T SubnetExtension] struct {
	// TODO: gcp 添加 vpcSelfLink字段，不要和 CloudVpcID 字段混用
	// CloudVpcID gcp 该字段为 self_link
	CloudVpcID string   `json:"cloud_vpc_id"`
	CloudID    string   `json:"cloud_id"`
	Name       string   `json:"name"`
	Ipv4Cidr   []string `json:"ipv4_cidr,omitempty"`
	Ipv6Cidr   []string `json:"ipv6_cidr,omitempty"`
	Memo       *string  `json:"memo,omitempty"`
	Extension  *T       `json:"extension"`
}

// SubnetExtension defines subnet extensional info.
type SubnetExtension interface {
	TCloudSubnetExtension | AwsSubnetExtension | GcpSubnetExtension | AzureSubnetExtension | HuaWeiSubnetExtension
}

// TCloudSubnetExtension defines tcloud subnet extensional info.
type TCloudSubnetExtension struct {
	IsDefault               bool    `json:"is_default"`
	Region                  string  `json:"region"`
	Zone                    string  `json:"zone"`
	CloudRouteTableID       *string `json:"cloud_route_table_id,omitempty"`
	CloudNetworkAclID       *string `json:"cloud_network_acl_id,omitempty"`
	AvailableIPAddressCount uint64  `json:"available_ip_address_count,omitempty"`
	TotalIpAddressCount     uint64  `json:"total_ip_address_count,omitempty"`
	UsedIpAddressCount      uint64  `json:"used_ip_address_count,omitempty"`
}

// AwsSubnetExtension defines aws subnet extensional info.
type AwsSubnetExtension struct {
	State                       string `json:"state"`
	Region                      string `json:"region"`
	Zone                        string `json:"zone"`
	IsDefault                   bool   `json:"is_default"`
	MapPublicIpOnLaunch         bool   `json:"map_public_ip_on_launch"`
	AssignIpv6AddressOnCreation bool   `json:"assign_ipv6_address_on_creation"`
	HostnameType                string `json:"hostname_type"`
	AvailableIPAddressCount     int64  `json:"available_ip_address_count"`
	TotalIpAddressCount         int64  `json:"total_ip_address_count"`
	UsedIpAddressCount          int64  `json:"used_ip_address_count"`
}

// GcpSubnetExtension defines gcp subnet extensional info.
type GcpSubnetExtension struct {
	SelfLink              string `json:"self_link"`
	Region                string `json:"region"`
	StackType             string `json:"stack_type"`
	Ipv6AccessType        string `json:"ipv6_access_type"`
	GatewayAddress        string `json:"gateway_address"`
	PrivateIpGoogleAccess bool   `json:"private_ip_google_access"`
	EnableFlowLogs        bool   `json:"enable_flow_logs"`
}

// AzureSubnetExtension defines azure subnet extensional info.
type AzureSubnetExtension struct {
	ResourceGroupName    string  `json:"resource_group"`
	CloudRouteTableID    *string `json:"cloud_route_table_id,omitempty"`
	NatGateway           string  `json:"nat_gateway,omitempty"`
	NetworkSecurityGroup string  `json:"network_security_group,omitempty"`
}

// HuaWeiSubnetExtension defines huawei subnet extensional info.
type HuaWeiSubnetExtension struct {
	Region       string   `json:"region"`
	Status       string   `json:"status"`
	DhcpEnable   bool     `json:"dhcp_enable"`
	GatewayIp    string   `json:"gateway_ip"`
	DnsList      []string `json:"dns_list"`
	NtpAddresses []string `json:"ntp_addresses"`
}

// TCloudSubnet defines tencent cloud subnet.
type TCloudSubnet Subnet[TCloudSubnetExtension]

// GetCloudID ...
func (vpc TCloudSubnet) GetCloudID() string {
	return vpc.CloudID
}

// AwsSubnet defines aws subnet.
type AwsSubnet Subnet[AwsSubnetExtension]

// GetCloudID ...
func (vpc AwsSubnet) GetCloudID() string {
	return vpc.CloudID
}

// GcpSubnet defines gcp subnet.
type GcpSubnet Subnet[GcpSubnetExtension]

// GetCloudID ...
func (vpc GcpSubnet) GetCloudID() string {
	return vpc.CloudID
}

// AzureSubnet defines azure subnet.
type AzureSubnet Subnet[AzureSubnetExtension]

// GetCloudID ...
func (vpc AzureSubnet) GetCloudID() string {
	return vpc.CloudID
}

// HuaWeiSubnet defines huawei subnet.
type HuaWeiSubnet Subnet[HuaWeiSubnetExtension]

// GetCloudID ...
func (vpc HuaWeiSubnet) GetCloudID() string {
	return vpc.CloudID
}
