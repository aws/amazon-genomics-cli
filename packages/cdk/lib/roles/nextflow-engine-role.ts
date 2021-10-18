import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { NextflowEngineBatchPolicy } from "./policies/nextflow-engine-batch-policy";
import { NextflowAdapterBatchPolicy } from "./policies/nextflow-adapter-batch-policy";
import { Arn, Stack } from "monocdk";

interface NextflowEngineRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
  batchJobPolicyArns: string[];
}

export class NextflowEngineRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowEngineRoleProps) {
    const nextflowJobDefinitionArn = Arn.format(
      {
        resource: "job-definition/*",
        service: "batch",
      },
      scope as Stack
    );
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowBatchPolicy: new NextflowEngineBatchPolicy({
          nextflowJobArn: nextflowJobDefinitionArn,
        }),
        NextflowBatchSubmitPolicy: new NextflowAdapterBatchPolicy({
          batchJobPolicyArns: [...props.batchJobPolicyArns, nextflowJobDefinitionArn],
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

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
