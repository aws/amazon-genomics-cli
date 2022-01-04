import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../common/BucketOperations";
import { HeadJobBatchPolicy } from "./policies/head-job-batch-policy";
import { NextflowAdapterBatchPolicy } from "./policies/nextflow-adapter-batch-policy";
import { batchArn } from "../util";
import { Role, ServicePrincipal, PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";

interface NextflowEngineRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
  batchJobPolicyArns: string[];
}

export class NextflowEngineRole extends Role {
  constructor(scope: Construct, id: string, props: NextflowEngineRoleProps) {
    super(scope, id, {
      assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowBatchSubmitPolicy: new NextflowAdapterBatchPolicy({
          batchJobPolicyArns: [...props.batchJobPolicyArns, batchArn(scope, "job-definition")],
        }),
        NextflowLogsPolicy: new PolicyDocument({
          assignSids: true,
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
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
