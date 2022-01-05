import { PolicyOptions } from "../types/engine-options";
import { BucketOperations } from "../common/BucketOperations";
import { Construct } from "constructs";
import { Role, ServicePrincipal, PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

export interface ToilJobRoleProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  policies: PolicyOptions;
}

// This role grants access to everything a Toil job needs to talk to the AWS
// job store and/or additional user data in S3.
export class ToilJobRole extends Role {
  constructor(scope: Construct, id: string, props: ToilJobRoleProps, additionalInlinePolicies?: { [key: string]: PolicyDocument }) {
    super(scope, id, {
      assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        // TODO: Remove this when Toil no longer uses its own SimpleDB domains
        ToilSimpleDBFullAccess: new PolicyDocument({
          assignSids: true,
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["sdb:*"],
              resources: ["*"],
            }),
          ],
        }),
        // TODO: Remove this when Toil is taught to use AGC buckets to store
        // its workflow state and doesn't need to make and destroy its own.
        ToilS3FullAccess: new PolicyDocument({
          assignSids: true,
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["s3:*"],
              resources: ["*"],
            }),
          ],
        }),
        ...additionalInlinePolicies,
      },
      ...props.policies,
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
