import { Role, ServicePrincipal, ManagedPolicy } from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";
import { BucketOperations } from "../common/BucketOperations";

export interface CromwellAdapterRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
}

export class CromwellAdapterRole extends Role {
  constructor(scope: Construct, id: string, props: CromwellAdapterRoleProps) {
    super(scope, id, {
      assumedBy: new ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaVPCAccessExecutionRole")],
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
