// Exclude CDK directory from Go module parsing
module ignore

go 1.19

require (
	github.com/aws/aws-cdk-go/awscdk/v2 v2.17.0
	github.com/aws/constructs-go/constructs/v10 v10.0.92
	github.com/aws/jsii-runtime-go v1.55.1
)

require (
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/aws/aws-cdk-go/awscdk v1.149.0-devpreview // indirect
	github.com/aws/constructs-go/constructs/v3 v3.3.246 // indirect
)
