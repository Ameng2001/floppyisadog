package consts

const (
	CookieName = "staffjoy-gateway"
	// uuidKey       = "uuid"
	// supportKey    = "support"
	// expirationKey = "exp"
	// for GRPC
	CurrentUserMetadata = "gateway-current-user-uuid"
	// header set for internal user id
	CurrentUserHeader = "Grpc-Metadata-Gateway-Current-User-Uuid"

	// AuthorizationHeader is the http request header
	// key used for accessing the internal authorization.
	AuthorizationHeader = "Authorization"

	// AuthorizationMetadata is the grpce metadadata key used
	// for accessing the internal authorization
	AuthorizationMetadata = "authorization"

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
