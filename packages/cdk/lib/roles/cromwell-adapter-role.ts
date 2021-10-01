import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { NextflowDescribeJobsBatchPolicy } from "./policies/nextflow-describe-jobs-batch-policy";
import { NextflowSubmitJobBatchPolicy } from "./policies/nextflow-submit-job-batch-policy";
import { BucketOperations } from "../../common/BucketOperations";

export interface CromwellAdapterRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class CromwellAdapterRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: CromwellAdapterRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
