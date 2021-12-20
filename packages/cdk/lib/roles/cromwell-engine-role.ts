import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../common/BucketOperations";
import { CromwellBatchPolicy } from "./policies/cromwell-batch-policy";
import { Arn, Aws, Stack } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Role, ServicePrincipal, PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

interface CromwellEngineRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
  jobQueueArn: string;
}

export class CromwellEngineRole extends Role {
  constructor(scope: Construct, id: string, props: CromwellEngineRoleProps) {
    const cromwellJobArn = Arn.format(
      {
        account: Aws.ACCOUNT_ID,
        region: Aws.REGION,
        partition: Aws.PARTITION,
        resource: "job-definition/*",
        service: "batch",
      },
      scope as Stack
    );
    super(scope, id, {
      assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        CromwellEngineBatchPolicy: new CromwellBatchPolicy({
          ...props,
          cromwellJobArn: cromwellJobArn,
        }),
        CromwellEcsDescribeInstances: new PolicyDocument({
          assignSids: true,
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["ecs:DescribeContainerInstances", "s3:ListAllMyBuckets"],
              resources: ["*"],
            }),
          ],
        }),
      },
      ...props.policies,
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
