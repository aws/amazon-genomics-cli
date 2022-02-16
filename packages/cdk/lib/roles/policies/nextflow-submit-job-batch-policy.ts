import { PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

export interface NextflowSubmitJobBatchPolicyProps {
  headJobDefinitionArn: string;
  jobQueueArn: string;
}

export class NextflowSubmitJobBatchPolicy extends PolicyDocument {
  constructor(props: NextflowSubmitJobBatchPolicyProps) {
    super({
      assignSids: true,
      statements: [
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:SubmitJob"],
          resources: [props.headJobDefinitionArn, props.jobQueueArn],
        }),
      ],
    });
  }
}
