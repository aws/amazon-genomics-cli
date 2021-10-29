import * as iam from "monocdk/aws-iam";

export class BatchPolicies {
  static readonly listAndDescribe = new iam.PolicyStatement({
    effect: iam.Effect.ALLOW,
    actions: ["batch:Describe*", "batch:ListJobs"],
    resources: ["*"],
  });
}
