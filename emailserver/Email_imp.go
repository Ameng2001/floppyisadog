package main

import (
	"context"
)

// EmailImp servant implementation
type EmailImp struct {
}

// Init servant init
func (imp *EmailImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *EmailImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *EmailImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *EmailImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
