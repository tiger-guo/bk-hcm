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

package huawei

import (
	ressync "hcm/cmd/hc-service/logics/res-sync"
	"hcm/cmd/hc-service/logics/res-sync/huawei"
	"hcm/cmd/hc-service/service/sync/handler"
	"hcm/pkg/adaptor/types/core"
	typecvm "hcm/pkg/adaptor/types/cvm"
	"hcm/pkg/api/hc-service/sync"
	"hcm/pkg/criteria/constant"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/rest"
)

// SyncCvmWithRelRes ....
func (svc *service) SyncCvmWithRelRes(cts *rest.Contexts) (interface{}, error) {
	return nil, handler.ResourceSync(cts, &cvmHandler{cli: svc.syncCli})
}

// cvmHandler cvm sync handler.
type cvmHandler struct {
	cli ressync.Interface

	// Prepare 构建参数
	request *sync.HuaWeiSyncReq
	syncCli huawei.Interface
	// offset cvm分页使用
	offset int32
}

var _ handler.Handler = new(cvmHandler)

// Prepare ...
func (hd *cvmHandler) Prepare(cts *rest.Contexts) error {
	request, syncCli, err := defaultPrepare(cts, hd.cli)
	if err != nil {
		return err
	}

	hd.request = request
	hd.syncCli = syncCli

	return nil
}

// Next ...
func (hd *cvmHandler) Next(kt *kit.Kit) ([]string, error) {
	listOpt := &typecvm.HuaWeiListOption{
		Region: hd.request.Region,
		Page: &core.HuaWeiCvmOffsetPage{
			Offset: hd.offset,
			Limit:  int32(constant.CloudResourceSyncMaxLimit),
		},
	}
	cvmResult, err := hd.syncCli.CloudCli().ListCvmNew(kt, listOpt)
	if err != nil {
		logs.Errorf("request adaptor list huawei cvm failed, err: %v, opt: %v, rid: %s", err, listOpt, kt.Rid)
		return nil, err
	}

	if len(cvmResult) == 0 {
		return nil, nil
	}

	cloudIDs := make([]string, 0, len(cvmResult))
	for _, one := range cvmResult {
		cloudIDs = append(cloudIDs, one.Id)
	}

	hd.offset += 1
	return cloudIDs, nil
}

// Sync ...
func (hd *cvmHandler) Sync(kt *kit.Kit, cloudIDs []string) error {
	params := &huawei.SyncBaseParams{
		AccountID: hd.request.AccountID,
		Region:    hd.request.Region,
		CloudIDs:  cloudIDs,
	}
	if _, err := hd.syncCli.CvmWithRelRes(kt, params, new(huawei.SyncCvmWithRelResOption)); err != nil {
		logs.Errorf("sync huawei cvm failed, err: %v, opt: %v, rid: %s", err, params, kt.Rid)
		return err
	}

	return nil
}

// RemoveDeleteFromCloud ...
func (hd *cvmHandler) RemoveDeleteFromCloud(kt *kit.Kit) error {
	if err := hd.syncCli.RemoveCvmDeleteFromCloud(kt, hd.request.AccountID, hd.request.Region); err != nil {
		logs.Errorf("remove cvm delete from cloud failed, err: %v, accountID: %s, region: %s, rid: %s", err,
			hd.request.AccountID, hd.request.Region, kt.Rid)
		return err
	}

	return nil
}

// Name ...
func (hd *cvmHandler) Name() enumor.CloudResourceType {
	return enumor.CvmCloudResType
}
