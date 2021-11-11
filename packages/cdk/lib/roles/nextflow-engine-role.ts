import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { HeadJobBatchPolicy } from "./policies/head-job-batch-policy";
import { NextflowAdapterBatchPolicy } from "./policies/nextflow-adapter-batch-policy";
import { batchArn } from "../util";

interface NextflowEngineRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
  batchJobPolicyArns: string[];
}

export class NextflowEngineRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowEngineRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowBatchSubmitPolicy: new NextflowAdapterBatchPolicy({
          batchJobPolicyArns: [...props.batchJobPolicyArns, batchArn(scope, "job-definition")],
        }),
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

    this.attachInlinePolicy(new HeadJobBatchPolicy(this, "NextflowHeadJobBatchPolicy"));

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
