package actionableerror

var AwsErrorMessageToSuggestedActionMap = map[string]string{
	"The security token included in the request is expire":                                             "Please refresh your aws credentials and try again",
	"api error AccessDeniedException: User: arn:aws:":                                                  "Please validate that you have sufficient permissions and try again",
	"so the toolkit stack must be deployed to the environment":                                         "Please bootstrap your account before retrying the current command",
	"an AWS region is required, but was not found":                                                     "Please either set the region to use with an AWS_REGION variable or run aws configure",
	"failed to retrieve credentials:":                                                                  "Please check that you have set some aws credentials and try again",
	"The security token included in the request is invalid":                                            "Please validate that the credentials you have set are valid",
	"unable to obtain bucket output name. The SSM parameter /agc/_common/bucket may be misconfigured.": "Please ensure your account is activated by running 'agc account activate'",
}
