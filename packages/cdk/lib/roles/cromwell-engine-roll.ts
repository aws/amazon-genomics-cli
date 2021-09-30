import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { S3ListAllBucketsPolicy } from "./policies/s3-list-all-buckets-policy";

interface CromwellEngineS3PolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
}

export class CromwellEngineRole extends iam.Role {
  readonly props: CromwellEngineS3PolicyProps;

  constructor(scope: cdk.Construct, id: string, props: CromwellEngineS3PolicyProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        EngineS3Policy: new S3ListAllBucketsPolicy(),
      },
      ...props.policies,
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
