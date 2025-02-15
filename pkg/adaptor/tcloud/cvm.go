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

package tcloud

import (
	"fmt"

	"hcm/pkg/adaptor/poller"
	"hcm/pkg/adaptor/types"
	"hcm/pkg/adaptor/types/core"
	typecvm "hcm/pkg/adaptor/types/cvm"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/tools/converter"
	"hcm/pkg/tools/slice"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// TODO: sync-todo  改好后统一删除ListDisk函数
// ListCvmNew list cvm.
// reference: https://cloud.tencent.com/document/api/213/15728
func (t *TCloud) ListCvmNew(kt *kit.Kit, opt *typecvm.TCloudListOption) ([]typecvm.TCloudCvm, error) {

	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "list option is required")
	}

	if err := opt.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return nil, fmt.Errorf("new tcloud vpc client failed, err: %v", err)
	}

	req := cvm.NewDescribeInstancesRequest()
	if len(opt.CloudIDs) != 0 {
		req.InstanceIds = common.StringPtrs(opt.CloudIDs)
		req.Limit = common.Int64Ptr(int64(core.TCloudQueryLimit))
	}

	if opt.Page != nil {
		req.Offset = common.Int64Ptr(int64(opt.Page.Offset))
		req.Limit = common.Int64Ptr(int64(opt.Page.Limit))
	}

	resp, err := client.DescribeInstancesWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("list tcloud instance failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	cvms := make([]typecvm.TCloudCvm, 0, len(resp.Response.InstanceSet))
	for _, one := range resp.Response.InstanceSet {
		cvms = append(cvms, typecvm.TCloudCvm{one})
	}

	return cvms, nil
}

// ListCvm list cvm.
// reference: https://cloud.tencent.com/document/api/213/15728
func (t *TCloud) ListCvm(kt *kit.Kit, opt *typecvm.TCloudListOption) ([]*cvm.Instance, error) {

	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "list option is required")
	}

	if err := opt.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return nil, fmt.Errorf("new tcloud vpc client failed, err: %v", err)
	}

	req := cvm.NewDescribeInstancesRequest()
	if len(opt.CloudIDs) != 0 {
		req.InstanceIds = common.StringPtrs(opt.CloudIDs)
		req.Limit = common.Int64Ptr(int64(core.TCloudQueryLimit))
	}

	if opt.Page != nil {
		req.Offset = common.Int64Ptr(int64(opt.Page.Offset))
		req.Limit = common.Int64Ptr(int64(opt.Page.Limit))
	}

	resp, err := client.DescribeInstancesWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("list tcloud instance failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	return resp.Response.InstanceSet, nil
}

// DeleteCvm reference: https://cloud.tencent.com/document/api/213/15723
func (t *TCloud) DeleteCvm(kt *kit.Kit, opt *typecvm.TCloudDeleteOption) error {

	if opt == nil {
		return errf.New(errf.InvalidParameter, "start cvm option is required")
	}

	if err := opt.Validate(); err != nil {
		return errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return fmt.Errorf("init tencent cloud client failed, err: %v", err)
	}

	req := cvm.NewTerminateInstancesRequest()
	req.InstanceIds = common.StringPtrs(opt.CloudIDs)

	_, err = client.TerminateInstancesWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("terminate cvm instance failed, err: %v, rid: %s", err, kt.Rid)
		return err
	}

	return nil
}

// StartCvm reference: https://cloud.tencent.com/document/api/213/15735
func (t *TCloud) StartCvm(kt *kit.Kit, opt *typecvm.TCloudStartOption) error {

	if opt == nil {
		return errf.New(errf.InvalidParameter, "start cvm option is required")
	}

	if err := opt.Validate(); err != nil {
		return errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return fmt.Errorf("init tencent cloud client failed, err: %v", err)
	}
	req := cvm.NewStartInstancesRequest()
	req.InstanceIds = common.StringPtrs(opt.CloudIDs)

	_, err = client.StartInstancesWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("start cvm failed, err: %v, ids: %v, rid: %s", err, opt.CloudIDs, kt.Rid)
		return err
	}

	// wait until all cvm done
	handler := &startCvmPollingHandler{
		opt.Region,
	}
	respPoller := poller.Poller[*TCloud, []*cvm.Instance, poller.BaseDoneResult]{Handler: handler}
	res, err := respPoller.PollUntilDone(t, kt, converter.SliceToPtr(opt.CloudIDs), types.NewBatchOperateCvmPollerOpt())
	if err != nil {
		logs.Errorf("poll start cvm failed, err: %v, res: %#v, rid: %s", err, res, kt.Rid)
		return err
	}

	return nil
}

// StopCvm reference: https://cloud.tencent.com/document/api/213/15743
func (t *TCloud) StopCvm(kt *kit.Kit, opt *typecvm.TCloudStopOption) error {

	if opt == nil {
		return errf.New(errf.InvalidParameter, "stop cvm option is required")
	}

	if err := opt.Validate(); err != nil {
		return errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return fmt.Errorf("init tencent cloud client failed, err: %v", err)
	}

	req := cvm.NewStopInstancesRequest()
	req.InstanceIds = common.StringPtrs(opt.CloudIDs)
	req.StopType = common.StringPtr(string(opt.StopType))
	req.StoppedMode = common.StringPtr(string(opt.StoppedMode))

	_, err = client.StopInstancesWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("stop cvm failed, err: %v, ids: %v, rid: %s", err, opt.CloudIDs, kt.Rid)
		return err
	}

	// wait until all cvm done
	handler := &stopCvmPollingHandler{
		opt.Region,
	}
	respPoller := poller.Poller[*TCloud, []*cvm.Instance, poller.BaseDoneResult]{Handler: handler}
	res, err := respPoller.PollUntilDone(t, kt, converter.SliceToPtr(opt.CloudIDs), types.NewBatchOperateCvmPollerOpt())
	if err != nil {
		logs.Errorf("poll stop cvm failed, err: %v, res: %#v, rid: %s", err, res, kt.Rid)
		return err
	}

	return nil
}

// RebootCvm reference: https://cloud.tencent.com/document/api/213/15742
func (t *TCloud) RebootCvm(kt *kit.Kit, opt *typecvm.TCloudRebootOption) error {

	if opt == nil {
		return errf.New(errf.InvalidParameter, "reboot cvm option is required")
	}

	if err := opt.Validate(); err != nil {
		return errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return fmt.Errorf("init tencent cloud client failed, err: %v", err)
	}

	req := cvm.NewRebootInstancesRequest()
	req.InstanceIds = common.StringPtrs(opt.CloudIDs)
	req.StopType = common.StringPtr(string(opt.StopType))

	_, err = client.RebootInstancesWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("reboot cvm failed, err: %v, ids: %v, rid: %s", err, opt.CloudIDs, kt.Rid)
		return err
	}

	// wait until all cvm are rebooted
	handler := &rebootCvmPollingHandler{
		opt.Region,
	}
	respPoller := poller.Poller[*TCloud, []*cvm.Instance, poller.BaseDoneResult]{Handler: handler}
	res, err := respPoller.PollUntilDone(t, kt, converter.SliceToPtr(opt.CloudIDs), types.NewBatchOperateCvmPollerOpt())
	if err != nil {
		logs.Errorf("poll reboot cvm failed, err: %v, res: %#v, rid: %s", err, res, kt.Rid)
		return err
	}

	return nil
}

// ResetCvmPwd reference: https://cloud.tencent.com/document/api/213/15736
func (t *TCloud) ResetCvmPwd(kt *kit.Kit, opt *typecvm.TCloudResetPwdOption) error {

	if opt == nil {
		return errf.New(errf.InvalidParameter, "reset pwd option is required")
	}

	if err := opt.Validate(); err != nil {
		return errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return fmt.Errorf("init tencent cloud client failed, err: %v", err)
	}

	req := cvm.NewResetInstancesPasswordRequest()
	req.InstanceIds = common.StringPtrs(opt.CloudIDs)
	req.Password = common.StringPtr(opt.Password)
	req.UserName = common.StringPtr(opt.UserName)
	req.ForceStop = common.BoolPtr(opt.ForceStop)

	_, err = client.ResetInstancesPasswordWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("reset cvm instance's password failed, err: %v, rid: %s", err, kt.Rid)
		return err
	}

	// wait until all cvm done
	handler := &resetpwdCvmPollingHandler{
		opt.Region,
	}
	respPoller := poller.Poller[*TCloud, []*cvm.Instance, poller.BaseDoneResult]{Handler: handler}
	res, err := respPoller.PollUntilDone(t, kt, converter.SliceToPtr(opt.CloudIDs), types.NewBatchOperateCvmPollerOpt())
	if err != nil {
		logs.Errorf("poll reset pwd cvm failed, err: %v, res: %#v, rid: %s", err, res, kt.Rid)
		return err
	}

	return nil
}

// CreateCvm reference: https://cloud.tencent.com/document/api/213/15730
// NOTE：返回实例`ID`列表并不代表实例创建成功，可根据 [DescribeInstances](https://cloud.tencent.com/document/api/213/15728)
// 接口查询返回的InstancesSet中对应实例的`ID`的状态来判断创建是否完成；如果实例状态由“PENDING(创建中)”变为“RUNNING(运行中)”，则为创建成功。
func (t *TCloud) CreateCvm(kt *kit.Kit, opt *typecvm.TCloudCreateOption) (*poller.BaseDoneResult, error) {
	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "create option is required")
	}

	if err := opt.Validate(); err != nil {
		return nil, errf.NewFromErr(errf.InvalidParameter, err)
	}

	client, err := t.clientSet.cvmClient(opt.Region)
	if err != nil {
		return nil, fmt.Errorf("init tencent cloud client failed, err: %v", err)
	}

	req := cvm.NewRunInstancesRequest()
	req.Placement = &cvm.Placement{
		Zone: common.StringPtr(opt.Zone),
	}
	req.DryRun = common.BoolPtr(opt.DryRun)
	req.InstanceType = common.StringPtr(opt.InstanceType)
	req.ImageId = common.StringPtr(opt.CloudImageID)
	req.InstanceCount = common.Int64Ptr(opt.RequiredCount)
	req.InstanceName = common.StringPtr(opt.Name)
	req.SecurityGroupIds = common.StringPtrs(opt.CloudSecurityGroupIDs)
	req.ClientToken = opt.ClientToken
	req.InstanceChargeType = common.StringPtr(string(opt.InstanceChargeType))
	req.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:    common.StringPtr(opt.CloudVpcID),
		SubnetId: common.StringPtr(opt.CloudSubnetID),
	}
	req.LoginSettings = &cvm.LoginSettings{
		Password: common.StringPtr(opt.Password),
	}
	req.InternetAccessible = &cvm.InternetAccessible{
		PublicIpAssigned: common.BoolPtr(opt.PublicIPAssigned),
	}

	req.SystemDisk = &cvm.SystemDisk{
		DiskId:   opt.SystemDisk.CloudDiskID,
		DiskSize: opt.SystemDisk.DiskSizeGB,
	}
	if len(opt.SystemDisk.DiskType) != 0 {
		req.SystemDisk.DiskType = common.StringPtr(string(opt.SystemDisk.DiskType))
	}

	if len(opt.DataDisk) != 0 {
		req.DataDisks = make([]*cvm.DataDisk, 0, len(opt.DataDisk))
		for _, one := range opt.DataDisk {
			disk := &cvm.DataDisk{
				DiskSize: one.DiskSizeGB,
				DiskId:   one.CloudDiskID,
			}

			if len(one.DiskType) != 0 {
				disk.DiskType = common.StringPtr(string(one.DiskType))
			}
			req.DataDisks = append(req.DataDisks, disk)
		}
	}

	if opt.InstanceChargePrepaid != nil {
		req.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
			Period: opt.InstanceChargePrepaid.Period,
		}

		if len(opt.InstanceChargePrepaid.RenewFlag) != 0 {
			req.InstanceChargePrepaid.RenewFlag = common.StringPtr(string(opt.InstanceChargePrepaid.RenewFlag))
		}
	}

	resp, err := client.RunInstancesWithContext(kt.Ctx, req)
	if err != nil {
		logs.Errorf("run tencent cloud cvm instance failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	// 参数预校验不生产资源
	if opt.DryRun {
		return new(poller.BaseDoneResult), nil
	}

	handler := &createCvmPollingHandler{
		opt.Region,
	}
	respPoller := poller.Poller[*TCloud, []*cvm.Instance, poller.BaseDoneResult]{Handler: handler}
	result, err := respPoller.PollUntilDone(t, kt, resp.Response.InstanceIdSet, types.NewBatchCreateCvmPollerOption())
	if err != nil {
		return nil, err
	}

	return result, nil
}

type startCvmPollingHandler struct {
	region string
}

// Done ...
func (h *startCvmPollingHandler) Done(cvms []*cvm.Instance) (bool, *poller.BaseDoneResult) {
	return done(cvms, "RUNNING")
}

// Poll ...
func (h *startCvmPollingHandler) Poll(client *TCloud, kt *kit.Kit, cloudIDs []*string) ([]*cvm.Instance, error) {
	return poll(client, kt, h.region, cloudIDs)
}

type stopCvmPollingHandler struct {
	region string
}

// Done ...
func (h *stopCvmPollingHandler) Done(cvms []*cvm.Instance) (bool, *poller.BaseDoneResult) {
	return done(cvms, "STOPPED")
}

// Poll ...
func (h *stopCvmPollingHandler) Poll(client *TCloud, kt *kit.Kit, cloudIDs []*string) ([]*cvm.Instance, error) {
	return poll(client, kt, h.region, cloudIDs)
}

type resetpwdCvmPollingHandler struct {
	region string
}

// Done ...
func (h *resetpwdCvmPollingHandler) Done(cvms []*cvm.Instance) (bool, *poller.BaseDoneResult) {
	return done(cvms, "RUNNING")
}

// Poll ...
func (h *resetpwdCvmPollingHandler) Poll(client *TCloud, kt *kit.Kit, cloudIDs []*string) ([]*cvm.Instance, error) {
	return poll(client, kt, h.region, cloudIDs)
}

type rebootCvmPollingHandler struct {
	region string
}

// Done ...
func (h *rebootCvmPollingHandler) Done(cvms []*cvm.Instance) (bool, *poller.BaseDoneResult) {
	return done(cvms, "RUNNING")
}

// Poll ...
func (h *rebootCvmPollingHandler) Poll(client *TCloud, kt *kit.Kit, cloudIDs []*string) ([]*cvm.Instance, error) {
	return poll(client, kt, h.region, cloudIDs)
}

func done(cvms []*cvm.Instance, succeed string) (bool, *poller.BaseDoneResult) {
	result := new(poller.BaseDoneResult)

	flag := true
	for _, instance := range cvms {
		// not done
		if converter.PtrToVal(instance.InstanceState) != succeed {
			flag = false
			continue
		}

		result.SuccessCloudIDs = append(result.SuccessCloudIDs, *instance.InstanceId)
	}

	return flag, result
}

func poll(client *TCloud, kt *kit.Kit, region string, cloudIDs []*string) ([]*cvm.Instance, error) {
	cloudIDSplit := slice.Split(cloudIDs, core.TCloudQueryLimit)

	cvms := make([]*cvm.Instance, 0, len(cloudIDs))
	for _, partIDs := range cloudIDSplit {
		req := cvm.NewDescribeInstancesRequest()
		req.InstanceIds = partIDs
		req.Limit = converter.ValToPtr(int64(core.TCloudQueryLimit))

		cvmCli, err := client.clientSet.cvmClient(region)
		if err != nil {
			return nil, err
		}

		resp, err := cvmCli.DescribeInstancesWithContext(kt.Ctx, req)
		if err != nil {
			return nil, err
		}

		cvms = append(cvms, resp.Response.InstanceSet...)
	}

	return cvms, nil
}

type createCvmPollingHandler struct {
	region string
}

func (h *createCvmPollingHandler) Done(cvms []*cvm.Instance) (bool, *poller.BaseDoneResult) {

	result := &poller.BaseDoneResult{
		SuccessCloudIDs: make([]string, 0),
		FailedCloudIDs:  make([]string, 0),
		UnknownCloudIDs: make([]string, 0),
	}
	flag := true
	for _, instance := range cvms {
		// 创建中
		if converter.PtrToVal(instance.InstanceState) == "PENDING" {
			flag = false
			result.UnknownCloudIDs = append(result.UnknownCloudIDs, *instance.InstanceId)
			continue
		}

		// 生产失败
		if converter.PtrToVal(instance.InstanceState) == "LAUNCH_FAILED" {
			result.FailedCloudIDs = append(result.FailedCloudIDs, *instance.InstanceId)
			result.FailedMessage = converter.PtrToVal(instance.LatestOperationErrorMsg)
			continue
		}

		result.SuccessCloudIDs = append(result.SuccessCloudIDs, *instance.InstanceId)
	}

	return flag, result
}

func (h *createCvmPollingHandler) Poll(client *TCloud, kt *kit.Kit, cloudIDs []*string) ([]*cvm.Instance, error) {

	cloudIDSplit := slice.Split(cloudIDs, core.TCloudQueryLimit)

	cvms := make([]*cvm.Instance, 0, len(cloudIDs))
	for _, partIDs := range cloudIDSplit {
		req := cvm.NewDescribeInstancesRequest()
		req.InstanceIds = partIDs
		req.Limit = converter.ValToPtr(int64(core.TCloudQueryLimit))

		cvmCli, err := client.clientSet.cvmClient(h.region)
		if err != nil {
			return nil, err
		}

		resp, err := cvmCli.DescribeInstancesWithContext(kt.Ctx, req)
		if err != nil {
			return nil, err
		}

		cvms = append(cvms, resp.Response.InstanceSet...)
	}

	if len(cvms) != len(cloudIDs) {
		return nil, fmt.Errorf("query cvm count: %d not equal return count: %d", len(cloudIDs), len(cvms))
	}

	return cvms, nil
}

var _ poller.PollingHandler[*TCloud, []*cvm.Instance, poller.BaseDoneResult] = new(createCvmPollingHandler)
