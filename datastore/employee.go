package datastore

import (
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gitkeng/ihttp"
	"github.com/gitkeng/ihttp/util/convutil"
	"github.com/gitkeng/ihttp/util/dateutil"
	"github.com/gitkeng/ihttp/util/stringutil"
)

const (
	PostgresContextName = "pgdb"
)

type Employee struct {
	EmployeeCode string  `json:"employee_code"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Age          int     `json:"age"`
	Email        string  `json:"email"`
	Department   string  `json:"department"`
	Salary       float64 `json:"salary"`
	UpdateTime   int64   `json:"update_time"`
}

func (emp *Employee) String() string {
	return stringutil.Json(*emp)
}

func (emp *Employee) ToMap() map[string]any {
	return convutil.Obj2Map(*emp)
}
func InsertEmployee(ctx ihttp.IContext, emp *Employee) (*Employee, error) {
	if emp == nil {
		return nil, errors.New("employee is require")
	}
	updateTime := dateutil.GetCurrentEpochTime()
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insertBuilder := builder.Insert("employee").Columns(
		"employee_code",
		"first_name",
		"last_name",
		"email",
		"age",
		"department",
		"salary",
		"update_time",
	).Values(
		emp.EmployeeCode,
		emp.FirstName,
		emp.LastName,
		emp.Email,
		emp.Age,
		emp.Department,
		emp.Salary,
		updateTime)

	sqlCmd, values, err := insertBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	ctx.Logger().Debugf("InsertEmployee Sql Command: %s", sqlCmd)

	dbStore, found := ctx.DB(PostgresContextName)
	if !found {
		return nil, fmt.Errorf("db context name %s not found", PostgresContextName)
	}

	txn, err := dbStore.Conn().Begin()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()
	stmt, err := txn.Prepare(sqlCmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	recordCount := int64(0)
	result, err := stmt.Exec(values...)
	if err != nil {
		txn.Rollback()
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		txn.Rollback()
		return nil, err
	}
	recordCount, err = result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if recordCount == 0 {
		return nil, fmt.Errorf("no record created")
	}

	return &Employee{
		EmployeeCode: emp.EmployeeCode,
		FirstName:    emp.FirstName,
		LastName:     emp.LastName,
		Email:        emp.Email,
		Age:          emp.Age,
		Department:   emp.Department,
		Salary:       emp.Salary,
		UpdateTime:   updateTime,
	}, nil

}

func GetEmployees(ctx ihttp.IContext,
	filters []ihttp.IQueryFilter,
	option ihttp.IQueryOption) ([]*Employee, int64, error) {

	searchBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	totalBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	searchQuery := searchBuilder.Select(
		"employee_code",
		"first_name",
		"last_name",
		"email",
		"age",
		"department",
		"salary",
		"update_time",
	).From("employee")
	totalQuery := totalBuilder.Select("count(employee_code) as total").From("employee")

	for _, filter := range filters {
		switch filter.GetField() {
		case "employee_code":
			searchQuery = searchQuery.Where(
				squirrel.And{
					squirrel.Eq{filter.GetField(): filter.GetValue()},
				})
			totalQuery = totalQuery.Where(
				squirrel.And{
					squirrel.Eq{filter.GetField(): filter.GetValue()},
				})
		case "first_name", "last_name", "email", "department":
			searchQuery = searchQuery.Where(
				squirrel.And{
					squirrel.ILike{filter.GetField(): fmt.Sprintf("%%%s%%", filter.GetValue())},
				})
			totalQuery = totalQuery.Where(
				squirrel.And{
					squirrel.ILike{filter.GetField(): fmt.Sprintf("%%%s%%", filter.GetValue())},
				})
		case "salary", "age", "update_time":
			if filter.GetFromValue() != nil && filter.GetToValue() == nil {
				if filterValue, ok := filter.GetFromValue().(float64); ok && filterValue >= 0 {
					searchQuery = searchQuery.Where(squirrel.And{
						squirrel.GtOrEq{filter.GetField(): filterValue}})
					totalQuery = totalQuery.Where(squirrel.And{
						squirrel.GtOrEq{filter.GetField(): filterValue}})
				}
			}
			if filter.GetFromValue() == nil && filter.GetToValue() != nil {
				if filterValue, ok := filter.GetFromValue().(float64); ok && filterValue >= 0 {
					searchQuery = searchQuery.Where(squirrel.And{
						squirrel.LtOrEq{filter.GetField(): filterValue}})
					totalQuery = totalQuery.Where(squirrel.And{
						squirrel.LtOrEq{filter.GetField(): filterValue}})
				}
			}
			if filter.GetFromValue() != nil && filter.GetToValue() != nil {
				filterFromValue, fromValueOk := filter.GetFromValue().(float64)
				filterToValue, toValueOK := filter.GetToValue().(float64)
				if fromValueOk && toValueOK && filterFromValue <= filterToValue {
					searchQuery = searchQuery.Where(squirrel.And{
						squirrel.GtOrEq{filter.GetField(): filterFromValue},
						squirrel.LtOrEq{filter.GetField(): filterToValue}})
					totalQuery = totalQuery.Where(squirrel.And{
						squirrel.GtOrEq{filter.GetField(): filterFromValue},
						squirrel.LtOrEq{filter.GetField(): filterToValue}})
				}
			}
		}
	}

	if option != nil {
		if option.GetLimit() > 0 {
			searchQuery = searchQuery.Limit(uint64(option.GetLimit()))
		}
		if option.GetOffset() > 0 {
			searchQuery = searchQuery.Offset(uint64(option.GetOffset()))
		}

		//sorting
		if len(option.GetSort()) > 0 {
			orders := make([]string, 0)
			for _, order := range option.GetSort() {
				orderStr := fmt.Sprintf("%s %s", order.GetField(), string(order.GetOrder()))
				orders = append(orders, orderStr)
			}
			searchQuery = searchQuery.OrderBy(orders...)
		}
	}
	sqlCmd, values, err := searchQuery.ToSql()
	if err != nil {
		return nil, -1, err
	}
	ctx.Logger().Debugf("GetEmployees Search SQL: %s", sqlCmd)
	ctx.Logger().Debugf("GetEmployees values: %v", values)

	dbStore, found := ctx.DB(PostgresContextName)
	if !found {
		return nil, -1, fmt.Errorf("db context name %s not found", PostgresContextName)
	}

	stmt, err := dbStore.Conn().Prepare(sqlCmd)
	if err != nil {
		return nil, -1, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()
	results := make([]*Employee, 0)
	for rows.Next() {
		result := &Employee{}
		if err := rows.Scan(
			&result.EmployeeCode,
			&result.FirstName,
			&result.LastName,
			&result.Email,
			&result.Age,
			&result.Department,
			&result.Salary,
			&result.UpdateTime); err != nil {
			return nil, -1, err
		}
		results = append(results, result)
	}

	totalResult, err := totalQuery.RunWith(dbStore.Conn()).Query()
	if err != nil {
		return nil, -1, err
	}
	defer totalResult.Close()
	var total int64
	for totalResult.Next() {
		if err := totalResult.Scan(&total); err != nil {
			return nil, -1, err
		}
	}

	return results, total, nil
}
