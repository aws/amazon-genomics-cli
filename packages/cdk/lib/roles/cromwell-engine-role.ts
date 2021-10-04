import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { S3ListAllBucketsPolicy } from "./policies/s3-list-all-buckets-policy";
import { CromwellBatchPolicy, CromwellBatchPolicyProps } from "./policies/cromwell-batch-policy";

interface CromwellEngineRoleProps extends CromwellBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
}

export class CromwellEngineRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: CromwellEngineRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        EngineS3Policy: new S3ListAllBucketsPolicy(),
        EngineBatchPolicy: new CromwellBatchPolicy(props),
      },
      ...props.policies,
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
