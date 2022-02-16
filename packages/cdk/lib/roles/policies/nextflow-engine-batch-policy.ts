import { PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

export interface NextflowBatchPolicyProps {
  nextflowJobArn: string;
}

export class NextflowEngineBatchPolicy extends PolicyDocument {
  constructor(props: NextflowBatchPolicyProps) {
    super({
      assignSids: true,
      statements: [
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:DescribeJobDefinitions", "batch:ListJobs", "batch:DescribeJobs", "batch:DescribeJobQueues", "batch:DescribeComputeEnvironments"],
          resources: ["*"],
        }),
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:RegisterJobDefinition"],
          resources: [props.nextflowJobArn],
        }),
      ],
    });
  }
}
