import * as iam from "monocdk/aws-iam";

export interface CromwellBatchPolicyProps {
  jobQueueArn: string;
  cromwellJobArn: string;
}

export class CromwellBatchPolicy extends iam.PolicyDocument {
  constructor(props: CromwellBatchPolicyProps) {
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
          resources: [props.cromwellJobArn],
        }),
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:SubmitJob"],
          resources: [props.cromwellJobArn, props.jobQueueArn],
        }),
      ],
    });
  }
}
