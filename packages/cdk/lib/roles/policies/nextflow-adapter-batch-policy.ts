import { PolicyDocument, PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";

export interface NextflowSubmitJobBatchPolicyProps {
  batchJobPolicyArns: string[];
}

export class NextflowAdapterBatchPolicy extends PolicyDocument {
  constructor(props: NextflowSubmitJobBatchPolicyProps) {
    super({
      assignSids: true,
      statements: [
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: ["batch:SubmitJob", "batch:TerminateJob"],
          resources: props.batchJobPolicyArns,
        }),
      ],
    });
  }
}
