package main

import (
	"github.com/gitkeng/ihttp"
	"ihttp-example/handler"
)

const (
	CreateEmployeeEndPoint = "/v1/employee"
	FilterEmployeeEndPoint = "/v1/employee/filter"
)

func registerEmployeeHandler(ms ihttp.IMicroservice) {
	ms.POST(CreateEmployeeEndPoint, handler.CreateEmployeeHandler)
	ms.POST(FilterEmployeeEndPoint, handler.FilterEmployeeHandler)
}
