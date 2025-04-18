package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/table"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var daoTmpl = templates.EditableHeaderTmpl + `package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"{{.EntityPackage}}"
	"github.com/unionj-cloud/toolkit/caller"
	"github.com/unionj-cloud/toolkit/sqlext/query"
	"github.com/unionj-cloud/toolkit/sqlext/wrapper"
	"github.com/unionj-cloud/toolkit/reflectutils"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/toolkit/templateutils"
	"strings"
	"math"
	"time"
)

var _ I{{.EntityName}}Dao = (*{{.EntityName}}Dao)(nil)

type {{.EntityName}}Dao struct {
	db wrapper.GddDB
}

func New{{.EntityName}}Dao(querier wrapper.GddDB) *{{.EntityName}}Dao {
	return &{{.EntityName}}Dao{
		db: querier,
	}
}

func (receiver *{{.EntityName}}Dao) BeforeSaveHook(ctx context.Context, data *entity.{{.EntityName}}) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) BeforeBulkSaveHook(ctx context.Context, data []*entity.{{.EntityName}}) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) AfterSaveHook(ctx context.Context, data *entity.{{.EntityName}}, lastInsertID int64, affected int64) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) AfterBulkSaveHook(ctx context.Context, data []*entity.{{.EntityName}}, lastInsertID int64, affected int64) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) BeforeUpdateManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) AfterUpdateManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where, affected int64) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) BeforeDeleteManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) AfterDeleteManyHook(ctx context.Context, data []*entity.{{.EntityName}}, where *query.Where, affected int64) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) BeforeReadManyHook(ctx context.Context, page *query.Page, where *query.Where) {
	// implement your business logic
}

func (receiver *{{.EntityName}}Dao) Insert(ctx context.Context, data *entity.{{.EntityName}}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Insert{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		{{- if eq .PkField.Type "int64"}}
		data.{{.PkField.Name}} = lastInsertID
		{{- else }}
		data.{{.PkField.Name}} = {{.PkField.Type}}(lastInsertID)
		{{- end }}
	}
	{{- end }}
	if affected, err = result.RowsAffected(); err == nil {
		{{- if .PkCol.Autoincrement }}
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
		{{- else }}
		receiver.AfterSaveHook(ctx, data, 0, affected)
		{{- end }}
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) InsertIgnore(ctx context.Context, data *entity.{{.EntityName}}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "InsertIgnore{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) BulkInsert(ctx context.Context, data []*entity.{{.EntityName}}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
		affected     int64
	)
	receiver.BeforeBulkSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Insert{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		for i, item :=range data {
			{{- if eq .PkField.Type "int64"}}
			item.{{.PkField.Name}} = lastInsertID + int64(i)
			{{- else }}
			item.{{.PkField.Name}} = {{.PkField.Type}}(lastInsertID) + {{.PkField.Type}}(i)
			{{- end }}
		}
	}
	{{- end }}
	if affected, err = result.RowsAffected(); err == nil {
		{{- if .PkCol.Autoincrement }}
		receiver.AfterBulkSaveHook(ctx, data, lastInsertID, affected)
		{{- else }}
		receiver.AfterBulkSaveHook(ctx, data, 0, affected)
		{{- end }}
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) BulkInsertIgnore(ctx context.Context, data []*entity.{{.EntityName}}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		affected     int64
	)
	receiver.BeforeBulkSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "InsertIgnore{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterBulkSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

// Upsert With ON DUPLICATE KEY UPDATE, the affected-rows value per row is 1 if the row is inserted as a new row,
// 2 if an existing row is updated, and 0 if an existing row is set to its current values.
// If you specify the CLIENT_FOUND_ROWS flag to the mysql_real_connect() C API function when connecting to mysqld,
// the affected-rows value is 1 (not 0) if an existing row is set to its current values.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (receiver *{{.EntityName}}Dao) Upsert(ctx context.Context, data *entity.{{.EntityName}}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Upsert{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		{{- if eq .PkField.Type "int64"}}
		data.{{.PkField.Name}} = lastInsertID
		{{- else }}
		data.{{.PkField.Name}} = {{.PkField.Type}}(lastInsertID)
		{{- end }}
	}
	{{- end }}
	if affected, err = result.RowsAffected(); err == nil {
		{{- if .PkCol.Autoincrement }}
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
		{{- else }}
		receiver.AfterSaveHook(ctx, data, 0, affected)
		{{- end }}
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) BulkUpsert(ctx context.Context, data []*entity.{{.EntityName}}) (int64, error) {
	var (
		statement    string
		updateClause string
		err          error
		result       sql.Result
		affected     int64
		args      []interface{}
	)
	receiver.BeforeBulkSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Insert{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	statement, args, err = receiver.db.BindNamed(statement, data)
	if err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if updateClause, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "UpdateClause{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	statement += "\n" + updateClause
	if result, err = receiver.db.ExecContext(ctx, statement, args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterBulkSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) BulkUpsertSelect(ctx context.Context, data []*entity.{{.EntityName}}, columns []string) (int64, error) {
	var (
		statement    string
		updateClause string
		err          error
		result       sql.Result
		affected     int64
		args      []interface{}
	)
	receiver.BeforeBulkSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Insert{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	statement, args, err = receiver.db.BindNamed(statement, data)
	if err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if updateClause, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "UpdateClauseSelect{{.EntityName}}", struct {
		Columns []string
	}{
		Columns: columns,
	}); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	statement += "\n" + updateClause
	if result, err = receiver.db.ExecContext(ctx, statement, args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterBulkSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) UpsertNoneZero(ctx context.Context, data *entity.{{.EntityName}}) (int64, error) {
	var (
		statement    string
		err          error
		result       sql.Result
		{{- if .PkCol.Autoincrement }}
		lastInsertID int64
		{{- end }}
		affected     int64
	)
	receiver.BeforeSaveHook(ctx, data)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(entity.{{.EntityName}}); !ok {
		return 0, errors.New("underlying type of data should be entity.{{.EntityName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Upsert{{.EntityName}}NoneZero", data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	{{- if .PkCol.Autoincrement }}
	if lastInsertID, err = result.LastInsertId(); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if lastInsertID > 0 {
		{{- if eq .PkField.Type "int64"}}
		data.{{.PkField.Name}} = lastInsertID
		{{- else }}
		data.{{.PkField.Name}} = {{.PkField.Type}}(lastInsertID)
		{{- end }}
	}
	{{- end }}
	if affected, err = result.RowsAffected(); err == nil {
		{{- if .PkCol.Autoincrement }}
		receiver.AfterSaveHook(ctx, data, lastInsertID, affected)
		{{- else }}
		receiver.AfterSaveHook(ctx, data, 0, affected)
		{{- end }}
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) DeleteMany(ctx context.Context, where query.Where) (int64, error) {
	var (
		err    error
		result sql.Result
		w      string
		args   []interface{}
		affected int64
	)
	receiver.BeforeDeleteManyHook(ctx, nil, &where)
	w, args = where.Sql()
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("delete from {{.TableName}} where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterDeleteManyHook(ctx, nil, &where, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) Update(ctx context.Context, data *entity.{{.EntityName}}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		affected  int64
	)
	receiver.BeforeSaveHook(ctx, data)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Update{{.EntityName}}", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) UpdateNoneZero(ctx context.Context, data *entity.{{.EntityName}}) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		affected  int64
	)
	receiver.BeforeSaveHook(ctx, data)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(entity.{{.EntityName}}); !ok {
		return 0, errors.New("underlying type of data should be entity.{{.EntityName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Update{{.EntityName}}NoneZero", data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if result, err = receiver.db.NamedExecContext(ctx, statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterSaveHook(ctx, data, 0, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) UpdateMany(ctx context.Context, data []*entity.{{.EntityName}}, where query.Where) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		q         string
		args      []interface{}
		wargs     []interface{}
		w         string
		affected  int64
	)
	receiver.BeforeUpdateManyHook(ctx, data, &where)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Update{{.EntityName}}s", nil); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if q, args, err = receiver.db.BindNamed(statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	w, wargs = where.Sql()
	if stringutils.IsNotEmpty(w) {
		q += " where " + w
	}
	args = append(args, wargs...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(q), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterUpdateManyHook(ctx, data, &where, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) UpdateManyNoneZero(ctx context.Context, data []*entity.{{.EntityName}}, where query.Where) (int64, error) {
	var (
		statement string
		err       error
		result    sql.Result
		q         string
		args      []interface{}
		wargs     []interface{}
		w         string
		affected  int64
	)
	receiver.BeforeUpdateManyHook(ctx, data, &where)
	value := reflectutils.ValueOf(data).Interface()
	if _, ok := value.(entity.{{.EntityName}}); !ok {
		return 0, errors.New("underlying type of data should be entity.{{.EntityName}}")
	}
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Update{{.EntityName}}sNoneZero", data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if q, args, err = receiver.db.BindNamed(statement, data); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	w, wargs = where.Sql()
	if stringutils.IsNotEmpty(w) {
		q += " where " + w
	}
	args = append(args, wargs...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(q), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterUpdateManyHook(ctx, data, &where, affected)
	}
	return affected, err
}

func (receiver *{{.EntityName}}Dao) Get(ctx context.Context, dest *entity.{{.EntityName}}, id {{.PkField.Type}}) error {
	var (
		statement string
		err       error
		{{.EntityName | ToLower}}      entity.{{.EntityName}}
	)
	if statement, err = templateutils.BlockMysql("{{.EntityName | ToLower}}dao.sql", {{.EntityName | ToLower}}daosql, "Get{{.EntityName}}", nil); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	if err = receiver.db.GetContext(ctx, &{{.EntityName | ToLower}}, receiver.db.Rebind(statement), id); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	return nil
}

func (receiver *{{.EntityName}}Dao) SelectMany(ctx context.Context, dest *[]entity.{{.EntityName}}, where query.Where) error {
	var (
		statements []string
		err       error
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, nil, &where)
    statements = append(statements, "select * from {{.TableName}}")
	if !where.IsEmpty() {
		statements = append(statements, "where")
		q, wargs := where.Sql()
		statements = append(statements, q)
		args = append(args, wargs...)
	}
	sqlStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.SelectContext(ctx, dest, receiver.db.Rebind(sqlStr), args...); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	return nil
}

func (receiver *{{.EntityName}}Dao) CountMany(ctx context.Context, where query.Where) (int, error) {
	var (
		statements []string
		err       error
		total     int
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, nil, &where)
	statements = append(statements, "select count(1) from {{.TableName}}")
    if !where.IsEmpty() {
		statements = append(statements, "where")
		q, wargs := where.Sql()
		statements = append(statements, q)
		args = append(args, wargs...)
	}
	sqlStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.GetContext(ctx, &total, receiver.db.Rebind(sqlStr), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	return total, nil
}

type {{.EntityName}}PageRet struct {
	Items    []entity.{{.EntityName}}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}

func (receiver *{{.EntityName}}Dao) PageMany(ctx context.Context, dest *{{.EntityName}}PageRet, page query.Page, where query.Where) error {
	var (
		statements []string
		err       error
		args       []interface{}
	)
	receiver.BeforeReadManyHook(ctx, &page, &where)
	statements = append(statements, "select * from {{.TableName}}")
    if !where.IsEmpty() {
		statements = append(statements, "where")
		q, wargs := where.Sql()
		statements = append(statements, q)
		args = append(args, wargs...)
	}
	p, pargs := page.Sql()
	statements = append(statements, p)
	args = append(args, pargs...)
	sqlStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.SelectContext(ctx, &dest.Items, receiver.db.Rebind(sqlStr), args...); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}
	
	args = nil
    statements = nil
	statements = append(statements, "select count(1) from {{.TableName}}")
    if !where.IsEmpty() {
		statements = append(statements, "where")
		q, wargs := where.Sql()
		statements = append(statements, q)
		args = append(args, wargs...)
	}
	sqlStr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(strings.Join(statements, " ")), "where"))
	if err = receiver.db.GetContext(ctx, &dest.Total, receiver.db.Rebind(sqlStr), args...); err != nil {
		return errors.Wrap(err, caller.NewCaller().String())
	}

	pageNo := 1
	if page.Size > 0 {
		pageNo = page.Offset/page.Size + 1
	}
	dest.PageNo = pageNo
	dest.PageSize = page.Size
	if dest.PageSize > 0 && math.Ceil(float64(dest.Total)/float64(dest.PageSize)) > float64(dest.PageNo) {
		dest.HasNext = true
	}
	return nil
}

func (receiver *{{.EntityName}}Dao) DeleteManySoft(ctx context.Context, where query.Where) (int64, error) {
	var (
		err      error
		result   sql.Result
		w        string
		args     []interface{}
		affected int64
	)
	receiver.BeforeDeleteManyHook(ctx, nil, &where)
	w, args = where.Sql()
	args = append([]interface{}{time.Now()}, args...)
	if result, err = receiver.db.ExecContext(ctx, receiver.db.Rebind(fmt.Sprintf("update {{.TableName}} set delete_at=? where %s;", w)), args...); err != nil {
		return 0, errors.Wrap(err, caller.NewCaller().String())
	}
	if affected, err = result.RowsAffected(); err == nil {
		receiver.AfterDeleteManyHook(ctx, nil, &where, affected)
	}
	return affected, err
}`

// GenDaoGo generates dao layer implementation code
func GenDaoGo(entityPath string, t table.Table, folder ...string) error {
	var (
		err      error
		dpkg     string
		daopath  string
		f        *os.File
		funcMap  map[string]interface{}
		tpl      *template.Template
		pkColumn table.Column
		df       string
	)
	df = "dao"
	if len(folder) > 0 {
		df = folder[0]
	}
	daopath = filepath.Join(filepath.Dir(entityPath), df)
	_ = os.MkdirAll(daopath, os.ModePerm)

	daofile := filepath.Join(daopath, strings.ToLower(t.Meta.Name)+"dao.go")
	if _, err = os.Stat(daofile); os.IsNotExist(err) {
		f, _ = os.Create(daofile)
		defer f.Close()

		dpkg = astutils.GetImportPath(entityPath)
		funcMap = make(map[string]interface{})
		funcMap["ToLower"] = strings.ToLower
		funcMap["ToSnake"] = strcase.ToSnake
		tpl, _ = template.New("dao.go.tmpl").Funcs(funcMap).Parse(daoTmpl)
		for _, column := range t.Columns {
			if column.Pk {
				pkColumn = column
				break
			}
		}
		_ = tpl.Execute(f, struct {
			EntityPackage string
			EntityName    string
			TableName     string
			PkField       astutils.FieldMeta
			PkCol         table.Column
			Version       string
		}{
			EntityPackage: dpkg,
			EntityName:    t.Meta.Name,
			TableName:     t.Name,
			PkField:       pkColumn.Meta,
			PkCol:         pkColumn,
			Version:       version.Release,
		})
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}
