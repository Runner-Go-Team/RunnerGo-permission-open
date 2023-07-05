// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"permission-open/internal/pkg/dal/model"
)

func newThirdNoticeGroupRelate(db *gorm.DB, opts ...gen.DOOption) thirdNoticeGroupRelate {
	_thirdNoticeGroupRelate := thirdNoticeGroupRelate{}

	_thirdNoticeGroupRelate.thirdNoticeGroupRelateDo.UseDB(db, opts...)
	_thirdNoticeGroupRelate.thirdNoticeGroupRelateDo.UseModel(&model.ThirdNoticeGroupRelate{})

	tableName := _thirdNoticeGroupRelate.thirdNoticeGroupRelateDo.TableName()
	_thirdNoticeGroupRelate.ALL = field.NewAsterisk(tableName)
	_thirdNoticeGroupRelate.ID = field.NewInt64(tableName, "id")
	_thirdNoticeGroupRelate.GroupID = field.NewString(tableName, "group_id")
	_thirdNoticeGroupRelate.NoticeID = field.NewString(tableName, "notice_id")
	_thirdNoticeGroupRelate.Params = field.NewString(tableName, "params")
	_thirdNoticeGroupRelate.CreatedAt = field.NewTime(tableName, "created_at")
	_thirdNoticeGroupRelate.UpdatedAt = field.NewTime(tableName, "updated_at")
	_thirdNoticeGroupRelate.DeletedAt = field.NewField(tableName, "deleted_at")

	_thirdNoticeGroupRelate.fillFieldMap()

	return _thirdNoticeGroupRelate
}

type thirdNoticeGroupRelate struct {
	thirdNoticeGroupRelateDo thirdNoticeGroupRelateDo

	ALL       field.Asterisk
	ID        field.Int64  // 主键id
	GroupID   field.String // 通知组id
	NoticeID  field.String // 通知id
	Params    field.String // 通知目标参数
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field

	fieldMap map[string]field.Expr
}

func (t thirdNoticeGroupRelate) Table(newTableName string) *thirdNoticeGroupRelate {
	t.thirdNoticeGroupRelateDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t thirdNoticeGroupRelate) As(alias string) *thirdNoticeGroupRelate {
	t.thirdNoticeGroupRelateDo.DO = *(t.thirdNoticeGroupRelateDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *thirdNoticeGroupRelate) updateTableName(table string) *thirdNoticeGroupRelate {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.GroupID = field.NewString(table, "group_id")
	t.NoticeID = field.NewString(table, "notice_id")
	t.Params = field.NewString(table, "params")
	t.CreatedAt = field.NewTime(table, "created_at")
	t.UpdatedAt = field.NewTime(table, "updated_at")
	t.DeletedAt = field.NewField(table, "deleted_at")

	t.fillFieldMap()

	return t
}

func (t *thirdNoticeGroupRelate) WithContext(ctx context.Context) *thirdNoticeGroupRelateDo {
	return t.thirdNoticeGroupRelateDo.WithContext(ctx)
}

func (t thirdNoticeGroupRelate) TableName() string { return t.thirdNoticeGroupRelateDo.TableName() }

func (t thirdNoticeGroupRelate) Alias() string { return t.thirdNoticeGroupRelateDo.Alias() }

func (t *thirdNoticeGroupRelate) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *thirdNoticeGroupRelate) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 7)
	t.fieldMap["id"] = t.ID
	t.fieldMap["group_id"] = t.GroupID
	t.fieldMap["notice_id"] = t.NoticeID
	t.fieldMap["params"] = t.Params
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
	t.fieldMap["deleted_at"] = t.DeletedAt
}

func (t thirdNoticeGroupRelate) clone(db *gorm.DB) thirdNoticeGroupRelate {
	t.thirdNoticeGroupRelateDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t thirdNoticeGroupRelate) replaceDB(db *gorm.DB) thirdNoticeGroupRelate {
	t.thirdNoticeGroupRelateDo.ReplaceDB(db)
	return t
}

type thirdNoticeGroupRelateDo struct{ gen.DO }

func (t thirdNoticeGroupRelateDo) Debug() *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Debug())
}

func (t thirdNoticeGroupRelateDo) WithContext(ctx context.Context) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t thirdNoticeGroupRelateDo) ReadDB() *thirdNoticeGroupRelateDo {
	return t.Clauses(dbresolver.Read)
}

func (t thirdNoticeGroupRelateDo) WriteDB() *thirdNoticeGroupRelateDo {
	return t.Clauses(dbresolver.Write)
}

func (t thirdNoticeGroupRelateDo) Session(config *gorm.Session) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Session(config))
}

func (t thirdNoticeGroupRelateDo) Clauses(conds ...clause.Expression) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t thirdNoticeGroupRelateDo) Returning(value interface{}, columns ...string) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t thirdNoticeGroupRelateDo) Not(conds ...gen.Condition) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t thirdNoticeGroupRelateDo) Or(conds ...gen.Condition) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t thirdNoticeGroupRelateDo) Select(conds ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t thirdNoticeGroupRelateDo) Where(conds ...gen.Condition) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t thirdNoticeGroupRelateDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *thirdNoticeGroupRelateDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t thirdNoticeGroupRelateDo) Order(conds ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t thirdNoticeGroupRelateDo) Distinct(cols ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t thirdNoticeGroupRelateDo) Omit(cols ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t thirdNoticeGroupRelateDo) Join(table schema.Tabler, on ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t thirdNoticeGroupRelateDo) LeftJoin(table schema.Tabler, on ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t thirdNoticeGroupRelateDo) RightJoin(table schema.Tabler, on ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t thirdNoticeGroupRelateDo) Group(cols ...field.Expr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t thirdNoticeGroupRelateDo) Having(conds ...gen.Condition) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t thirdNoticeGroupRelateDo) Limit(limit int) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t thirdNoticeGroupRelateDo) Offset(offset int) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t thirdNoticeGroupRelateDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t thirdNoticeGroupRelateDo) Unscoped() *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Unscoped())
}

func (t thirdNoticeGroupRelateDo) Create(values ...*model.ThirdNoticeGroupRelate) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t thirdNoticeGroupRelateDo) CreateInBatches(values []*model.ThirdNoticeGroupRelate, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t thirdNoticeGroupRelateDo) Save(values ...*model.ThirdNoticeGroupRelate) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t thirdNoticeGroupRelateDo) First() (*model.ThirdNoticeGroupRelate, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThirdNoticeGroupRelate), nil
	}
}

func (t thirdNoticeGroupRelateDo) Take() (*model.ThirdNoticeGroupRelate, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThirdNoticeGroupRelate), nil
	}
}

func (t thirdNoticeGroupRelateDo) Last() (*model.ThirdNoticeGroupRelate, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThirdNoticeGroupRelate), nil
	}
}

func (t thirdNoticeGroupRelateDo) Find() ([]*model.ThirdNoticeGroupRelate, error) {
	result, err := t.DO.Find()
	return result.([]*model.ThirdNoticeGroupRelate), err
}

func (t thirdNoticeGroupRelateDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ThirdNoticeGroupRelate, err error) {
	buf := make([]*model.ThirdNoticeGroupRelate, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t thirdNoticeGroupRelateDo) FindInBatches(result *[]*model.ThirdNoticeGroupRelate, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t thirdNoticeGroupRelateDo) Attrs(attrs ...field.AssignExpr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t thirdNoticeGroupRelateDo) Assign(attrs ...field.AssignExpr) *thirdNoticeGroupRelateDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t thirdNoticeGroupRelateDo) Joins(fields ...field.RelationField) *thirdNoticeGroupRelateDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t thirdNoticeGroupRelateDo) Preload(fields ...field.RelationField) *thirdNoticeGroupRelateDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t thirdNoticeGroupRelateDo) FirstOrInit() (*model.ThirdNoticeGroupRelate, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThirdNoticeGroupRelate), nil
	}
}

func (t thirdNoticeGroupRelateDo) FirstOrCreate() (*model.ThirdNoticeGroupRelate, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThirdNoticeGroupRelate), nil
	}
}

func (t thirdNoticeGroupRelateDo) FindByPage(offset int, limit int) (result []*model.ThirdNoticeGroupRelate, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t thirdNoticeGroupRelateDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t thirdNoticeGroupRelateDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t thirdNoticeGroupRelateDo) Delete(models ...*model.ThirdNoticeGroupRelate) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *thirdNoticeGroupRelateDo) withDO(do gen.Dao) *thirdNoticeGroupRelateDo {
	t.DO = *do.(*gen.DO)
	return t
}