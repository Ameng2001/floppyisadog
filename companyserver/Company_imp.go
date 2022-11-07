package main

import (
	"context"

	"github.com/floppyisadog/companyserver/tars-protocol/companyserver"
)

// CompanyImp servant implementation
type CompanyImp struct {
}

// Init servant init
func (imp *CompanyImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *CompanyImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *CompanyImp) CreateCompany(ctx context.Context, req *companyserver.CreateCompanyRequest, rsp *companyserver.CompanyInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListCompanies(ctx context.Context, req *companyserver.CompanyListRequest, rsp *companyserver.CompanyList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetCompany(ctx context.Context, req *companyserver.GetCompanyRequest, rsp *companyserver.CompanyInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) UpdateCompany(ctx context.Context, req *companyserver.CompanyInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) CreateTeam(ctx context.Context, req *companyserver.CreateTeamRequest, rsp *companyserver.TeamInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListTeams(ctx context.Context, req *companyserver.TeamListRequest, rsp *companyserver.TeamList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetTeam(ctx context.Context, req *companyserver.GetTeamRequest, rsp *companyserver.TeamInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) UpdateTeam(ctx context.Context, req *companyserver.TeamInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetWorkerTeamInfo(ctx context.Context, req *companyserver.Worker, rsp *companyserver.Worker) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) CreateJob(ctx context.Context, req *companyserver.CreateJobRequest, rsp *companyserver.JobInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListJobs(ctx context.Context, req *companyserver.JobListRequest, rsp *companyserver.JobList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetJob(ctx context.Context, req *companyserver.GetJobRequest, rsp *companyserver.JobInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) UpdateJob(ctx context.Context, req *companyserver.JobInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) CreateShift(ctx context.Context, req *companyserver.CreateShiftRequest, rsp *companyserver.ShiftInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListShifts(ctx context.Context, req *companyserver.ShiftListRequest, rsp *companyserver.ShiftList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListWorkerShifts(ctx context.Context, req *companyserver.WorkerShiftListRequest, rsp *companyserver.ShiftList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) BulkPublishShifts(ctx context.Context, req *companyserver.BulkPublishShiftsRequest, rsp *companyserver.ShiftList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetShift(ctx context.Context, req *companyserver.GetShiftRequest, rsp *companyserver.ShiftInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) DeleteShift(ctx context.Context, req *companyserver.GetShiftRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) UpdateShift(ctx context.Context, req *companyserver.ShiftInfo, rsp *companyserver.ShiftInfo) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) CreateDirectory(ctx context.Context, req *companyserver.NewDirectoryEntry, rsp *companyserver.DirectoryEntry) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) Directory(ctx context.Context, req *companyserver.DirectoryListRequest, rsp *companyserver.DirectoryList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetAssociations(ctx context.Context, req *companyserver.DirectoryListRequest, rsp *companyserver.AssociationList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetDirectoryEntry(ctx context.Context, req *companyserver.DirectoryEntryRequest, rsp *companyserver.DirectoryEntry) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) UpdateDirectoryEntry(ctx context.Context, req *companyserver.DirectoryEntry, rsp *companyserver.DirectoryEntry) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListAdmins(ctx context.Context, req *companyserver.AdminListRequest, rsp *companyserver.Admins) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) CreateAdmin(ctx context.Context, req *companyserver.DirectoryEntryRequest, rsp *companyserver.DirectoryEntry) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetAdmin(ctx context.Context, req *companyserver.DirectoryEntryRequest, rsp *companyserver.DirectoryEntry) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) DeleteAdmin(ctx context.Context, req *companyserver.DirectoryEntryRequest) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListWorkers(ctx context.Context, req *companyserver.WorkerListRequest, rsp *companyserver.Workers) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetWorker(ctx context.Context, req *companyserver.Worker, rsp *companyserver.DirectoryEntry) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) DeleteWorker(ctx context.Context, req *companyserver.Worker) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) CreateWorker(ctx context.Context, req *companyserver.Worker, rsp *companyserver.DirectoryEntry) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) ListTimeZones(ctx context.Context, req *companyserver.TimeZoneList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GrowthGraph(ctx context.Context, req *companyserver.GrowthGraphResponse) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetWorkerOf(ctx context.Context, req *companyserver.WorkerOfRequest, rsp *companyserver.WorkerOfList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}

func (imp *CompanyImp) GetAdminOf(ctx context.Context, req *companyserver.AdminOfRequest, rsp *companyserver.AdminOfList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
