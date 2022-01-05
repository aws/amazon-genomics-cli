import { ToilBatchPolicy } from "./policies/toil-batch-policy";
import { ToilJobRole, ToilJobRoleProps } from "./toil-job-role";
import { Arn, Aws, Stack } from "aws-cdk-lib";
import { Construct } from "constructs";
import { PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

interface ToilEngineRoleProps extends ToilJobRoleProps {
  // This is the queue to which we are authorizing jobs to be submitted by
  // something with this role.
  jobQueueArn: string;
  // And this other role can be assigned by this role
  jobRoleArn: string;
}

// This role grants access to Toil job stores, but also the access needed to
// launch jobs on AWS Batch that themselves have a ToilJobRole role assigned.
export class ToilEngineRole extends ToilJobRole {
  constructor(scope: Construct, id: string, props: ToilEngineRoleProps) {
    const toilJobArn = Arn.format(
      {
        account: Aws.ACCOUNT_ID,
        region: Aws.REGION,
        partition: Aws.PARTITION,
        resource: "job-definition/*",
        service: "batch",
      },
      scope as Stack
    );
    super(scope, id, props, {
      ToilEngineBatchPolicy: new ToilBatchPolicy({
        ...props,
        toilJobArn: toilJobArn,
      }),
      // TODO: Can we restrict this to allow passing the role only to jobs?
      ToilIamPassJobRole: new PolicyDocument({
        assignSids: true,
        statements: [
          new PolicyStatement({
            effect: Effect.ALLOW,
            actions: ["iam:PassRole"],
            resources: [props.jobRoleArn],
          }),
        ],
      }),
      ToilEcsDescribeInstances: new PolicyDocument({
        assignSids: true,
        statements: [
          new PolicyStatement({
            effect: Effect.ALLOW,
            actions: ["ecs:DescribeContainerInstances", "s3:ListAllMyBuckets"],
            resources: ["*"],
          }),
        ],
      }),
    });
  }
}
