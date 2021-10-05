import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { NextflowDescribeJobsBatchPolicy } from "./policies/nextflow-describe-jobs-batch-policy";
import { NextflowSubmitJobBatchPolicy, NextflowSubmitJobBatchPolicyProps } from "./policies/nextflow-submit-job-batch-policy";
import { BucketOperations } from "../../common/BucketOperations";

export interface NextflowAdapterRoleProps extends NextflowSubmitJobBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class NextflowAdapterRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowAdapterRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowDescribeJobsPolicy: new NextflowDescribeJobsBatchPolicy(),
        NextflowSubmitJobsPolicy: new NextflowSubmitJobBatchPolicy(props),
      },
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
