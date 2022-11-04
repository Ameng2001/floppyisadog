package main

import (
	"context"
)

// CompanyImp servant implementation
type CompanyImp struct {
}

// Init servant init
func (imp *CompanyImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *CompanyImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *CompanyImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *CompanyImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
