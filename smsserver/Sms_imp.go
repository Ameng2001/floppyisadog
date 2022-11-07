package main

import (
	"context"

	"github.com/floppyisadog/smsserver/tars-protocol/smsserver"
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

func (imp *SmsImp) QueueSend(ctx context.Context, req *smsserver.SmsRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
