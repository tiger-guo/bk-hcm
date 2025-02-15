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

package constant

// Note:
// This scope is used to define all the constant keys which is used inside and outside
// the HCM system.
const (
	// BatchOperationMaxLimit 批量操作最大上限，包括批量创建、批量更新、批量删除。
	BatchOperationMaxLimit = 100
	// RelResourceOperationMaxLimit 批量操作最大上限，包括批量创建、批量更新、批量删除。
	RelResourceOperationMaxLimit = 500
	// SGBatchOperationMaxLimit 批量操作最大上限，包括批量创建、批量更新、批量删除。
	SGBatchOperationMaxLimit = 1000

	// CloudResourceSyncMaxLimit 单次云资源同步最大数量限制。
	CloudResourceSyncMaxLimit = 100
	// SyncConcurrencyDefaultMaxLimit 同步并发最大限制
	SyncConcurrencyDefaultMaxLimit = 10
)
