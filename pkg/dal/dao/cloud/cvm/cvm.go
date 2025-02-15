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

package cvm

import (
	"fmt"

	"hcm/pkg/api/core"
	"hcm/pkg/criteria/enumor"
	"hcm/pkg/criteria/errf"
	"hcm/pkg/dal/dao/audit"
	idgenerator "hcm/pkg/dal/dao/id-generator"
	"hcm/pkg/dal/dao/orm"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/dal/dao/types"
	"hcm/pkg/dal/table"
	tableaudit "hcm/pkg/dal/table/audit"
	tablecvm "hcm/pkg/dal/table/cloud/cvm"
	"hcm/pkg/dal/table/utils"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/runtime/filter"

	"github.com/jmoiron/sqlx"
)

// Interface only used for cvm.
type Interface interface {
	BatchCreateWithTx(kt *kit.Kit, tx *sqlx.Tx, models []*tablecvm.Table) ([]string, error)
	Update(kt *kit.Kit, expr *filter.Expression, model *tablecvm.Table) error
	UpdateByIDWithTx(kt *kit.Kit, tx *sqlx.Tx, id string, model *tablecvm.Table) error
	List(kt *kit.Kit, opt *types.ListOption) (*types.ListCvmDetails, error)
	ListWithTx(kt *kit.Kit, tx *sqlx.Tx, opt *types.ListOption) (*types.ListCvmDetails, error)
	DeleteWithTx(kt *kit.Kit, tx *sqlx.Tx, expr *filter.Expression) error
}

var _ Interface = new(Dao)

// Dao cvm dao.
type Dao struct {
	Orm   orm.Interface
	IDGen idgenerator.IDGenInterface
	Audit audit.Interface
}

// BatchCreateWithTx cvm.
func (dao Dao) BatchCreateWithTx(kt *kit.Kit, tx *sqlx.Tx, models []*tablecvm.Table) ([]string, error) {

	ids, err := dao.IDGen.Batch(kt, table.CvmTable, len(models))
	if err != nil {
		return nil, err
	}
	for index, model := range models {
		model.ID = ids[index]

		if err := model.InsertValidate(); err != nil {
			return nil, err
		}
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s)	VALUES(%s)`, table.CvmTable,
		tablecvm.TableColumns.ColumnExpr(), tablecvm.TableColumns.ColonNameExpr())

	if err = dao.Orm.Txn(tx).BulkInsert(kt.Ctx, sql, models); err != nil {
		logs.Errorf("insert %s failed, err: %v, rid: %s", table.CvmTable, err, kt.Rid)
		return nil, fmt.Errorf("insert %s failed, err: %v", table.CvmTable, err)
	}

	// create audit.
	audits := make([]*tableaudit.AuditTable, 0, len(models))
	for _, one := range models {
		audits = append(audits, &tableaudit.AuditTable{
			ResID:      one.ID,
			CloudResID: one.CloudID,
			ResName:    one.Name,
			ResType:    enumor.CvmAuditResType,
			Action:     enumor.Create,
			BkBizID:    one.BkBizID,
			Vendor:     one.Vendor,
			AccountID:  one.AccountID,
			Operator:   kt.User,
			Source:     kt.GetRequestSource(),
			Rid:        kt.Rid,
			AppCode:    kt.AppCode,
			Detail: &tableaudit.BasicDetail{
				Data: one,
			},
		})
	}
	if err = dao.Audit.BatchCreateWithTx(kt, tx, audits); err != nil {
		logs.Errorf("batch create audit failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	return ids, nil
}

// Update cvm.
func (dao Dao) Update(kt *kit.Kit, expr *filter.Expression, model *tablecvm.Table) error {

	if expr == nil {
		return errf.New(errf.InvalidParameter, "filter expr is nil")
	}

	if err := model.UpdateValidate(); err != nil {
		return err
	}

	whereExpr, whereValue, err := expr.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return err
	}

	opts := utils.NewFieldOptions().AddBlankedFields("memo").AddIgnoredFields(types.DefaultIgnoredFields...)
	setExpr, toUpdate, err := utils.RearrangeSQLDataWithOption(model, opts)
	if err != nil {
		return fmt.Errorf("prepare parsed sql set filter expr failed, err: %v", err)
	}

	sql := fmt.Sprintf(`UPDATE %s %s %s`, model.TableName(), setExpr, whereExpr)

	_, err = dao.Orm.AutoTxn(kt, func(txn *sqlx.Tx, opt *orm.TxnOption) (interface{}, error) {
		effected, err := dao.Orm.Txn(txn).Update(kt.Ctx, sql, tools.MapMerge(toUpdate, whereValue))
		if err != nil {
			logs.ErrorJson("update cvm failed, err: %v, filter: %s, rid: %v", err, expr, kt.Rid)
			return nil, err
		}

		if effected == 0 {
			logs.ErrorJson("update cvm, but record not found, filter: %v, rid: %v", expr, kt.Rid)
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateByIDWithTx cvm.
func (dao Dao) UpdateByIDWithTx(kt *kit.Kit, tx *sqlx.Tx, id string, model *tablecvm.Table) error {
	if len(id) == 0 {
		return errf.New(errf.InvalidParameter, "id is required")
	}

	if err := model.UpdateValidate(); err != nil {
		return err
	}

	opts := utils.NewFieldOptions().AddBlankedFields("memo").AddIgnoredFields(types.DefaultIgnoredFields...)
	setExpr, toUpdate, err := utils.RearrangeSQLDataWithOption(model, opts)
	if err != nil {
		return fmt.Errorf("prepare parsed sql set filter expr failed, err: %v", err)
	}

	sql := fmt.Sprintf(`UPDATE %s %s where id = :id`, model.TableName(), setExpr)

	toUpdate["id"] = id
	_, err = dao.Orm.Txn(tx).Update(kt.Ctx, sql, toUpdate)
	if err != nil {
		logs.ErrorJson("update cvm failed, err: %v, id: %s, rid: %v", err, id, kt.Rid)
		return err
	}

	return nil
}

// List cvm.
func (dao Dao) List(kt *kit.Kit, opt *types.ListOption) (*types.ListCvmDetails, error) {
	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "list options is nil")
	}

	columnTypes := tablecvm.TableColumns.ColumnTypes()
	columnTypes["extension.resource_group_name"] = enumor.String
	columnTypes["extension.zones"] = enumor.Json
	if err := opt.Validate(filter.NewExprOption(filter.RuleFields(columnTypes)),
		core.DefaultPageOption); err != nil {
		return nil, err
	}

	whereExpr, whereValue, err := opt.Filter.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return nil, err
	}

	if opt.Page.Count {
		// this is a count request, then do count operation only.
		sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, table.CvmTable, whereExpr)

		count, err := dao.Orm.Do().Count(kt.Ctx, sql, whereValue)
		if err != nil {
			logs.ErrorJson("count cvm failed, err: %v, filter: %s, rid: %s", err, opt.Filter, kt.Rid)
			return nil, err
		}

		return &types.ListCvmDetails{Count: count}, nil
	}

	pageExpr, err := types.PageSQLExpr(opt.Page, types.DefaultPageSQLOption)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s %s %s`, tablecvm.TableColumns.FieldsNamedExpr(opt.Fields),
		table.CvmTable, whereExpr, pageExpr)

	details := make([]tablecvm.Table, 0)
	if err = dao.Orm.Do().Select(kt.Ctx, &details, sql, whereValue); err != nil {
		return nil, err
	}

	return &types.ListCvmDetails{Details: details}, nil
}

// ListWithTx cvm with tx.
func (dao Dao) ListWithTx(kt *kit.Kit, tx *sqlx.Tx, opt *types.ListOption) (*types.ListCvmDetails, error) {
	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "list options is nil")
	}

	columnTypes := tablecvm.TableColumns.ColumnTypes()
	columnTypes["extension.resource_group_name"] = enumor.String
	columnTypes["extension.zones"] = enumor.Json
	if err := opt.Validate(filter.NewExprOption(filter.RuleFields(columnTypes)),
		core.DefaultPageOption); err != nil {
		return nil, err
	}

	whereExpr, whereValue, err := opt.Filter.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return nil, err
	}

	if opt.Page.Count {
		// this is a count request, then do count operation only.
		sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, table.CvmTable, whereExpr)

		count, err := dao.Orm.Txn(tx).Count(kt.Ctx, sql, whereValue)
		if err != nil {
			logs.ErrorJson("count cvm failed, err: %v, filter: %s, rid: %s", err, opt.Filter, kt.Rid)
			return nil, err
		}

		return &types.ListCvmDetails{Count: count}, nil
	}

	pageExpr, err := types.PageSQLExpr(opt.Page, types.DefaultPageSQLOption)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s %s %s`, tablecvm.TableColumns.FieldsNamedExpr(opt.Fields),
		table.CvmTable, whereExpr, pageExpr)

	details := make([]tablecvm.Table, 0)
	if err = dao.Orm.Txn(tx).Select(kt.Ctx, &details, sql, whereValue); err != nil {
		return nil, err
	}

	return &types.ListCvmDetails{Details: details}, nil
}

// DeleteWithTx cvm.
func (dao Dao) DeleteWithTx(kt *kit.Kit, tx *sqlx.Tx, expr *filter.Expression) error {
	if expr == nil {
		return errf.New(errf.InvalidParameter, "filter expr is required")
	}

	whereExpr, whereValue, err := expr.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`DELETE FROM %s %s`, table.CvmTable, whereExpr)
	if _, err = dao.Orm.Txn(tx).Delete(kt.Ctx, sql, whereValue); err != nil {
		logs.ErrorJson("delete cvm failed, err: %v, filter: %s, rid: %s", err, expr, kt.Rid)
		return err
	}

	return nil
}

// ListCvm TODO: 考虑之后这种跨表查询是否可以直接引用对象的 List 函数，而不是再写一个。
func ListCvm(kt *kit.Kit, orm orm.Interface, ids []string) (map[string]tablecvm.Table, error) {
	sql := fmt.Sprintf(`SELECT %s FROM %s where id in (:ids)`, tablecvm.TableColumns.FieldsNamedExpr(nil),
		table.CvmTable)

	cvms := make([]tablecvm.Table, 0)
	if err := orm.Do().Select(kt.Ctx, &cvms, sql, map[string]interface{}{"ids": ids}); err != nil {
		return nil, err
	}

	idCvmMap := make(map[string]tablecvm.Table, len(ids))
	for _, one := range cvms {
		idCvmMap[one.ID] = one
	}

	return idCvmMap, nil
}
