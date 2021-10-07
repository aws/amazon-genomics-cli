import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { NextflowAdapterBatchPolicy, NextflowSubmitJobBatchPolicyProps } from "./policies/nextflow-adapter-batch-policy";
import { BucketOperations } from "../../common/BucketOperations";
import { Arn, ArnComponents } from "monocdk";
import { Stack } from "monocdk";

export interface NextflowAdapterRoleProps extends NextflowSubmitJobBatchPolicyProps {
  account: string;
  region: string;
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class NextflowAdapterRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowAdapterRoleProps) {
    const nextflowJobDefinitionArn = Arn.format(
      {
        account: props.account,
        region: props.region,
        resource: "job-definition/*",
        service: "batch",
      },
      scope as Stack
    );
    const nextflowJobArn = Arn.format(
      {
        account: props.account,
        region: props.region,
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
              actions: ["batch:DescribeJobs", "logs:GetQueryResults"],
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
