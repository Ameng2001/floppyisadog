package main

import (
	"context"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
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

func (imp *AccountImp) Create(ctx context.Context, req *accountserver.CreateAccountRequest, rsp *accountserver.AccountInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) List(ctx context.Context, req *accountserver.GetAccountListRequest, rsp *accountserver.AccountList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) Get(ctx context.Context, req *accountserver.GetAccountRequest, rsp *accountserver.AccountInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) Update(ctx context.Context, req *accountserver.AccountInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) UpdatePassword(ctx context.Context, req *accountserver.UpdatePasswordRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) RequestPasswordReset(ctx context.Context, req *accountserver.PasswordResetRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) RequestEmailChange(ctx context.Context, req *accountserver.EmailChangeRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) VerifyPassword(ctx context.Context, req *accountserver.VerifyPasswordRequest, rsp *accountserver.AccountInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) ChangeEmail(ctx context.Context, req *accountserver.EmailConfirmation) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) GetOrCreate(ctx context.Context, req *accountserver.GetOrCreateRequest, rsp *accountserver.AccountInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) GetAccountByPhonenumber(ctx context.Context, req *accountserver.GetAccountByPhonenumberRequest, rsp *accountserver.AccountInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) TrackEvent(ctx context.Context, req *accountserver.TrackEventRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) SyncUser(ctx context.Context, req *accountserver.SyncUserRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
