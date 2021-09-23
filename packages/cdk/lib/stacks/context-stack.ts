import { Stack, StackProps } from "monocdk";
import { Vpc } from "monocdk/aws-ec2";
import { Construct } from "constructs";
import { getCommonParameter } from "../util";
import { VPC_PARAMETER_NAME } from "../constants";
import { ContextAppParameters } from "../env";
import { BatchStack } from "./nested/batch-stack";
import { CromwellEngineStack, CromwellEngineStackProps } from "./nested/cromwell-engine-stack";
import { NextflowEngineStack, NextflowEngineStackProps } from "./nested/nextflow-engine-stack";
import { ManagedPolicy } from "monocdk/aws-iam";

export interface ContextStackProps extends StackProps {
  readonly contextParameters: ContextAppParameters;
}

export class ContextStack extends Stack {
  constructor(scope: Construct, id: string, props: ContextStackProps) {
    super(scope, id, props);

    const vpcId = getCommonParameter(this, VPC_PARAMETER_NAME);
    const vpc = Vpc.fromLookup(this, "Vpc", { vpcId });

    const batchStack = new BatchStack(this, "Batch", { vpc, contextParameters: props.contextParameters });

    const commonProps: CromwellEngineStackProps | NextflowEngineStackProps = {
      vpc,
      contextParameters: props.contextParameters,
      jobQueue: batchStack.batchWorkers.jobQueue,
      policyOptions: {
        managedPolicies: [
          // TODO: Can these be scoped down?
          ManagedPolicy.fromAwsManagedPolicyName("AmazonEC2ContainerRegistryReadOnly"),
          ManagedPolicy.fromAwsManagedPolicyName("AmazonECS_FullAccess"),
          ManagedPolicy.fromAwsManagedPolicyName("AmazonElasticFileSystemFullAccess"),
          ManagedPolicy.fromAwsManagedPolicyName("AmazonS3ReadOnlyAccess"),
          ManagedPolicy.fromAwsManagedPolicyName("AWSBatchFullAccess"),
        ],
      },
    };
    const engineName = props.contextParameters.engineName;
    switch (engineName) {
      case "cromwell":
        new CromwellEngineStack(this, engineName, {
          ...commonProps,
        }).outputToParent(this);
        break;
      case "nextflow":
        new NextflowEngineStack(this, engineName, {
          ...commonProps,
          headQueue: batchStack.batchHead.jobQueue,
        }).outputToParent(this);
        break;
      default:
        throw Error(`Engine '${engineName}' is not supported`);
    }
  }
}
