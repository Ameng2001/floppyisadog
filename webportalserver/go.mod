module github.com/floppyisadog/webportalserver

go 1.19

require (
	github.com/GeertJohan/go.rice v1.0.3
	github.com/TarsCloud/TarsGo v1.3.6
	github.com/floppyisadog/accountserver v0.0.0-00010101000000-000000000000
	github.com/floppyisadog/appcommon v0.0.0-00010101000000-000000000000
	github.com/floppyisadog/companyserver v0.0.0-00010101000000-000000000000
	github.com/gin-contrib/gzip v0.0.6
	github.com/gin-contrib/sessions v0.0.5
	github.com/gin-gonic/gin v1.8.1
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/utrack/gin-csrf v0.0.0-20190424104817-40fb8d2c8fca
)

require (
	github.com/daaku/go.zipexe v1.0.2 // indirect
	github.com/dchest/uniuri v0.0.0-20160212164326-8902c56451e9 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	golang.org/x/crypto v0.0.0-20210920023735-84f357641f63 // indirect
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf // indirect
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/floppyisadog/appcommon => ../appcommon

replace github.com/floppyisadog/accountserver => ../accountserver

replace github.com/floppyisadog/companyserver => ../companyserver
