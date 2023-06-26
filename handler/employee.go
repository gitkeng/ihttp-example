package handler

import (
	"fmt"
	"ihttp-example/datastore"
	"ihttp-example/request"
	"ihttp-example/response"
	"net/http"
	"strconv"

	"github.com/gitkeng/ihttp"
	"github.com/gitkeng/ihttp/util/uuid"
	"github.com/jinzhu/copier"
)

const (
	tagCreateEmployeeHandler = "CreateEmployeeHandler"
	tagFilterEmployeeHandler = "FilterEmployeeHandler"
)

func CreateEmployeeHandler(ctx ihttp.IContext) error {
	req := &request.CreateEmployeeRequest{}
	if err := ctx.Bind(req); err != nil {
		return ctx.Response(
			ihttp.ErrorLevel,
			tagCreateEmployeeHandler,
			http.StatusBadRequest,
			strconv.Itoa(http.StatusBadRequest),
			"invalid req",
			err,
		)
	}

	// copy req
	empEntity := &datastore.Employee{}
	if err := copier.Copy(empEntity, req); err != nil {
		return ctx.Response(
			ihttp.ErrorLevel,
			tagCreateEmployeeHandler,
			http.StatusInternalServerError,
			strconv.Itoa(http.StatusInternalServerError),
			err.Error(),
			err,
		)
	}

	// set another require field
	empEntity.EmployeeCode = uuid.NewUUID()

	if emp, err := datastore.InsertEmployee(ctx, empEntity); err != nil {
		return ctx.Response(
			ihttp.ErrorLevel,
			tagCreateEmployeeHandler,
			http.StatusInternalServerError,
			strconv.Itoa(http.StatusInternalServerError),
			"insert employee failed",
			err,
		)
	} else {
		//generate response
		resp := &response.EmployeeResponse{}
		if err := copier.Copy(resp, emp); err != nil {
			return ctx.Response(
				ihttp.ErrorLevel,
				tagCreateEmployeeHandler,
				http.StatusInternalServerError,
				strconv.Itoa(http.StatusInternalServerError),
				fmt.Sprintf("generate response fail for employee %s", emp.EmployeeCode),
				err,
			)
		}

		return ctx.Response(
			ihttp.DebugLevel,
			tagCreateEmployeeHandler,
			http.StatusOK,
			strconv.Itoa(http.StatusOK),
			"create employee success",
			nil,
			ihttp.Field{
				Key:   "employee",
				Value: resp,
			},
		)
	}

}

func FilterEmployeeHandler(ctx ihttp.IContext) error {
	req := &ihttp.FilterRequest{}
	if err := ctx.Bind(req); err != nil {
		return ctx.Response(
			ihttp.ErrorLevel,
			tagFilterEmployeeHandler,
			http.StatusBadRequest,
			strconv.Itoa(http.StatusBadRequest),
			"invalid req",
			err,
		)
	}

	if employees, total, err := datastore.GetEmployees(ctx, req.GetFilters(), req.GetOption()); err != nil {

		return ctx.Response(
			ihttp.ErrorLevel,
			tagFilterEmployeeHandler,
			http.StatusInternalServerError,
			strconv.Itoa(http.StatusInternalServerError),
			"filter employee failed",
			err,
		)
	} else {

		//generate response response
		resp := make([]*response.EmployeeResponse, 0)
		for idx, _ := range employees {
			emp := &response.EmployeeResponse{}
			if err := copier.Copy(emp, employees[idx]); err != nil {
				return ctx.Response(
					ihttp.ErrorLevel,
					tagFilterEmployeeHandler,
					http.StatusInternalServerError,
					strconv.Itoa(http.StatusInternalServerError),
					fmt.Sprintf("generate response fail for employee %s", employees[idx].EmployeeCode),
					err,
				)
			}
			resp = append(resp, emp)
		}

		return ctx.Response(
			ihttp.DebugLevel,
			tagFilterEmployeeHandler,
			http.StatusOK,
			strconv.Itoa(http.StatusOK),
			"filter employee success",
			nil,
			ihttp.Field{
				Key:   "employees",
				Value: resp,
			},
			ihttp.Field{
				Key:   "total",
				Value: total,
			},
		)
	}

}
