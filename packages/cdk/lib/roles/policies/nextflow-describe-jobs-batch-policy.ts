import { PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

export class NextflowDescribeJobsBatchPolicy extends PolicyDocument {
  constructor() {
    super({
      assignSids: true,
      statements: [
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:DescribeJobs"],
          resources: ["*"],
        }),
      ],
    });
  }
}
