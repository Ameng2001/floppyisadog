package main

import (
	"context"
)

// foauthImp servant implementation
type foauthImp struct {
}

// Init servant init
func (imp *foauthImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *foauthImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *foauthImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *foauthImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
