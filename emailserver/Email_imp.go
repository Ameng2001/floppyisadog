package main

import (
	"context"

	"github.com/floppyisadog/emailserver/tars-protocol/emailserver"
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

func (imp *EmailImp) Send(ctx context.Context, req *emailserver.EmailRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
