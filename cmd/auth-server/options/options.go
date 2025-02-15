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

package options

import (
	"sync"

	"hcm/pkg/cc"
	"hcm/pkg/runtime/flags"

	"github.com/spf13/pflag"
)

// Option defines the app's runtime flag options.
type Option struct {
	Sys *cc.SysOption
	// DisableAuth defines whether iam authorization is disabled
	DisableAuth bool
}

// InitOptions init auth server's options from command flags.
func InitOptions() *Option {
	fs := pflag.CommandLine
	sysOpt := flags.SysFlags(fs)

	disableAuth := true
	fs.BoolVar(&disableAuth, "disable-auth", false, "defines whether iam authorization is disabled. Auth is enabled "+
		"by default, disable-auth=true needs to set explicitly to disable iam authorization. Note: disable auth "+
		"may cause security problems, please do not use it unless you understand and accept the risks.")

	// parses the command-line flags from os.Args[1:]. must be called after all flags are defined
	// and before flags are accessed by the program.
	pflag.Parse()

	// check if the command-line flag is show current version info cmd.
	sysOpt.CheckV()

	return &Option{Sys: sysOpt, DisableAuth: disableAuth}
}

// DisableWriteOption defines which biz's write operations needs to be disabled when authorized
type DisableWriteOption struct {
	// IsDisabled defines if the write operations needs to be disabled
	IsDisabled bool
	// IsAll defines if all write operations needs to be disabled
	IsAll bool
	// BizIDMap defines the biz ids that needs to be disabled as map keys
	BizIDMap sync.Map
}
