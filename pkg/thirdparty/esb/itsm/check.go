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

package itsm

import (
	"context"
	"fmt"

	"hcm/pkg/thirdparty/esb/types"
)

type tokenVerifiedResp struct {
	types.BaseResponse `json:",inline"`
	Data               struct {
		IsPassed bool `json:"is_passed"`
	} `json:"data"`
}

func (i *itsm) VerifyToken(ctx context.Context, token string) (bool, error) {
	req := map[string]string{"token": token}
	resp := new(tokenVerifiedResp)
	header := types.GetCommonHeader(i.config)
	err := i.client.Post().
		SubResourcef("/itsm/token/verify/").
		WithContext(ctx).
		WithHeaders(*header).
		Body(req).
		Do().Into(resp)
	if err != nil {
		return false, err
	}
	if !resp.Result || resp.Code != 0 {
		return false, fmt.Errorf("verify token failed, code: %d, msg: %s, rid: %s", resp.Code, resp.Message, resp.Rid)
	}

	return resp.Data.IsPassed, nil
}
