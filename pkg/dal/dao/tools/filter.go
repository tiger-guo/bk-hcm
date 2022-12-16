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

package tools

import (
	"hcm/pkg/runtime/filter"
)

// EqualExpression 生成资源字段等于查询的过滤条件，即fieldName=value
func EqualExpression(fieldName string, value interface{}) *filter.Expression {
	return &filter.Expression{
		Op: filter.And,
		Rules: []filter.RuleFactory{
			filter.AtomRule{Field: fieldName, Op: filter.Equal.Factory(), Value: value},
		},
	}
}

// ContainersExpression 生成资源字段包含的过滤条件，即fieldName in (1,2,3)
func ContainersExpression(fieldName string, values interface{}) *filter.Expression {
	return &filter.Expression{
		Op: filter.And,
		Rules: []filter.RuleFactory{
			filter.AtomRule{Field: fieldName, Op: filter.In.Factory(), Value: values},
		},
	}
}

// DefaultSqlWhereOption define sql where option.
var DefaultSqlWhereOption = &filter.SQLWhereOption{
	Priority: filter.Priority{"id"},
}
