package main

import (
	"context"
)

// BotImp servant implementation
type BotImp struct {
}

// Init servant init
func (imp *BotImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *BotImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *BotImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
func (imp *BotImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
