import * as iam from "monocdk/aws-iam";

export interface NextflowSubmitJobBatchPolicyProps {
  headJobDefinitionArn: string;
  jobQueueArn: string;
}

export class NextflowSubmitJobBatchPolicy extends iam.PolicyDocument {
  constructor(props: NextflowSubmitJobBatchPolicyProps) {
    super({
      assignSids: true,
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:SubmitJob"],
          resources: [props.headJobDefinitionArn, props.jobQueueArn],
        }),
      ],
    });
  }
}
