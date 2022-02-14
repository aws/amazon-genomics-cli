import { Policy, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";
import { batchArn } from "../../util";
import { BatchPolicies } from "./batch-policies";

export class HeadJobBatchPolicy extends Policy {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      statements: [
        BatchPolicies.listAndDescribe,
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:RegisterJobDefinition", "batch:DeregisterJobDefinition"],
          resources: [batchArn(scope, "job-definition")],
        }),
      ],
    });
  }
}
