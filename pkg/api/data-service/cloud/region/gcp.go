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

package region

import (
	"errors"
	"fmt"

	"hcm/pkg/api/core"
	"hcm/pkg/api/core/cloud"
	"hcm/pkg/criteria/constant"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/validator"
	"hcm/pkg/rest"
	"hcm/pkg/runtime/filter"
)

// -------------------------- Create --------------------------

// GcpRegionCreateReq define gcp region create request.
type GcpRegionCreateReq struct {
	Regions []GcpRegionBatchCreate `json:"regions" validate:"required"`
}

// GcpRegionBatchCreate define gcp region rule when create.
type GcpRegionBatchCreate struct {
	Vendor     enumor.Vendor `json:"vendor" validate:"required"`
	RegionID   string        `json:"region_id" validate:"required"`
	RegionName string        `json:"region_name" validate:"required"`
	Status     string        `json:"status"`
	SelfLink   string        `json:"self_link"`
}

// Validate gcp region create request.
func (req *GcpRegionCreateReq) Validate() error {
	if len(req.Regions) == 0 {
		return errors.New("regions is required")
	}

	if len(req.Regions) > constant.BatchOperationMaxLimit {
		return fmt.Errorf("regions count should <= %d", constant.BatchOperationMaxLimit)
	}

	return nil
}

// -------------------------- Update --------------------------

// GcpRegionBatchUpdateReq define gcp region batch update request.
type GcpRegionBatchUpdateReq struct {
	Regions []GcpRegionBatchUpdate `json:"regions" validate:"required"`
}

// GcpRegionBatchUpdate gcp region batch update option.
type GcpRegionBatchUpdate struct {
	ID         string        `json:"id" validate:"required"`
	Vendor     enumor.Vendor `json:"vendor" validate:"required"`
	RegionID   string        `json:"region_id"`
	RegionName string        `json:"region_name"`
	Status     string        `json:"status"`
	SelfLink   string        `json:"self_link"`
}

// Validate gcp region batch update request.
func (req *GcpRegionBatchUpdateReq) Validate() error {
	if len(req.Regions) == 0 {
		return errors.New("regions is required")
	}

	if len(req.Regions) > constant.BatchOperationMaxLimit {
		return fmt.Errorf("regions count should <= %d", constant.BatchOperationMaxLimit)
	}

	return nil
}

// -------------------------- List --------------------------

// GcpRegionListReq gcp region list req.
type GcpRegionListReq struct {
	Field  []string           `json:"field" validate:"omitempty"`
	Filter *filter.Expression `json:"filter" validate:"required"`
	Page   *core.BasePage     `json:"page" validate:"required"`
}

// Validate gcp region list request.
func (req *GcpRegionListReq) Validate() error {
	return validator.Validate.Struct(req)
}

// GcpRegionListResult define gcp region list result.
type GcpRegionListResult struct {
	Count   uint64            `json:"count"`
	Details []cloud.GcpRegion `json:"details"`
}

// GcpRegionListResp define gcp region list resp.
type GcpRegionListResp struct {
	rest.BaseResp `json:",inline"`
	Data          *GcpRegionListResult `json:"data"`
}
