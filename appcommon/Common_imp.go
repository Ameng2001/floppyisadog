package main

import (
	"context"
)

// CommonImp servant implementation
type CommonImp struct {
}

// Init servant init
func (imp *CommonImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *CommonImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *CommonImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *CommonImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
