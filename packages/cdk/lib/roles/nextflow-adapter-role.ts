import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { NextflowAdapterBatchPolicy, NextflowSubmitJobBatchPolicyProps } from "./policies/nextflow-adapter-batch-policy";
import { BucketOperations } from "../../common/BucketOperations";
import { Arn } from "monocdk";
import { Stack } from "monocdk";

export interface NextflowAdapterRoleProps extends NextflowSubmitJobBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class NextflowAdapterRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowAdapterRoleProps) {
    const nextflowJobDefinitionArn = Arn.format(
      {
        resource: "job-definition/*",
        service: "batch",
      },
      scope as Stack
    );
    const nextflowJobArn = Arn.format(
      {
        resource: "job/*",
        service: "batch",
      },
      scope as Stack
    );
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowDescribeJobsPolicy: new iam.PolicyDocument({
          assignSids: true,
          statements: [
            new iam.PolicyStatement({
              effect: iam.Effect.ALLOW,
              actions: ["batch:DescribeJobs", "batch:ListJobs", "logs:GetQueryResults"],
              resources: ["*"],
            }),
          ],
        }),
        NextflowSubmitJobsPolicy: new NextflowAdapterBatchPolicy({
          batchJobPolicyArns: [...props.batchJobPolicyArns, nextflowJobDefinitionArn, nextflowJobArn],
        }),
      },
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
