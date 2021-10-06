import * as iam from "monocdk/aws-iam";

export interface NextflowSubmitJobBatchPolicyProps {
  account: string;
  region: string;
  submitJobPolicyArns: string[];
}

export class NextflowSubmitJobBatchPolicy extends iam.PolicyDocument {
  constructor(props: NextflowSubmitJobBatchPolicyProps) {
    const nextflowJobArn = `arn:aws:batch:${props.region}:${props.account}:job-definition/nf-*:*`;

    super({
      assignSids: true,
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["batch:SubmitJob"],
          resources: [...props.submitJobPolicyArns, nextflowJobArn],
        }),
      ],
    });
  }
}
