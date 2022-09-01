module github.com/aws/amazon-genomics-cli

go 1.16

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
	github.com/kr/pretty v0.2.1 // indirect
	github.com/rs/zerolog v1.22.0
	github.com/rsc/wes_client v0.0.0-00010101000000-000000000000
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/rsc/wes_client => ./../wes_api/client
