module github.com/jennyservices/jenny

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/certifi/gocertifi v0.0.0-20180905225744-ee1a9a0726d2 // indirect
	github.com/cloudflare/roughtime v0.0.0-20181123203841-94176ac0b23c
	github.com/d4l3k/messagediff v1.2.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.2.0
	github.com/go-kit/kit v0.8.0
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-openapi/inflect v0.17.2
	github.com/go-openapi/loads v0.17.2
	github.com/go-openapi/spec v0.17.2
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/gddo v0.0.0-20181116215533-9bd4a3295021
	github.com/gorilla/mux v1.6.2 // indirect
	github.com/gorilla/schema v1.0.2
	github.com/hashicorp/hcl v1.0.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.0.0
	github.com/oklog/ulid v1.3.1
	github.com/opentracing/opentracing-go v1.0.2
	github.com/pkg/errors v0.8.0
	github.com/sergi/go-diff v1.0.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
	golang.org/x/tools v0.0.0-20181213151202-c779628d65d9
	google.golang.org/grpc v1.17.0
	gopkg.in/square/go-jose.v2 v2.3.1
	roughtime.googlesource.com/roughtime.git v0.0.0-20190418172256-51f6971f5f06
	sevki.org/x v1.0.0
)

replace roughtime.googlesource.com/go/client => roughtime.googlesource.com/roughtime.git/go/client v0.0.0-20190418172256-51f6971f5f06
