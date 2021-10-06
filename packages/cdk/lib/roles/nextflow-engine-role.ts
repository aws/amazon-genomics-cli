import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { NextflowBatchPolicy, NextflowBatchPolicyProps } from "./policies/nextflow-batch-policy";
import { NextflowSubmitJobBatchPolicy, NextflowSubmitJobBatchPolicyProps } from "./policies/nextflow-submit-job-batch-policy";

interface NextflowEngineRoleProps extends NextflowBatchPolicyProps, NextflowSubmitJobBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
}

export class NextflowEngineRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowEngineRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowBatchPolicy: new NextflowBatchPolicy(props),
        NextflowBatchSubmitPolicy: new NextflowSubmitJobBatchPolicy(props),
        NextflowLogsPolicy: new iam.PolicyDocument({
          assignSids: true,
          statements: [
            new iam.PolicyStatement({
              effect: iam.Effect.ALLOW,
              actions: ["logs:GetQueryResults", "logs:StopQuery"],
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
