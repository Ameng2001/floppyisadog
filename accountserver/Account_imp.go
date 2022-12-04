package main

import (
	"context"
	"strings"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/floppyisadog/accountserver/managers/configmgr"
	"github.com/floppyisadog/accountserver/managers/logmgr"
	"github.com/floppyisadog/accountserver/models"
	"github.com/floppyisadog/accountserver/repos"
	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils"
	"github.com/floppyisadog/appcommon/utils/crypto"
	"github.com/floppyisadog/appcommon/utils/database"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/emailserver/tars-protocol/emailserver"
	"github.com/floppyisadog/smsserver/tars-protocol/smsserver"
	"github.com/jinzhu/gorm"
)

const (
	minPasswordLength = 6
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

	switch authz {
	case consts.AuthorizationSupportUser:
	case consts.AuthorizationWWWService:
	case consts.AuthorizationCompanyService:
	default:
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
	uid := strings.Replace(uuid.String(), "-", "", -1) //去掉uuid中的-，否则mysql查询不出
	if err != nil {
		return codes.Unknown, codes.ErrGenerateUUID
	}

	accountDAO := &models.Account{
		BaseGormModel: models.BaseGormModel{
			ID:        uid,
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
	rsp.Uuid = accountDAO.ID
	rsp.Name = accountDAO.Name
	rsp.Email = accountDAO.Email
	rsp.Confirmed_and_active = accountDAO.ConfirmAndActive
	rsp.Member_since = accountserver.Timestamp{
		Seconds: int64(accountDAO.MemberSince.Second()),
		Nanos:   int32(accountDAO.MemberSince.Nanosecond()),
	}
	rsp.Support = accountDAO.Support
	rsp.Phonenumber = accountDAO.PhoneNumber
	rsp.Photo_url = accountDAO.PhotoUrl

	return codes.OK, nil
}
func (imp *AccountImp) List(ctx context.Context, req *accountserver.GetAccountListRequest, rsp *accountserver.AccountList) (int32, error) {
	//Doing something in your function
	//TODO
	return 0, nil
}
func (imp *AccountImp) Get(ctx context.Context, req *accountserver.GetAccountRequest, rsp *accountserver.AccountInfo) (int32, error) {
	_, authz, err := helpers.GetAuth(ctx)
	if err != nil {
		return codes.Unknown, codes.ErrAuthorizedFailed
	}
	switch authz {
	case consts.AuthorizationWWWService:
	case consts.AuthorizationAccountService:
	case consts.AuthorizationCompanyService:
	case consts.AuthorizationWhoamiService:
	case consts.AuthorizationBotService:
	case consts.AuthorizationAuthenticatedUser:
		uuid, err := helpers.GetCurrentUserUUIDFromMetadata(ctx)
		if err != nil {
			return codes.PermissionDenied, err
		}
		if uuid != req.Uuid {
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
	case consts.AuthorizationSupportUser:
	case consts.AuthorizationSuperpowersService:
		if environment.GetCurrEnv().Name != "development" {
			logmgr.RWARN("Development service trying to connect outside development environment")
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
	default:
		return codes.PermissionDenied, codes.ErrPermissionDenied
	}

	if req.Uuid == "" {
		logmgr.RERROR("lack of uuid in the request")
		return codes.InvalidArgument, codes.ErrInvalidArgument
	}

	accountDAO, nofound := repos.FindAccountByUUID(req.Uuid)
	if nofound {
		return codes.NotFound, codes.ErrAccountNotFound
	}

	// TODO: 有没有优雅的方法DAO to DTO
	rsp.Uuid = accountDAO.ID
	rsp.Name = accountDAO.Name
	rsp.Email = accountDAO.Email
	rsp.Confirmed_and_active = accountDAO.ConfirmAndActive
	rsp.Member_since = accountserver.Timestamp{
		Seconds: int64(accountDAO.MemberSince.Second()),
		Nanos:   int32(accountDAO.MemberSince.Nanosecond()),
	}
	rsp.Support = accountDAO.Support
	rsp.Phonenumber = accountDAO.PhoneNumber
	rsp.Photo_url = accountDAO.PhotoUrl

	return codes.OK, nil
}
func (imp *AccountImp) Update(ctx context.Context, req *accountserver.AccountInfo) (int32, error) {
	_, authz, err := helpers.GetAuth(ctx)
	if err != nil {
		logmgr.RERROR("can not get authz\n")
		return codes.Unknown, codes.ErrAuthorizedFailed
	}
	switch authz {
	case consts.AuthorizationWWWService:
	case consts.AuthorizationCompanyService:
	case consts.AuthorizationAuthenticatedUser:
		uuid, err := helpers.GetCurrentUserUUIDFromMetadata(ctx)
		if err != nil {
			logmgr.RERROR("can not get uuid for metadata\n")
			return codes.PermissionDenied, err
		}
		if uuid != req.Uuid {
			logmgr.RERROR("uuid mismatch\n")
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
	case consts.AuthorizationSupportUser:
	case consts.AuthorizationSuperpowersService:
		if environment.GetCurrEnv().Name != "development" {
			logmgr.RWARN("Development service trying to connect outside development environment")
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
	default:
		return codes.PermissionDenied, codes.ErrPermissionDenied
	}

	existing, nofound := repos.FindAccountByUUID(req.Uuid)
	if nofound {
		// This handles 404 and everything!
		logmgr.RERROR("cannot get account by uuid\n")
		return codes.NotFound, codes.ErrAccountNotFound
	}

	// Some validations
	if req.Phonenumber, err = utils.ParseAndFormatPhonenumber(req.Phonenumber); err != nil {
		logmgr.RERROR("format phonenumber error\n")
		return codes.InvalidArgument, codes.ErrInvalidArgument
	}
	existing.PhoneNumber = req.Phonenumber

	if req.Member_since.Seconds != int64(existing.MemberSince.Second()) {
		logmgr.RERROR("You cannot modify the member_since date\n")
		return codes.PermissionDenied, codes.ErrPermissionDenied
	}

	req.Email = strings.ToLower(req.Email)
	if req.Email != "" && (req.Email != existing.Email) {
		// Check to see if account exists
		_, nofound := repos.FindAccountByEmail(req.Email)
		if !nofound {
			logmgr.RERROR("A user with that email already exists. Try a password reset\n")
			return codes.AlreadyExists, codes.ErrAccountAlreadyExists
		}
	}
	existing.Email = req.Email

	if req.Phonenumber != "" && (req.Phonenumber != existing.PhoneNumber) {
		_, nofound := repos.FindAccountByPhonenumber(req.Phonenumber)
		if !nofound {
			logmgr.RERROR("A user with that phonenumber already exists. Try a password reset.\n")
			return codes.AlreadyExists, codes.ErrAccountAlreadyExists
		}
	}
	existing.PhoneNumber = req.Phonenumber

	if authz == consts.AuthorizationAuthenticatedUser {
		if (req.Confirmed_and_active != existing.ConfirmAndActive) && (!existing.ConfirmAndActive) {
			logmgr.RERROR("You cannot activate this account.\n")
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
		if req.Support != existing.Support {
			logmgr.RERROR("You cannot change the support parameter.\n")
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
		if req.Photo_url != existing.PhotoUrl {
			logmgr.RERROR("You cannot change the photo through this endpoint (see docs).\n")
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
	}

	existing.Name = req.Name
	existing.ConfirmAndActive = req.Confirmed_and_active
	existing.Support = req.Support
	existing.PhotoUrl = utils.GenerateGravatarURL(req.Email)

	if ok := repos.UpdateAccount(existing); !ok {
		logmgr.RERROR("Could not update the user account.\n")
		return codes.Unknown, codes.ErrUpdateAccountError
	}

	go imp.syncUser(existing.ID)

	// If account is being activated, or if phone number is changed by current user - send text
	if req.Confirmed_and_active && len(req.Phonenumber) > 0 && req.Phonenumber != existing.PhoneNumber {
		imp.sendSmsGreeting(req.Phonenumber)
	}

	return codes.OK, nil
}
func (imp *AccountImp) UpdatePassword(ctx context.Context, req *accountserver.UpdatePasswordRequest) (int32, error) {
	_, authz, err := helpers.GetAuth(ctx)
	if err != nil {
		logmgr.RERROR("can not get authz\n")
		return codes.Unknown, codes.ErrAuthorizedFailed
	}
	switch authz {
	case consts.AuthorizationAuthenticatedUser:
		uuid, err := helpers.GetCurrentUserUUIDFromMetadata(ctx)
		if err != nil {
			logmgr.RERROR("can not get uuid for metadata\n")
			return codes.PermissionDenied, err
		}
		if uuid != req.Uuid {
			logmgr.RERROR("uuid mismatch\n")
			return codes.PermissionDenied, codes.ErrPermissionDenied
		}
	case consts.AuthorizationWWWService:
	case consts.AuthorizationSupportUser:
	default:
		return codes.PermissionDenied, codes.ErrPermissionDenied
	}

	// Verify inputs
	if req.Uuid == "" {
		return codes.InvalidArgument, codes.ErrInvalidArgument
	}
	if len(req.Password) < minPasswordLength {
		return codes.InvalidArgument, codes.ErrInvalidArgument
	}
	salt, err := crypto.NewSalt()
	if err != nil {
		logmgr.RERROR("Failed to generate a salt")
		return codes.Unknown, codes.ErrInternal
	}

	pwHash, err := crypto.HashPassword(salt, []byte(req.Password))
	if err != nil {
		logmgr.RERROR("Failed to hash the password")
		return codes.Unknown, codes.ErrInternal
	}

	// Run the update . . .
	if ok := repos.UpdateAccountFields(req.Uuid, map[string]interface{}{"password_hash": pwHash}); !ok {
		logmgr.RERROR("Failed to read the database.\n")
		return codes.Unknown, codes.ErrUpdateAccountError
	}

	return codes.OK, nil
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
