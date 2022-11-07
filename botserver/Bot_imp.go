package main

import (
	"context"

	"github.com/floppyisadog/botserver/tars-protocol/botserver"
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

func (imp *BotImp) OnboardWorker(ctx context.Context, req *botserver.OnboardWorkerRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *BotImp) AlertNewShift(ctx context.Context, req *botserver.AlertNewShiftRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *BotImp) AlertNewShifts(ctx context.Context, req *botserver.AlertNewShiftsRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *BotImp) AlertRemovedShift(ctx context.Context, req *botserver.AlertRemovedShiftRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *BotImp) AlertRemovedShifts(ctx context.Context, req *botserver.AlertRemovedShiftsRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *BotImp) AlertChangedShift(ctx context.Context, req *botserver.AlertChangedShiftRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
