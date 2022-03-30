import { PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

export interface ToilBatchPolicyProps {
  jobQueueArn: string;
  toilJobDefinitionArnPattern: string;
}

export class ToilBatchPolicy extends PolicyDocument {
  constructor(props: ToilBatchPolicyProps) {
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
          actions: ["batch:RegisterJobDefinition", "batch:DeregisterJobDefinition"],
          resources: [props.toilJobDefinitionArnPattern],
        }),
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:SubmitJob"],
          resources: [props.toilJobDefinitionArnPattern, props.jobQueueArn],
        }),
      ],
    });
  }
}
