package main

import (
	"context"
)

// AccountImp servant implementation
type AccountImp struct {
}

// Init servant init
func (imp *AccountImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *AccountImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *AccountImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *AccountImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
