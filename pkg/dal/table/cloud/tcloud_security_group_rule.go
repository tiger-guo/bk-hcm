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

package cloud

import (
	"errors"
	"time"

	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/validator"
	"hcm/pkg/dal/table"
	"hcm/pkg/dal/table/utils"
)

// TCloudSGRuleColumns defines all the tcloud security group rule table's columns.
var TCloudSGRuleColumns = utils.MergeColumns(nil, TCloudSGRuleColumnDescriptor)

// TCloudSGRuleColumnDescriptor is TCloud Security Group Rule's column descriptors.
var TCloudSGRuleColumnDescriptor = utils.ColumnDescriptors{
	{Column: "id", NamedC: "id", Type: enumor.String},
	{Column: "rule_index", NamedC: "rule_index", Type: enumor.Numeric},
	{Column: "version", NamedC: "version", Type: enumor.String},
	{Column: "protocol", NamedC: "protocol", Type: enumor.String},
	{Column: "cloud_service_id", NamedC: "cloud_service_id", Type: enumor.String},
	{Column: "cloud_service_group_id", NamedC: "cloud_service_group_id", Type: enumor.String},
	{Column: "ipv4_cidr", NamedC: "ipv4_cidr", Type: enumor.String},
	{Column: "ipv6_cidr", NamedC: "ipv6_cidr", Type: enumor.String},
	{Column: "cloud_target_security_group_id", NamedC: "cloud_target_security_group_id", Type: enumor.String},
	{Column: "cloud_address_id", NamedC: "cloud_address_id", Type: enumor.String},
	{Column: "cloud_address_group_id", NamedC: "cloud_address_group_id", Type: enumor.String},
	{Column: "action", NamedC: "action", Type: enumor.String},
	{Column: "memo", NamedC: "memo", Type: enumor.String},
	{Column: "type", NamedC: "type", Type: enumor.String},
	{Column: "region", NamedC: "region", Type: enumor.String},
	{Column: "cloud_security_group_id", NamedC: "cloud_security_group_id", Type: enumor.String},
	{Column: "security_group_id", NamedC: "security_group_id", Type: enumor.String},
	{Column: "account_id", NamedC: "account_id", Type: enumor.String},
	{Column: "creator", NamedC: "creator", Type: enumor.String},
	{Column: "reviser", NamedC: "reviser", Type: enumor.String},
	{Column: "created_at", NamedC: "created_at", Type: enumor.Time},
	{Column: "updated_at", NamedC: "updated_at", Type: enumor.Time},
}

// TCloudSecurityGroupRuleTable define tcloud security group rule table.
type TCloudSecurityGroupRuleTable struct {
	ID                         string     `db:"id" validate:"lte=64"`
	PolicyIndex                int64      `db:"policy_index"`
	Version                    string     `db:"version"`
	Type                       string     `db:"type" validate:"lte=20"`
	CloudSecurityGroupID       string     `db:"cloud_security_group_id" validate:"lte=255"`
	SecurityGroupID            string     `db:"security_group_id" validate:"lte=64"`
	AccountID                  string     `db:"account_id" validate:"lte=64"`
	Action                     string     `db:"action" validate:"lte=10"`
	Protocol                   *string    `db:"protocol" validate:"lte=10"`
	Port                       *string    `db:"port" validate:"lte=255"`
	CloudServiceID             *string    `db:"cloud_service_id" validate:"lte=255"`
	CloudServiceGroupID        *string    `db:"cloud_service_group_id" validate:"lte=255"`
	IPv4Cidr                   *string    `db:"ipv4_cidr" validate:"lte=255"`
	IPv6Cidr                   *string    `db:"ipv6_cidr" validate:"lte=255"`
	CloudTargetSecurityGroupID *string    `db:"cloud_target_security_group_id" validate:"lte=255"`
	CloudAddressID             *string    `db:"cloud_address_id" validate:"lte=255"`
	CloudAddressGroupID        *string    `db:"cloud_address_group_id" validate:"lte=255"`
	Region                     string     `db:"region" validate:"lte=20"`
	Memo                       *string    `db:"memo" validate:"lte=64"`
	Creator                    string     `db:"creator" validate:"lte=64"`
	Reviser                    string     `db:"reviser" validate:"lte=64"`
	CreatedAt                  *time.Time `db:"created_at" validate:"excluded_unless"`
	UpdatedAt                  *time.Time `db:"updated_at" validate:"excluded_unless"`
}

// TableName return aws security group rule table name.
func (t TCloudSecurityGroupRuleTable) TableName() table.Name {
	return table.TCloudSecurityGroupRuleTable
}

// InsertValidate aws security group rule table when insert.
func (t TCloudSecurityGroupRuleTable) InsertValidate() error {
	// length validate.
	if err := validator.Validate.Struct(t); err != nil {
		return err
	}

	if len(t.ID) == 0 {
		return errors.New("id is required")
	}

	if len(t.Region) == 0 {
		return errors.New("region is required")
	}

	if t.PolicyIndex == 0 {
		return errors.New("rule index is required")
	}

	if len(t.Version) == 0 {
		return errors.New("version is required")
	}

	if len(t.Action) == 0 {
		return errors.New("action is required")
	}

	if len(t.CloudSecurityGroupID) == 0 {
		return errors.New("cloud security group id is required")
	}

	if len(t.SecurityGroupID) == 0 {
		return errors.New("security group id is required")
	}

	if len(t.AccountID) == 0 {
		return errors.New("account id is required")
	}

	if len(t.Type) == 0 {
		return errors.New("type is required")
	}

	if len(t.Creator) == 0 {
		return errors.New("creator is required")
	}

	if len(t.Reviser) == 0 {
		return errors.New("reviser is required")
	}

	return nil
}

// UpdateValidate aws security group rule table when update.
func (t TCloudSecurityGroupRuleTable) UpdateValidate() error {
	// length validate.
	if err := validator.Validate.Struct(t); err != nil {
		return err
	}

	if len(t.Creator) != 0 {
		return errors.New("creator can not update")
	}

	return nil
}
