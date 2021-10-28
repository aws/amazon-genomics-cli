import * as iam from "monocdk/aws-iam";

export class NextflowDescribeJobsBatchPolicy extends iam.PolicyDocument {
  constructor() {
    super({
      assignSids: true,
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:DescribeJobs"],
          resources: ["*"],
        }),
      ],
    });
  }
}
