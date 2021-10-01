import * as iam from "monocdk/aws-iam";
import { NextflowAdapterRoleProps } from "../nextflow-adapter-role";

export class NextflowSubmitJobBatchPolicy extends iam.PolicyDocument {
  constructor(props: NextflowAdapterRoleProps) {
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
