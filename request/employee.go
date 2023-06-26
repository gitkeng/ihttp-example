package request

import (
	"fmt"

	"github.com/gitkeng/ihttp"
	"github.com/gitkeng/ihttp/util/convutil"
	"github.com/gitkeng/ihttp/util/stringutil"
)

var (
	err0001 = ihttp.NewError("err0001",
		"first_name is require")
	err0002 = ihttp.NewError("err0002",
		"last_name is require")
	err0003 = ihttp.NewError("err0003",
		"email is require")
	err0004 = ihttp.NewError("err0004",
		"department is require")
	err0005 = func(salary float64) ihttp.Error {
		return ihttp.NewError("err0005", fmt.Sprintf("invalid salary: %.2f", salary))
	}
)

type CreateEmployeeRequest struct {
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	Age        int     `json:"age"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
}

func (emp *CreateEmployeeRequest) String() string {
	return stringutil.Json(*emp)
}

func (emp *CreateEmployeeRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*emp)
}

func (emp *CreateEmployeeRequest) Validate() error {
	var errs ihttp.Errors
	if stringutil.IsEmptyString(emp.FirstName) {
		errs = append(errs, err0001)
	}
	if stringutil.IsEmptyString(emp.LastName) {
		errs = append(errs, err0002)
	}
	if stringutil.IsEmptyString(emp.Email) {
		errs = append(errs, err0003)
	}
	if stringutil.IsEmptyString(emp.Department) {
		errs = append(errs, err0004)
	}
	if emp.Salary <= 0 {
		errs = append(errs, err0005(emp.Salary))
	}
	if len(errs) > 0 {
		return &errs
	}
	return nil
}
