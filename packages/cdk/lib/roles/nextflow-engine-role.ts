import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../../common/BucketOperations";
import { NextflowLogsPolicy } from "./policies/nextflow-logs-policy";

interface NextflowEngineRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
}

export class NextflowEngineRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowEngineRoleProps) {
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowLogPolicy: new NextflowLogsPolicy(),
      },
      ...props.policies,
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
