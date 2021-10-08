import { Stack, StackProps } from "monocdk";
import { IVpc, Vpc } from "monocdk/aws-ec2";
import { Construct } from "constructs";
import { getCommonParameter } from "../util";
import { VPC_PARAMETER_NAME } from "../constants";
import { ContextAppParameters } from "../env";
import { BatchStack, BatchStackProps } from "./nested/batch-stack";
import { CromwellEngineStack } from "./nested/cromwell-engine-stack";
import { NextflowEngineStack } from "./nested/nextflow-engine-stack";
import { ManagedPolicy } from "monocdk/aws-iam";

export interface ContextStackProps extends StackProps {
  readonly contextParameters: ContextAppParameters;
}

export class ContextStack extends Stack {
  private readonly vpc: IVpc;

  constructor(scope: Construct, id: string, props: ContextStackProps) {
    super(scope, id, props);

    const vpcId = getCommonParameter(this, VPC_PARAMETER_NAME);
    this.vpc = Vpc.fromLookup(this, "Vpc", { vpcId });

    const { contextParameters } = props;
    const { engineName } = contextParameters;

    switch (engineName) {
      case "cromwell":
        this.renderCromwellStack(props);
        break;
      case "nextflow":
        this.renderNextflowStack(props);
        break;
      default:
        throw Error(`Engine '${engineName}' is not supported`);
    }
  }

  private renderCromwellStack(props: ContextStackProps) {
    const batchProps = this.getCromwellBatchProps(props);
    const batchStack = this.renderBatchStack(batchProps);

    let jobQueue;
    if (props.contextParameters.requestSpotInstances) {
      jobQueue = batchStack.batchSpot.jobQueue;
    } else {
      jobQueue = batchStack.batchOnDemand.jobQueue;
    }

    const commonEngineProps = this.getCommonEngineProps(props);
    new CromwellEngineStack(this, "cromwell", {
      jobQueue,
      ...commonEngineProps,
    }).outputToParent(this);
  }

  private renderNextflowStack(props: ContextStackProps) {
    const batchProps = this.getNextflowBatchProps(props);
    const batchStack = this.renderBatchStack(batchProps);

    let jobQueue, headQueue;
    if (props.contextParameters.requestSpotInstances) {
      jobQueue = batchStack.batchSpot.jobQueue;
      headQueue = batchStack.batchOnDemand.jobQueue;
    } else {
      headQueue = jobQueue = batchStack.batchOnDemand.jobQueue;
    }

    const commonEngineProps = this.getCommonEngineProps(props);
    new NextflowEngineStack(this, "nextflow", {
      ...commonEngineProps,
      jobQueue,
      headQueue,
    }).outputToParent(this);
  }

  private getCromwellBatchProps(props: ContextStackProps) {
    const commonBatchProps = this.getCommonBatchProps(props);
    const { requestSpotInstances } = props.contextParameters;

    return {
      ...commonBatchProps,
      createSpotBatch: requestSpotInstances,
      createOnDemandBatch: !requestSpotInstances,
    };
  }

  private getCommonBatchProps(props: ContextStackProps) {
    const { contextParameters } = props;
    return {
      vpc: this.vpc,
      contextParameters,
    };
  }

  private getNextflowBatchProps(props: ContextStackProps) {
    const commonBatchProps = this.getCommonBatchProps(props);
    const { requestSpotInstances } = props.contextParameters;
    return {
      ...commonBatchProps,
      createSpotBatch: requestSpotInstances,
      createOnDemandBatch: true,
    };
  }

  private renderBatchStack(props: BatchStackProps) {
    return new BatchStack(this, "Batch", props);
  }
  private getCommonEngineProps(props: ContextStackProps) {
    return {
      vpc: this.vpc,
      contextParameters: props.contextParameters,
      policyOptions: {
        managedPolicies: [],
      },
    };
  }
}
