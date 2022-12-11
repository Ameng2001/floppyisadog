package consts

const (
	CookieName = "staffjoy-gateway"
	//和tars verify接口一致(
	//1. X-Verify-UID： verify 接口返回的 uid；
	//2. X-Verify-Data：verify 接口返回的 context，可以为空；
	//3. X-Verify-Token：身份票据 token 信息，即请求 verify 接口中的 token。)
	// 网关中配置filterheaders = X-Verify-UID|X-Verify-Data
	// tars context 中存储uuid，在网关的verify回调中设置
	CurrentUserMetadata = "X-Verify-UID"
	// http request header中存储uuid，在webportal的middleware中设置
	CurrentUserHeader = "X-Verify-UID"

	// AuthorizationHeader is the http request header
	// key used for accessing the internal authorization.
	// http request header中存储用于授权的身份信息，在webportal的middleware中设置
	AuthorizationHeader = "X-Verify-Data"
	// AuthorizationMetadata is the grpce metadadata key used
	// for accessing the internal authorization
	// tars context中存储用于授权的身份信息，在网关的verify回调中设置
	AuthorizationMetadata = "X-Verify-Data"

	// AuthorizationAnonymousWeb is set as the Authorization header to denote that
	// a request is being made bu an unauthenticated web user
	AuthorizationAnonymousWeb = "gateway-anonymous"

	// AuthorizationAuthenticatedUser is set as the Authorization header to denote that
	// a request is being made by an authenticated web user
	AuthorizationAuthenticatedUser = "gateway-authenticated"

	// AuthorizationSupportUser is set as the Authorization header to denote that
	// a request is being made by a Staffjoy team me
	AuthorizationSupportUser = "gateway-support"

	// AuthorizationWWWService is set as the Authorization header to denote that
	// a request is being made by the www login / signup system
	AuthorizationWWWService = "www-service"

	// AuthorizationCompanyService is set as the Authorization header to denote
	// that a request is being made by the company api/server
	AuthorizationCompanyService = "company-service"

	// AuthorizationSuperpowersService is set as the Authorization header to
	// denote that a request is being made by the dev-only superpowers service
	AuthorizationSuperpowersService = "superpowers-service"

	// AuthorizationWhoamiService is set as the Authorization heade to denote that
	// a request is being made by the whoami microservice
	AuthorizationWhoamiService = "whoami-service"

	// AuthorizationBotService is set as the Authorization header to denote that
	// a request is being made by the bot microservice
	AuthorizationBotService = "bot-service"

	// AuthorizationAccountService is set as the Authorization header to denote that
	// a request is being made by the account service
	AuthorizationAccountService = "account-service"

	// AuthorizationICalService is set as the Authorization header to denote that
	// a request is being made by the ical service
	AuthorizationICalService = "ical-service"
)
