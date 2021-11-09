import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { NextflowAdapterBatchPolicy, NextflowSubmitJobBatchPolicyProps } from "./policies/nextflow-adapter-batch-policy";
import { batchArn } from "../util";
import { BucketOperations } from "../../common/BucketOperations";

export interface NextflowAdapterRoleProps extends NextflowSubmitJobBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class NextflowAdapterRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowAdapterRoleProps) {
    const nextflowJobDefinitionArn = batchArn(scope, "job-definition");
    const nextflowJobArn = batchArn(scope, "job");

    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
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
      managedPolicies: [iam.ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaVPCAccessExecutionRole")],
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
