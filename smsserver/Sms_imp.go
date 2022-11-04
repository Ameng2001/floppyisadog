package main

import (
	"context"
)

// SmsImp servant implementation
type SmsImp struct {
}

// Init servant init
func (imp *SmsImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *SmsImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *SmsImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *SmsImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
