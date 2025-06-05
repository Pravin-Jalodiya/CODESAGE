require (
	github.com/aws/aws-lambda-go v1.36.1
	github.com/aws/aws-sdk-go-v2/config v1.29.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.39.6
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.18.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.32.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.34.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.55 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.16.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.29 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.29 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.24.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.10 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

module codesage

go 1.21

toolchain go1.22.5
