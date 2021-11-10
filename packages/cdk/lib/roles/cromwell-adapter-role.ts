import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { BucketOperations } from "../../common/BucketOperations";

export interface CromwellAdapterRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class CromwellAdapterRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: CromwellAdapterRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [iam.ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaVPCAccessExecutionRole")],
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
