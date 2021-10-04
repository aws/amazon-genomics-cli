import * as iam from "monocdk/aws-iam";

export interface CromwellBatchPolicyProps {
  account: string;
  region: string;
  jobQueueArn: string;
}

export class CromwellBatchPolicy extends iam.PolicyDocument {
  constructor(props: CromwellBatchPolicyProps) {
    const cromwellJobArn = `arn:aws:batch:${props.region}:${props.account}:job-definition/cromwell_*`;
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
          resources: [cromwellJobArn],
        }),
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:SubmitJob"],
          resources: [`${cromwellJobArn}:*`, props.jobQueueArn],
        }),
      ],
    });
  }
}
