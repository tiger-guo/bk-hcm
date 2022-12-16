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

package validator

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	// qualifiedNameFmt hcm resource's name format.
	// '.' And '/' as reserved characters, users are absolutely not allowed to create
	qualifiedNameFmt        = "(" + lowEnglish + qnameExtNameFmt + "*)?" + lowEnglish
	qnameExtNameFmt  string = "[\u4E00-\u9FA5A-Za-z0-9-_]"
)

// qualifiedNameRegexp hcm resource's name regexp.
var qualifiedNameRegexp = regexp.MustCompile("^" + qualifiedNameFmt + "$")

// ValidateSecurityGroupName validate security group name's length and format.
func ValidateSecurityGroupName(name string) error {
	if len(name) < 1 {
		return errors.New("invalid name, length should >= 1")
	}

	if len(name) > 128 {
		return errors.New("invalid name, length should <= 60")
	}

	if !qualifiedNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid name: %s, only allows to include chinese、english、numbers、underscore (_)"+
				"、hyphen (-), and must start and end with an chinese、english、numbers", name)
	}

	return nil
}
