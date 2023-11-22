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

package cloudselection

import (
	"net/http"

	"hcm/cmd/data-service/service/capability"
	"hcm/pkg/dal/dao"
	"hcm/pkg/rest"
)

// InitService initial the service
func InitService(cap *capability.Capability) {
	svc := &service{
		dao: cap.Dao,
	}

	h := rest.NewHandler()

	h.Add("ListScheme", http.MethodPost, "/clouds/selections/schemes/list",
		svc.ListScheme)
	h.Add("BatchDeleteScheme", http.MethodDelete, "/clouds/selections/schemes/batch",
		svc.BatchDeleteScheme)
	h.Add("CreateScheme", http.MethodPost, "/clouds/selections/schemes/create",
		svc.CreateScheme)
	h.Add("UpdateScheme", http.MethodPatch, "/clouds/selections/schemes/{id}",
		svc.UpdateScheme)

	h.Add("ListIdc", http.MethodPost, "/clouds/selections/idcs/list",
		svc.ListIdc)

	h.Add("ListBizType", http.MethodPost, "/clouds/selections/biz_types/list",
		svc.ListBizType)

	h.Load(cap.WebService)
}

type service struct {
	dao dao.Set
}
