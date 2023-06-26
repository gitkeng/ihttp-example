package response

import (
	"github.com/gitkeng/ihttp/util/convutil"
	"github.com/gitkeng/ihttp/util/stringutil"
)

type EmployeeResponse struct {
	EmployeeCode string  `json:"employee_code"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Age          int     `json:"age"`
	Email        string  `json:"email"`
	Department   string  `json:"department"`
	Salary       float64 `json:"salary"`
}

func (emp *EmployeeResponse) String() string {
	return stringutil.Json(*emp)
}

func (emp *EmployeeResponse) ToMap() map[string]any {
	return convutil.Obj2Map(*emp)
}
