import * as iam from "monocdk/aws-iam";

export interface NextflowBatchPolicyProps {
  account: string;
  region: string;
}

export class NextflowBatchPolicy extends iam.PolicyDocument {
  constructor(props: NextflowBatchPolicyProps) {
    const nextflowJobArn = `arn:aws:batch:${props.region}:${props.account}:job-definition/nf-ubuntu:*`;

    super({
      assignSids: true,
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:DescribeJobDefinitions", "batch:ListJobs", "batch:DescribeJobs", "batch:DescribeJobQueues", "batch:DescribeComputeEnvironments"],
          resources: ["*"],
        }),
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:RegisterJobDefinition"],
          resources: [nextflowJobArn],
        }),
      ],
    });
  }
}
