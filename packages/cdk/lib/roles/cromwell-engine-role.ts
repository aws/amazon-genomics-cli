import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { CromwellBatchPolicy } from "./policies/cromwell-batch-policy";
import { Arn, Stack } from "monocdk";

interface CromwellEngineRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
  jobQueueArn: string;
}

export class CromwellEngineRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: CromwellEngineRoleProps) {
    const cromwellJobArn = Arn.format(
      {
        resource: "job-definition/*",
        service: "batch",
      },
      scope as Stack
    );
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        CromwellEngineBatchPolicy: new CromwellBatchPolicy({
          ...props,
          cromwellJobArn: cromwellJobArn,
        }),
        CromwellEcsDescribeInstances: new iam.PolicyDocument({
          assignSids: true,
          statements: [
            new iam.PolicyStatement({
              effect: iam.Effect.ALLOW,
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
