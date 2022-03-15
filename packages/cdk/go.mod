// Exclude CDK directory from Go module parsing
module ignore

go 1.17

require (
	github.com/aws/aws-cdk-go/awscdk v1.148.0-devpreview
	github.com/aws/aws-cdk-go/awscdk/v2 v2.16.0
	github.com/aws/constructs-go/constructs/v10 v10.0.88
	github.com/aws/constructs-go/constructs/v3 v3.3.243
	github.com/aws/jsii-runtime-go v1.55.0
)

require github.com/Masterminds/semver/v3 v3.1.1 // indirect
