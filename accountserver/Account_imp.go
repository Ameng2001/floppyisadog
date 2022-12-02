package main

import (
	"context"
	"strings"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/floppyisadog/accountserver/managers/configmgr"
	"github.com/floppyisadog/accountserver/models"
	"github.com/floppyisadog/accountserver/repos"
	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils"
	"github.com/floppyisadog/appcommon/utils/crypto"
	"github.com/floppyisadog/appcommon/utils/database"
	"github.com/floppyisadog/emailserver/tars-protocol/emailserver"
	"github.com/floppyisadog/smsserver/tars-protocol/smsserver"
	"github.com/jinzhu/gorm"
)

// AccountImp servant implementation
type AccountImp struct {
	config   *configmgr.Config
	db       *gorm.DB
	smsPrx   *smsserver.Sms
	emailPrx *emailserver.Email
}

// Init servant init
func (imp *AccountImp) Init() error {
	//initialize servant here:
	imp.config = configmgr.GetConfig()
	imp.db = database.GetDB()

	comm := tars.NewCommunicator()
	smsObj := imp.config.Outerfactory["SmsObj"]
	imp.smsPrx = new(smsserver.Sms)
	comm.StringToProxy(smsObj, imp.smsPrx)

	emailObj := imp.config.Outerfactory["EmailObj"]
	imp.emailPrx = new(emailserver.Email)
	comm.StringToProxy(emailObj, imp.emailPrx)

	return nil
}

// Destroy servant destroy
func (imp *AccountImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *AccountImp) Create(ctx context.Context, req *accountserver.CreateAccountRequest, rsp *accountserver.AccountInfo) (int32, error) {
	_, authz, err := helpers.GetAuth(ctx)
	if err != nil {
		return codes.Unknown, codes.ErrAuthorizedFailed
	}

	ok := imp.authorizeClient(authz)
	if !ok {
		return codes.PermissionDenied, codes.ErrPermissionDenied
	}

	if len(req.Email)+len(req.Phonenumber)+len(req.Name) == 0 {
		return codes.InvalidArgument, codes.ErrInvalidArgument
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if len(req.Email) > 0 && !strings.Contains(req.Email, "@") {
		return codes.InvalidArgument, codes.ErrInvalidArgument
	}

	req.Phonenumber, err = utils.ParseAndFormatPhonenumber(req.Phonenumber)
	if err != nil {
		return codes.InvalidArgument, codes.ErrInvalidArgument
	}

	if req.Email != "" {
		_, nofound := repos.FindAccountByEmail(req.Email)
		if !nofound {
			return codes.AlreadyExists, codes.ErrAccountAlreadyExists
		}
	}

	if req.Phonenumber != "" {
		_, nofound := repos.FindAccountByPhonenumber(req.Phonenumber)
		if !nofound {
			return codes.AlreadyExists, codes.ErrAccountAlreadyExists
		}
	}

	uuid, err := crypto.NewUUID()
	if err != nil {
		return codes.Unknown, codes.ErrGenerateUUID
	}

	accountDAO := &models.Account{
		BaseGormModel: models.BaseGormModel{
			ID:        uuid.String(),
			CreatedAt: time.Now().UTC(),
		},
		Email:       req.Email,
		Name:        req.Name,
		PhoneNumber: req.Phonenumber,
		PhotoUrl:    utils.GenerateGravatarURL(req.Email),
		MemberSince: time.Now().UTC(),
	}

	if err := imp.db.Create(accountDAO).Error; err != nil {
		return codes.Unknown, codes.ErrCreateAccount
	}

	go imp.syncUser(accountDAO.ID)

	if err := imp.sendActiveEmail(accountDAO); err != nil {
		return codes.Unknown, codes.ErrSendActiveEmail
	}

	//TODO
	//AuditLog

	return codes.OK, nil
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

func (imp *AccountImp) authorizeClient(authz string) bool {
	switch authz {
	case consts.AuthorizationSupportUser:
	case consts.AuthorizationWWWService:
	case consts.AuthorizationCompanyService:
	default:
		return false
	}

	return true
}
