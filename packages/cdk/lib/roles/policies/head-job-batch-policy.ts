import * as iam from "monocdk/aws-iam";
import { Construct } from "monocdk";
import { batchArn } from "../../util";
import { BatchPolicies } from "./batch-policies";

export class HeadJobBatchPolicy extends iam.Policy {
  constructor(scope: Construct, id: string) {
    super(scope, id, {
      statements: [
        BatchPolicies.listAndDescribe,
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:RegisterJobDefinition", "batch:DeregisterJobDefinition"],
          resources: [batchArn(scope, "job-definition")],
        }),
      ],
    });
  }
}
