import { Size, Stack, StackProps } from "aws-cdk-lib";
import { IVpc, SubnetSelection, Vpc } from "aws-cdk-lib/aws-ec2";
import { Construct } from "constructs";
import { getCommonParameter, getCommonParameterList, subnetSelectionFromIds } from "../util";
import {
  ENGINE_CROMWELL,
  ENGINE_MINIWDL,
  ENGINE_NEXTFLOW,
  ENGINE_SNAKEMAKE,
  VPC_NUMBER_SUBNETS_PARAMETER_NAME,
  VPC_PARAMETER_NAME,
  VPC_SUBNETS_PARAMETER_NAME,
} from "../constants";
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
  private readonly iops: Size;
  private readonly subnets: SubnetSelection;

  constructor(scope: Construct, id: string, props: ContextStackProps) {
    super(scope, id, props);

    const vpcId = getCommonParameter(this, VPC_PARAMETER_NAME);
    this.vpc = Vpc.fromLookup(this, "Vpc", { vpcId });
    const subnetIds = getCommonParameterList(this, VPC_SUBNETS_PARAMETER_NAME, VPC_NUMBER_SUBNETS_PARAMETER_NAME);
    this.subnets = subnetSelectionFromIds(this, subnetIds);

    const { contextParameters } = props;
    const { engineName } = contextParameters;
    const { filesystemType } = contextParameters;
    const { fsProvisionedThroughput } = contextParameters;
    this.iops = Size.mebibytes(fsProvisionedThroughput!);

    switch (engineName) {
      case ENGINE_CROMWELL:
        if (filesystemType != "S3") {
          throw Error(`'Cromwell' requires filesystem type 'S3'`);
        }
        if (contextParameters.usePublicSubnets) {
          throw Error(`'Cromwell cannot be securely deployed using public subnets'`);
        }
        this.renderCromwellStack(props);
        break;
      case ENGINE_NEXTFLOW:
        if (filesystemType != "S3") {
          throw Error(`'Nextflow' requires filesystem type 'S3'`);
        }
        this.renderNextflowStack(props);
        break;
      case ENGINE_MINIWDL:
        if (filesystemType != "EFS") {
          throw Error(`'MiniWDL' requires filesystem type 'EFS'`);
        }
        if (contextParameters.usePublicSubnets) {
          throw Error(`'miniwdl is not currently supported using public subnets, please file a github issue detailing your use case'`);
        }
        this.renderMiniwdlStack(props);
        break;
      case ENGINE_SNAKEMAKE:
        if (filesystemType != "EFS") {
          throw Error(`'Snakemake' requires filesystem type 'EFS'`);
        }
        if (contextParameters.usePublicSubnets) {
          throw Error(`'snakemake is not currently supported using public subnets, please file a github issue detailing your use case'`);
        }
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
      subnets: this.subnets,
      iops: this.iops,
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
      subnets: this.subnets,
      iops: this.iops,
      contextParameters: props.contextParameters,
      policyOptions: {
        managedPolicies: [],
      },
    };
  }
}
