import * as iam from "monocdk/aws-iam";

export interface NextflowBatchPolicyProps {
  nextflowJobArn: string;
}

export class NextflowEngineBatchPolicy extends iam.PolicyDocument {
  constructor(props: NextflowBatchPolicyProps) {
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
          resources: [props.nextflowJobArn],
        }),
      ],
    });
  }
}
