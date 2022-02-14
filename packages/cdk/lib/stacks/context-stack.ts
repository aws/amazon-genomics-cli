import { Stack, StackProps } from "aws-cdk-lib";
import { IVpc, Vpc } from "aws-cdk-lib/aws-ec2";
import { Construct } from "constructs";
import { getCommonParameter } from "../util";
import { ENGINE_CROMWELL, ENGINE_MINIWDL, ENGINE_NEXTFLOW, ENGINE_SNAKEMAKE, VPC_PARAMETER_NAME } from "../constants";
import { ContextAppParameters } from "../env";
import { BatchConstruct, BatchConstructProps } from "./engines/batch-construct";
import { CromwellEngineConstruct } from "./engines/cromwell-engine-construct";
import { NextflowEngineConstruct } from "./engines/nextflow-engine-construct";
import { MiniwdlEngineConstruct } from "./engines/miniwdl-engine-construct";
import { SnakemakeEngineConstruct } from "./engines/snakemake-engine-construct";

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
      case ENGINE_CROMWELL:
        this.renderCromwellStack(props);
        break;
      case ENGINE_NEXTFLOW:
        this.renderNextflowStack(props);
        break;
      case ENGINE_MINIWDL:
        this.renderMiniwdlStack(props);
        break;
      case ENGINE_SNAKEMAKE:
        this.renderSnakemakeStack(props);
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
    new CromwellEngineConstruct(this, ENGINE_CROMWELL, {
      jobQueue,
      ...commonEngineProps,
    }).outputToParent();
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
    new NextflowEngineConstruct(this, ENGINE_NEXTFLOW, {
      ...commonEngineProps,
      jobQueue,
      headQueue,
    }).outputToParent();
  }

  private renderMiniwdlStack(props: ContextStackProps) {
    const commonEngineProps = this.getCommonEngineProps(props);
    new MiniwdlEngineConstruct(this, ENGINE_MINIWDL, {
      ...commonEngineProps,
    }).outputToParent();
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

  private renderSnakemakeStack(props: ContextStackProps) {
    const commonEngineProps = this.getCommonEngineProps(props);
    new SnakemakeEngineConstruct(this, ENGINE_SNAKEMAKE, {
      ...commonEngineProps,
    }).outputToParent();
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
      parent: this,
    };
  }

  private renderBatchStack(props: BatchConstructProps) {
    return new BatchConstruct(this, "Batch", props);
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
