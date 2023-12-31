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

func newStressPlan(db *gorm.DB, opts ...gen.DOOption) stressPlan {
	_stressPlan := stressPlan{}

	_stressPlan.stressPlanDo.UseDB(db, opts...)
	_stressPlan.stressPlanDo.UseModel(&model.StressPlan{})

	tableName := _stressPlan.stressPlanDo.TableName()
	_stressPlan.ALL = field.NewAsterisk(tableName)
	_stressPlan.ID = field.NewInt64(tableName, "id")
	_stressPlan.PlanID = field.NewString(tableName, "plan_id")
	_stressPlan.TeamID = field.NewString(tableName, "team_id")
	_stressPlan.RankID = field.NewInt64(tableName, "rank_id")
	_stressPlan.PlanName = field.NewString(tableName, "plan_name")
	_stressPlan.TaskType = field.NewInt32(tableName, "task_type")
	_stressPlan.TaskMode = field.NewInt32(tableName, "task_mode")
	_stressPlan.Status = field.NewInt32(tableName, "status")
	_stressPlan.CreateUserID = field.NewString(tableName, "create_user_id")
	_stressPlan.RunUserID = field.NewString(tableName, "run_user_id")
	_stressPlan.Remark = field.NewString(tableName, "remark")
	_stressPlan.RunCount = field.NewInt64(tableName, "run_count")
	_stressPlan.CreatedAt = field.NewTime(tableName, "created_at")
	_stressPlan.UpdatedAt = field.NewTime(tableName, "updated_at")
	_stressPlan.DeletedAt = field.NewField(tableName, "deleted_at")

	_stressPlan.fillFieldMap()

	return _stressPlan
}

type stressPlan struct {
	stressPlanDo stressPlanDo

	ALL          field.Asterisk
	ID           field.Int64  // 主键ID
	PlanID       field.String // 计划ID
	TeamID       field.String // 团队ID
	RankID       field.Int64  // 序号ID
	PlanName     field.String // 计划名称
	TaskType     field.Int32  // 计划类型：1-普通任务，2-定时任务
	TaskMode     field.Int32  // 压测类型: 1-并发模式，2-阶梯模式，3-错误率模式，4-响应时间模式，5-每秒请求数模式，6-每秒事务数模式
	Status       field.Int32  // 计划状态1:未开始,2:进行中
	CreateUserID field.String // 创建人id
	RunUserID    field.String // 运行人id
	Remark       field.String // 备注
	RunCount     field.Int64  // 运行次数
	CreatedAt    field.Time   // 创建时间
	UpdatedAt    field.Time   // 修改时间
	DeletedAt    field.Field  // 删除时间

	fieldMap map[string]field.Expr
}

func (s stressPlan) Table(newTableName string) *stressPlan {
	s.stressPlanDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s stressPlan) As(alias string) *stressPlan {
	s.stressPlanDo.DO = *(s.stressPlanDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *stressPlan) updateTableName(table string) *stressPlan {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.PlanID = field.NewString(table, "plan_id")
	s.TeamID = field.NewString(table, "team_id")
	s.RankID = field.NewInt64(table, "rank_id")
	s.PlanName = field.NewString(table, "plan_name")
	s.TaskType = field.NewInt32(table, "task_type")
	s.TaskMode = field.NewInt32(table, "task_mode")
	s.Status = field.NewInt32(table, "status")
	s.CreateUserID = field.NewString(table, "create_user_id")
	s.RunUserID = field.NewString(table, "run_user_id")
	s.Remark = field.NewString(table, "remark")
	s.RunCount = field.NewInt64(table, "run_count")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")
	s.DeletedAt = field.NewField(table, "deleted_at")

	s.fillFieldMap()

	return s
}

func (s *stressPlan) WithContext(ctx context.Context) *stressPlanDo {
	return s.stressPlanDo.WithContext(ctx)
}

func (s stressPlan) TableName() string { return s.stressPlanDo.TableName() }

func (s stressPlan) Alias() string { return s.stressPlanDo.Alias() }

func (s *stressPlan) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *stressPlan) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 15)
	s.fieldMap["id"] = s.ID
	s.fieldMap["plan_id"] = s.PlanID
	s.fieldMap["team_id"] = s.TeamID
	s.fieldMap["rank_id"] = s.RankID
	s.fieldMap["plan_name"] = s.PlanName
	s.fieldMap["task_type"] = s.TaskType
	s.fieldMap["task_mode"] = s.TaskMode
	s.fieldMap["status"] = s.Status
	s.fieldMap["create_user_id"] = s.CreateUserID
	s.fieldMap["run_user_id"] = s.RunUserID
	s.fieldMap["remark"] = s.Remark
	s.fieldMap["run_count"] = s.RunCount
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
	s.fieldMap["deleted_at"] = s.DeletedAt
}

func (s stressPlan) clone(db *gorm.DB) stressPlan {
	s.stressPlanDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s stressPlan) replaceDB(db *gorm.DB) stressPlan {
	s.stressPlanDo.ReplaceDB(db)
	return s
}

type stressPlanDo struct{ gen.DO }

func (s stressPlanDo) Debug() *stressPlanDo {
	return s.withDO(s.DO.Debug())
}

func (s stressPlanDo) WithContext(ctx context.Context) *stressPlanDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s stressPlanDo) ReadDB() *stressPlanDo {
	return s.Clauses(dbresolver.Read)
}

func (s stressPlanDo) WriteDB() *stressPlanDo {
	return s.Clauses(dbresolver.Write)
}

func (s stressPlanDo) Session(config *gorm.Session) *stressPlanDo {
	return s.withDO(s.DO.Session(config))
}

func (s stressPlanDo) Clauses(conds ...clause.Expression) *stressPlanDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s stressPlanDo) Returning(value interface{}, columns ...string) *stressPlanDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s stressPlanDo) Not(conds ...gen.Condition) *stressPlanDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s stressPlanDo) Or(conds ...gen.Condition) *stressPlanDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s stressPlanDo) Select(conds ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s stressPlanDo) Where(conds ...gen.Condition) *stressPlanDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s stressPlanDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *stressPlanDo {
	return s.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (s stressPlanDo) Order(conds ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s stressPlanDo) Distinct(cols ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s stressPlanDo) Omit(cols ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s stressPlanDo) Join(table schema.Tabler, on ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s stressPlanDo) LeftJoin(table schema.Tabler, on ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s stressPlanDo) RightJoin(table schema.Tabler, on ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s stressPlanDo) Group(cols ...field.Expr) *stressPlanDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s stressPlanDo) Having(conds ...gen.Condition) *stressPlanDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s stressPlanDo) Limit(limit int) *stressPlanDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s stressPlanDo) Offset(offset int) *stressPlanDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s stressPlanDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *stressPlanDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s stressPlanDo) Unscoped() *stressPlanDo {
	return s.withDO(s.DO.Unscoped())
}

func (s stressPlanDo) Create(values ...*model.StressPlan) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s stressPlanDo) CreateInBatches(values []*model.StressPlan, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s stressPlanDo) Save(values ...*model.StressPlan) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s stressPlanDo) First() (*model.StressPlan, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.StressPlan), nil
	}
}

func (s stressPlanDo) Take() (*model.StressPlan, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.StressPlan), nil
	}
}

func (s stressPlanDo) Last() (*model.StressPlan, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.StressPlan), nil
	}
}

func (s stressPlanDo) Find() ([]*model.StressPlan, error) {
	result, err := s.DO.Find()
	return result.([]*model.StressPlan), err
}

func (s stressPlanDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.StressPlan, err error) {
	buf := make([]*model.StressPlan, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s stressPlanDo) FindInBatches(result *[]*model.StressPlan, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s stressPlanDo) Attrs(attrs ...field.AssignExpr) *stressPlanDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s stressPlanDo) Assign(attrs ...field.AssignExpr) *stressPlanDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s stressPlanDo) Joins(fields ...field.RelationField) *stressPlanDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s stressPlanDo) Preload(fields ...field.RelationField) *stressPlanDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s stressPlanDo) FirstOrInit() (*model.StressPlan, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.StressPlan), nil
	}
}

func (s stressPlanDo) FirstOrCreate() (*model.StressPlan, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.StressPlan), nil
	}
}

func (s stressPlanDo) FindByPage(offset int, limit int) (result []*model.StressPlan, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s stressPlanDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s stressPlanDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s stressPlanDo) Delete(models ...*model.StressPlan) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *stressPlanDo) withDO(do gen.Dao) *stressPlanDo {
	s.DO = *do.(*gen.DO)
	return s
}
