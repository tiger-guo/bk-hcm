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

// Package account Package service defines service service.
package account

import (
	"net/http"

	"hcm/cmd/hc-service/service/capability"
	cloudadaptor "hcm/cmd/hc-service/service/cloud-adaptor"
	"hcm/pkg/rest"
)

// InitAccountService initial the service service
func InitAccountService(cap *capability.Capability) {
	svc := &service{
		ad: cap.CloudAdaptor,
	}

	h := rest.NewHandler()
	// 账号检查
	h.Add("TCloudAccountCheck", http.MethodPost, "/vendors/tcloud/accounts/check", svc.TCloudAccountCheck)
	h.Add("AwsAccountCheck", http.MethodPost, "/vendors/aws/accounts/check", svc.AwsAccountCheck)
	h.Add("HuaWeiAccountCheck", http.MethodPost, "/vendors/huawei/accounts/check", svc.HuaWeiAccountCheck)
	h.Add("GcpAccountCheck", http.MethodPost, "/vendors/gcp/accounts/check", svc.GcpAccountCheck)
	h.Add("AzureAccountCheck", http.MethodPost, "/vendors/azure/accounts/check", svc.AzureAccountCheck)

	// 获取账号配额
	h.Add("GetTCloudAccountZoneQuota", http.MethodPost, "/vendors/tcloud/accounts/zones/quotas",
		svc.GetTCloudAccountZoneQuota)
	h.Add("GetHuaWeiAccountRegionQuota", http.MethodPost, "/vendors/huawei/accounts/regions/quotas",
		svc.GetHuaWeiAccountRegionQuota)
	h.Add("GetGcpAccountRegionQuota", http.MethodPost, "/vendors/gcp/accounts/regions/quotas",
		svc.GetGcpAccountRegionQuota)

	h.Load(cap.WebService)
}

type service struct {
	ad *cloudadaptor.CloudAdaptorClient
}
