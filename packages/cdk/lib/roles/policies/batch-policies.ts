import { PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

export class BatchPolicies {
  static readonly listAndDescribe = new PolicyStatement({
    effect: Effect.ALLOW,
    actions: ["batch:Describe*", "batch:ListJobs"],
    resources: ["*"],
  });
}
