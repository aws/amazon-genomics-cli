import * as iam from "monocdk/aws-iam";

export interface NextflowSubmitJobBatchPolicyProps {
  batchJobPolicyArns: string[];
}

export class NextflowAdapterBatchPolicy extends iam.PolicyDocument {
  constructor(props: NextflowSubmitJobBatchPolicyProps) {
    super({
      assignSids: true,
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:SubmitJob", "batch:TerminateJob"],
          resources: props.batchJobPolicyArns,
        }),
      ],
    });
  }
}
