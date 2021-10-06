import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { CromwellBatchPolicy, CromwellBatchPolicyProps } from "./policies/cromwell-batch-policy";

interface CromwellEngineRoleProps extends CromwellBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
}

export class CromwellEngineRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: CromwellEngineRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        CromwellEngineBatchPolicy: new CromwellBatchPolicy(props),
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
