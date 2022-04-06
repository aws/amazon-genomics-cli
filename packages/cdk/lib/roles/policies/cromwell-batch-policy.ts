import { Effect, PolicyDocument, PolicyStatement } from "aws-cdk-lib/aws-iam";

export interface CromwellBatchPolicyProps {
  jobQueueArn: string;
  cromwellJobArn: string;
}

export class CromwellBatchPolicy extends PolicyDocument {
  constructor(props: CromwellBatchPolicyProps) {
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
          resources: [props.cromwellJobArn],
        }),
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:SubmitJob"],
          resources: [props.cromwellJobArn, props.jobQueueArn],
        }),
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:TerminateJob", "batch:CancelJob"],
          // can only be restricted to "job*" resources, but we cannot know the IDs of jobs that Cromwell has started here.
          // https://docs.aws.amazon.com/service-authorization/latest/reference/list_awsbatch.html
          resources: ["*"],
        }),
      ],
    });
  }
}
