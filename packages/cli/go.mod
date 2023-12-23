module github.com/aws/amazon-genomics-cli

go 1.19

require (
	github.com/AlecAivazis/survey/v2 v2.2.15
	github.com/antihax/optional v1.0.0
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/aws/aws-sdk-go-v2 v1.8.1
	github.com/aws/aws-sdk-go-v2/config v1.6.1
	github.com/aws/aws-sdk-go-v2/credentials v1.3.3
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.1.4
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.2.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.3.2
	github.com/aws/aws-sdk-go-v2/service/batch v1.6.0
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.5.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.5.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.4.2
	github.com/aws/aws-sdk-go-v2/service/ecr v1.4.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.11.1
	github.com/aws/aws-sdk-go-v2/service/ssm v1.7.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.6.2
	github.com/aws/smithy-go v1.7.0
	github.com/blang/semver/v4 v4.0.0
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/fatih/color v1.12.0
	github.com/golang/mock v1.6.0
	github.com/jeremywohl/flatten v1.0.1
	github.com/rs/zerolog v1.22.0
	github.com/rsc/wes_client v0.0.0-00010101000000-000000000000
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/net v0.1.0
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.0
)

require (
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.2.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.2.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.3.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 // indirect
	golang.org/x/sys v0.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/rsc/wes_client => ./../wes_api/client
