package main

import (
	"context"
)

// WebPortalImp servant implementation
type WebPortalImp struct {
}

// Init servant init
func (imp *WebPortalImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *WebPortalImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *WebPortalImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *WebPortalImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
