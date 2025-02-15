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

package cc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"hcm/pkg/logs"
	"hcm/pkg/tools/ssl"
	"hcm/pkg/version"

	etcd3 "go.etcd.io/etcd/client/v3"
)

// Service defines Setting related runtime.
type Service struct {
	Etcd Etcd `yaml:"etcd"`
}

// trySetDefault set the Setting default value if user not configured.
func (s *Service) trySetDefault() {
	s.Etcd.trySetDefault()
}

// validate Setting related runtime.
func (s Service) validate() error {
	if err := s.Etcd.validate(); err != nil {
		return err
	}

	return nil
}

// Etcd defines etcd related runtime
type Etcd struct {
	// Endpoints is a list of URLs.
	Endpoints []string `yaml:"endpoints"`
	// DialTimeoutMS is the timeout seconds for failing
	// to establish a connection.
	DialTimeoutMS uint `yaml:"dialTimeoutMS"`
	// Username is a user's name for authentication.
	Username string `yaml:"username"`
	// Password is a password for authentication.
	Password string    `yaml:"password"`
	TLS      TLSConfig `yaml:"tls"`
}

// trySetDefault set the etcd default value if user not configured.
func (es *Etcd) trySetDefault() {
	if len(es.Endpoints) == 0 {
		es.Endpoints = []string{"127.0.0.1:2379"}
	}

	if es.DialTimeoutMS == 0 {
		es.DialTimeoutMS = 200
	}
}

// ToConfig convert to etcd config.
func (es Etcd) ToConfig() (etcd3.Config, error) {
	var tlsC *tls.Config
	if es.TLS.Enable() {
		var err error
		tlsC, err = ssl.ClientTLSConfVerify(es.TLS.InsecureSkipVerify, es.TLS.CAFile, es.TLS.CertFile,
			es.TLS.KeyFile, es.TLS.Password)
		if err != nil {
			return etcd3.Config{}, fmt.Errorf("init etcd tls config failed, err: %v", err)
		}
	}

	c := etcd3.Config{
		Endpoints:            es.Endpoints,
		AutoSyncInterval:     0,
		DialTimeout:          time.Duration(es.DialTimeoutMS) * time.Millisecond,
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  tlsC,
		Username:             es.Username,
		Password:             es.Password,
		RejectOldCluster:     false,
		DialOptions:          nil,
		Context:              nil,
		LogConfig:            nil,
		PermitWithoutStream:  false,
	}

	return c, nil
}

// validate etcd runtime
func (es Etcd) validate() error {
	if len(es.Endpoints) == 0 {
		return errors.New("etcd endpoints is not set")
	}

	if err := es.TLS.validate(); err != nil {
		return fmt.Errorf("etcd tls, %v", err)
	}

	return nil
}

// Limiter defines the request limit options
type Limiter struct {
	// QPS should >=1
	QPS uint `yaml:"qps"`
	// Burst should >= 1;
	Burst uint `yaml:"burst"`
}

// validate if the limiter is valid or not.
func (lm Limiter) validate() error {
	if lm.QPS <= 0 {
		return errors.New("invalid QPS value, should >= 1")
	}

	if lm.Burst <= 0 {
		return errors.New("invalid Burst value, should >= 1")
	}

	return nil
}

// trySetDefault try set the default value of limiter
func (lm *Limiter) trySetDefault() {
	if lm.QPS == 0 {
		lm.QPS = 500
	}

	if lm.Burst == 0 {
		lm.Burst = 500
	}
}

// DataBase defines database related runtime
type DataBase struct {
	Resource ResourceDB `yaml:"resource"`
	// MaxSlowLogLatencyMS defines the max tolerance in millisecond to execute
	// the database command, if the cost time of execute have >= the MaxSlowLogLatencyMS
	// then this request will be logged.
	MaxSlowLogLatencyMS uint `yaml:"maxSlowLogLatencyMS"`
	// Limiter defines request's to ORM's limitation for each sharding, and
	// each sharding have the independent request limitation.
	Limiter *Limiter `yaml:"limiter"`
}

// trySetDefault set the sharding default value if user not configured.
func (s *DataBase) trySetDefault() {
	s.Resource.trySetDefault()

	if s.MaxSlowLogLatencyMS == 0 {
		s.MaxSlowLogLatencyMS = 100
	}

	if s.Limiter == nil {
		s.Limiter = new(Limiter)
	}

	s.Limiter.trySetDefault()
}

// validate sharding runtime
func (s DataBase) validate() error {
	if err := s.Resource.validate(); err != nil {
		return err
	}

	if s.MaxSlowLogLatencyMS <= 0 {
		return errors.New("invalid maxSlowLogLatencyMS")
	}

	if s.Limiter != nil {
		if err := s.Limiter.validate(); err != nil {
			return fmt.Errorf("sharding.limiter is invalid, %v", err)
		}
	}

	return nil
}

// ResourceDB defines database related runtime.
type ResourceDB struct {
	// Endpoints is a seed list of host:port addresses of database nodes.
	Endpoints []string `yaml:"endpoints"`
	Database  string   `yaml:"database"`
	User      string   `yaml:"user"`
	Password  string   `yaml:"password"`
	// DialTimeoutSec is timeout in seconds to wait for a
	// response from the db server
	// all the timeout default value reference:
	// https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html
	DialTimeoutSec    uint      `yaml:"dialTimeoutSec"`
	ReadTimeoutSec    uint      `yaml:"readTimeoutSec"`
	WriteTimeoutSec   uint      `yaml:"writeTimeoutSec"`
	MaxIdleTimeoutMin uint      `yaml:"maxIdleTimeoutMin"`
	MaxOpenConn       uint      `yaml:"maxOpenConn"`
	MaxIdleConn       uint      `yaml:"maxIdleConn"`
	TLS               TLSConfig `yaml:"tls"`
}

// trySetDefault set the database's default value if user not configured.
func (ds *ResourceDB) trySetDefault() {
	if len(ds.Endpoints) == 0 {
		ds.Endpoints = []string{"127.0.0.1:3306"}
	}

	if ds.DialTimeoutSec == 0 {
		ds.DialTimeoutSec = 15
	}

	if ds.ReadTimeoutSec == 0 {
		ds.ReadTimeoutSec = 10
	}

	if ds.WriteTimeoutSec == 0 {
		ds.WriteTimeoutSec = 10
	}

	if ds.MaxOpenConn == 0 {
		ds.MaxOpenConn = 500
	}

	if ds.MaxIdleConn == 0 {
		ds.MaxIdleConn = 5
	}
}

// validate database runtime.
func (ds ResourceDB) validate() error {
	if len(ds.Endpoints) == 0 {
		return errors.New("database endpoints is not set")
	}

	if len(ds.Database) == 0 {
		return errors.New("database is not set")
	}

	if (ds.DialTimeoutSec > 0 && ds.DialTimeoutSec < 1) || ds.DialTimeoutSec > 60 {
		return errors.New("invalid database dialTimeoutMS, should be in [1:60]s")
	}

	if (ds.ReadTimeoutSec > 0 && ds.ReadTimeoutSec < 1) || ds.ReadTimeoutSec > 60 {
		return errors.New("invalid database readTimeoutMS, should be in [1:60]s")
	}

	if (ds.WriteTimeoutSec > 0 && ds.WriteTimeoutSec < 1) || ds.WriteTimeoutSec > 30 {
		return errors.New("invalid database writeTimeoutMS, should be in [1:30]s")
	}

	if err := ds.TLS.validate(); err != nil {
		return fmt.Errorf("database tls, %v", err)
	}

	return nil
}

// LogOption defines log's related configuration
type LogOption struct {
	LogDir           string `yaml:"logDir"`
	MaxPerFileSizeMB uint32 `yaml:"maxPerFileSizeMB"`
	MaxPerLineSizeKB uint32 `yaml:"maxPerLineSizeKB"`
	MaxFileNum       uint   `yaml:"maxFileNum"`
	LogAppend        bool   `yaml:"logAppend"`
	// log the log to std err only, it can not be used with AlsoToStdErr
	// at the same time.
	ToStdErr bool `yaml:"toStdErr"`
	// log the log to file and also to std err. it can not be used with ToStdErr
	// at the same time.
	AlsoToStdErr bool `yaml:"alsoToStdErr"`
	Verbosity    uint `yaml:"verbosity"`
}

// trySetDefault set the log's default value if user not configured.
func (log *LogOption) trySetDefault() {
	if len(log.LogDir) == 0 {
		log.LogDir = "./"
	}

	if log.MaxPerFileSizeMB == 0 {
		log.MaxPerFileSizeMB = 500
	}

	if log.MaxPerLineSizeKB == 0 {
		log.MaxPerLineSizeKB = 5
	}

	if log.MaxFileNum == 0 {
		log.MaxFileNum = 5
	}
}

// Logs convert it to logs.LogConfig.
func (log LogOption) Logs() logs.LogConfig {
	l := logs.LogConfig{
		LogDir:             log.LogDir,
		LogMaxSize:         log.MaxPerFileSizeMB,
		LogLineMaxSize:     log.MaxPerLineSizeKB,
		LogMaxNum:          log.MaxFileNum,
		RestartNoScrolling: log.LogAppend,
		ToStdErr:           log.ToStdErr,
		AlsoToStdErr:       log.AlsoToStdErr,
		Verbosity:          log.Verbosity,
	}

	return l
}

// Network defines all the network related options
type Network struct {
	// BindIP is ip where server working on
	BindIP string `yaml:"bindIP"`
	// Port is port where server listen to http port.
	Port uint      `yaml:"port"`
	TLS  TLSConfig `yaml:"tls"`
}

// trySetFlagBindIP try set flag bind ip, bindIP only can set by one of the flag or configuration file.
func (n *Network) trySetFlagBindIP(ip net.IP) error {
	if len(ip) != 0 {
		if len(n.BindIP) != 0 {
			return errors.New("bind ip only can set by one of the flags or configuration file")
		}

		n.BindIP = ip.String()
		return nil
	}

	return nil
}

// trySetDefault set the network's default value if user not configured.
func (n *Network) trySetDefault() {
	if len(n.BindIP) == 0 {
		n.BindIP = "127.0.0.1"
	}
}

// validate network options
func (n Network) validate() error {
	if len(n.BindIP) == 0 {
		return errors.New("network bindIP is not set")
	}

	if ip := net.ParseIP(n.BindIP); ip == nil {
		return errors.New("invalid network bindIP")
	}

	if err := n.TLS.validate(); err != nil {
		return fmt.Errorf("network tls, %v", err)
	}

	return nil
}

// TLSConfig defines tls related options.
type TLSConfig struct {
	// Server should be accessed without verifying the TLS certificate.
	// For testing only.
	InsecureSkipVerify bool `yaml:"insecureSkipVerify"`
	// Server requires TLS client certificate authentication
	CertFile string `yaml:"certFile"`
	// Server requires TLS client certificate authentication
	KeyFile string `yaml:"keyFile"`
	// Trusted root certificates for server
	CAFile string `yaml:"caFile"`
	// the password to decrypt the certificate
	Password string `yaml:"password"`
}

// Enable test tls if enable.
func (tls TLSConfig) Enable() bool {
	if len(tls.CertFile) == 0 &&
		len(tls.KeyFile) == 0 &&
		len(tls.CAFile) == 0 {
		return false
	}

	return true
}

// validate tls configs
func (tls TLSConfig) validate() error {
	if !tls.Enable() {
		return nil
	}

	// TODO: add tls config validate.

	return nil
}

// SysOption is the system's normal option, which is parsed from
// flag commandline.
type SysOption struct {
	ConfigFile string
	// BindIP Setting startup bind ip.
	BindIP net.IP
	// Versioned Setting if show current version info.
	Versioned bool
}

// CheckV check if show current version info.
func (s SysOption) CheckV() {
	if s.Versioned {
		version.ShowVersion()
		os.Exit(0)
	}
}

// IAM defines all the iam related runtime.
type IAM struct {
	// Endpoints is a seed list of host:port addresses of iam nodes.
	Endpoints []string `yaml:"endpoints"`
	// AppCode blueking belong to hcm's appcode.
	AppCode string `yaml:"appCode"`
	// AppSecret blueking belong to hcm app's secret.
	AppSecret string    `yaml:"appSecret"`
	TLS       TLSConfig `yaml:"tls"`
}

// validate iam runtime.
func (s IAM) validate() error {
	if len(s.Endpoints) == 0 {
		return errors.New("iam endpoints is not set")
	}

	if len(s.AppCode) == 0 {
		return errors.New("iam appcode is not set")
	}

	if len(s.AppSecret) == 0 {
		return errors.New("iam app secret is not set")
	}

	if err := s.TLS.validate(); err != nil {
		return fmt.Errorf("iam tls validate failed, err: %v", err)
	}

	return nil
}

// Web 服务依赖所需特有配置， 包括登录、静态文件等配置的定义
type Web struct {
	StaticFileDirPath string `yaml:"staticFileDirPath"`

	BkLoginCookieName string `yaml:"bkLoginCookieName"`
	BkLoginUrl        string `yaml:"bkLoginUrl"`
	BkComponentApiUrl string `yaml:"bkComponentApiUrl"`
	BkItsmUrl         string `yaml:"bkItsmUrl"`
}

func (s Web) validate() error {
	if len(s.BkLoginUrl) == 0 {
		return errors.New("bk_login_url is not set")
	}

	if len(s.BkComponentApiUrl) == 0 {
		return errors.New("bk_component_api_url is not set")
	}

	if len(s.BkItsmUrl) == 0 {
		return errors.New("bk_itsm_url is not set")
	}

	return nil
}

// Esb defines the esb related runtime.
type Esb struct {
	// Endpoints is a seed list of host:port addresses of esb nodes.
	Endpoints []string `yaml:"endpoints"`
	// AppCode is the BlueKing app code of hcm to request esb.
	AppCode string `yaml:"appCode"`
	// AppSecret is the BlueKing app secret of hcm to request esb.
	AppSecret string `yaml:"appSecret"`
	// User is the BlueKing user of hcm to request esb.
	User string    `yaml:"user"`
	TLS  TLSConfig `yaml:"tls"`
}

// validate esb runtime.
func (s Esb) validate() error {
	if len(s.Endpoints) == 0 {
		return errors.New("esb endpoints is not set")
	}
	if len(s.AppCode) == 0 {
		return errors.New("esb app code is not set")
	}
	if len(s.AppSecret) == 0 {
		return errors.New("esb app secret is not set")
	}
	if len(s.User) == 0 {
		return errors.New("esb user is not set")
	}
	if err := s.TLS.validate(); err != nil {
		return fmt.Errorf("validate esb tls failed, err: %v", err)
	}
	return nil
}

// AesGcm Aes Gcm加密
type AesGcm struct {
	Key   string `yaml:"key"`
	Nonce string `yaml:"nonce"`
}

func (a AesGcm) validate() error {
	if len(a.Key) != 16 && len(a.Key) != 32 {
		return errors.New("invalid key, should be 16 or 32 bytes")
	}

	if len(a.Nonce) != 12 {
		return errors.New("invalid nonce, should be 12 bytes")
	}

	return nil
}

// Crypto 定义项目里需要用到的加密，包括选择的算法等
// TODO: 这里默认只支持AES Gcm算法，后续需要支持国密等的选择，可能还需要支持根据不同场景配置不同（比如不同场景，加密的密钥等都不一样）
type Crypto struct {
	AesGcm AesGcm `yaml:"aesGcm"`
}

func (c Crypto) validate() error {
	if err := c.AesGcm.validate(); err != nil {
		return err
	}

	return nil
}

// CloudResource 云资源配置
type CloudResource struct {
	Sync CloudResourceSync `yaml:"sync"`
}

func (c CloudResource) validate() error {
	if err := c.Sync.validate(); err != nil {
		return err
	}

	return nil
}

// CloudResourceSync 云资源同步配置
type CloudResourceSync struct {
	Enable                       bool   `yaml:"enable"`
	SyncIntervalMin              uint64 `yaml:"syncIntervalMin"`
	SyncFrequencyLimitingTimeMin uint64 `yaml:"syncFrequencyLimitingTimeMin"`
}

func (c CloudResourceSync) validate() error {
	if c.Enable {
		if c.SyncFrequencyLimitingTimeMin < 10 {
			return errors.New("syncFrequencyLimitingTimeMin must > 10")
		}
	}

	return nil
}

// Recycle configuration.
type Recycle struct {
	AutoDeleteTime uint `yaml:"autoDeleteTimeHour"`
}

func (a Recycle) validate() error {
	if a.AutoDeleteTime == 0 {
		return errors.New("autoDeleteTimeHour must > 0")
	}

	return nil
}

// BillConfig 账号账单配置
type BillConfig struct {
	Enable          bool   `yaml:"enable"`
	SyncIntervalMin uint64 `yaml:"syncIntervalMin"`
}

func (c BillConfig) validate() error {
	if c.Enable && c.SyncIntervalMin < 1 {
		return errors.New("BillConfig.SyncIntervalMin must >= 1")
	}

	return nil
}
