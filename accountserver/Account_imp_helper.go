package main

import (
	"fmt"
	"net/url"

	"github.com/floppyisadog/accountserver/managers/logmgr"
	"github.com/floppyisadog/accountserver/models"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/utils/crypto"
	"github.com/floppyisadog/appcommon/utils/environment"
)

const (
	resetPasswordTmpl   = "<div>We received a request to reset the password on your account. To do so, click the below link. If you did not request this change, no action is needed. <br/> <a href=\"%s\">%s</a></div>"
	activateAccountTmpl = "<div><p>Hi %s, and welcome to Staffjoy!</p><a href=\"%s\">Please click here to finish setting up your account.</a></p></div><br/><br/><div>If you have trouble clicking on the link, please copy and paste this link into your browser: <br/><a href=\"%s\">%s</a></div>"
	confirmEmailTmpl    = "<div>Hi %s!</div>To confirm your new email address, <a href=\"%s\">please click here</a>.</div><br/><br/><div>If you have trouble clicking on the link, please copy and paste this link into your browser: <br/><a href=\"%s\">%s</a></div>"
)

func (imp *AccountImp) syncUser(uuid string) {
	//TODO
}

func (imp *AccountImp) sendActiveEmail(accountDAO *models.Account) error {
	if len(accountDAO.Email) <= 0 {
		return codes.ErrSendActiveEmail
	}

	//generate a activate token
	logmgr.RINFO("Begin generate active token\n")
	token, err := crypto.EmailConfirmationToken(accountDAO.ID, accountDAO.Email, imp.config.SigningToken)
	if err != nil {
		logmgr.RERROR("Gnetate active token error:(%v)\n", err)
		return err
	}
	logmgr.RINFO("End genetate active token:(%s)\n", token)

	link := url.URL{Host: "www." + environment.GetCurrEnv().ExternalApex, Path: fmt.Sprintf("/activate/%s", token), Scheme: "http"}
	// Send verification email
	emailName := accountDAO.Name
	if emailName == "" {
		emailName = "there"
	}
	emailBody := fmt.Sprintf(activateAccountTmpl, emailName, link.String(), link.String(), link.String())
	logmgr.RINFO("Send active mail:(%s)\n", emailBody)
	logmgr.DINFOF("Send active mail:(%s)\n", emailBody)

	//call email rpc service
	md := make(map[string]string)
	md[consts.AuthorizationMetadata] = consts.AuthorizationAccountService

	//TODO
	// ret, err := imp.emailPrx.SendWithContext(
	// 	context.Background(),
	// 	&emailserver.EmailRequest{
	// 		To:        accountDAO.Email,
	// 		Name:      accountDAO.Name,
	// 		Subject:   "Activate your Staffjoy account",
	// 		Html_body: emailBody,
	// 	},
	// 	md,
	// )
	// if ret != codes.OK {
	// 	logmgr.RERROR("Failed to call sendemail - %v", err)
	// 	return codes.ErrSendActiveEmail
	// }

	return nil
}
