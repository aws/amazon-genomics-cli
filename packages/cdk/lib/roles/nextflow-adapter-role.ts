import * as cdk from "monocdk";
import * as iam from "monocdk/aws-iam";
import { NextflowAdapterBatchPolicy, NextflowSubmitJobBatchPolicyProps } from "./policies/nextflow-adapter-batch-policy";
import { BucketOperations } from "../../common/BucketOperations";
import { Arn, ArnComponents } from "monocdk";
import { Stack } from "monocdk";

export interface NextflowAdapterRoleProps extends NextflowSubmitJobBatchPolicyProps {
  readOnlyBucketArns: string[];
  readWriteBucketArns: string[];
  components: ArnComponents;
}

export class NextflowAdapterRole extends iam.Role {
  constructor(scope: cdk.Construct, id: string, props: NextflowAdapterRoleProps) {
    const nextflowJobDefinitionArn = Arn.format(props.components, scope as Stack);
    const nextflowJobArn = Arn.format(
      {
        ...props.components,
        resource: "job/*",
      },
      scope as Stack
    );
    super(scope, id, {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      inlinePolicies: {
        NextflowDescribeJobsPolicy: new iam.PolicyDocument({
          assignSids: true,
          statements: [
            new iam.PolicyStatement({
              effect: iam.Effect.ALLOW,
              actions: ["batch:DescribeJobs", "logs:GetQueryResults"],
              resources: ["*"],
            }),
          ],
        }),
        NextflowSubmitJobsPolicy: new NextflowAdapterBatchPolicy({
          batchJobPolicyArns: [...props.batchJobPolicyArns, nextflowJobDefinitionArn, nextflowJobArn],
        }),
      },
    });

    BucketOperations.grantBucketAccess(this, this, props.readOnlyBucketArns, true);
    BucketOperations.grantBucketAccess(this, this, props.readWriteBucketArns);
  }
}
