import { NextflowAdapterBatchPolicy, NextflowSubmitJobBatchPolicyProps } from "./policies/nextflow-adapter-batch-policy";
import { batchArn } from "../util";
import { BucketOperations } from "../common/BucketOperations";
import { Role, ServicePrincipal, PolicyDocument, PolicyStatement, Effect, ManagedPolicy } from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";

export interface NextflowAdapterRoleProps extends NextflowSubmitJobBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class NextflowAdapterRole extends Role {
  constructor(scope: Construct, id: string, props: NextflowAdapterRoleProps) {
    const nextflowJobDefinitionArn = batchArn(scope, "job-definition");
    const nextflowJobArn = batchArn(scope, "job");

    super(scope, id, {
      assumedBy: new ServicePrincipal("lambda.amazonaws.com"),
      inlinePolicies: {
        NextflowDescribeJobsPolicy: new PolicyDocument({
          assignSids: true,
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["batch:DescribeJobs", "batch:ListJobs", "logs:GetQueryResults"],
              resources: ["*"],
            }),
          ],
        }),
        NextflowSubmitJobsPolicy: new NextflowAdapterBatchPolicy({
          batchJobPolicyArns: [...props.batchJobPolicyArns, nextflowJobDefinitionArn, nextflowJobArn],
        }),
      },
      managedPolicies: [ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaVPCAccessExecutionRole")],
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
