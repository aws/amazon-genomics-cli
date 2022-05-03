import { PolicyStatement, Effect } from "aws-cdk-lib/aws-iam";
import { CromwellBatchPolicy } from "./cromwell-batch-policy";

export interface ToilBatchPolicyProps {
  jobQueueArn: string;
  toilJobArnPattern: string;
}

export class ToilBatchPolicy extends CromwellBatchPolicy {
  constructor(props: ToilBatchPolicyProps) {
    // To avoid adding more policies allowing access to "*", we are based on
    // the Cromwell policy set. When the permissions for that get locked
    // down to the minimum required to use Batch, we will inherit those
    // improvements.
    super({
      jobQueueArn: props.jobQueueArn,
      cromwellJobArn: props.toilJobArnPattern,
    });

    // The only additional thing we need is to be able to deregister job
    // definitions, which Cromwell doesn't do.
    this.addStatements(
      new PolicyStatement({
        effect: Effect.ALLOW,
        actions: ["batch:DeregisterJobDefinition"],
        resources: [props.toilJobArnPattern],
      })
    );
  }
}
