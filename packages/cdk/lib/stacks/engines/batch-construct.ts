import { IMachineImage, IVpc, SubnetSelection } from "aws-cdk-lib/aws-ec2";
import { Stack } from "aws-cdk-lib";
import { Batch } from "../../constructs";
import { ContextAppParameters } from "../../env";
import { BucketOperations } from "../../common/BucketOperations";
import { IRole } from "aws-cdk-lib/aws-iam";
import { ComputeResourceType } from "@aws-cdk/aws-batch-alpha";
import { Construct } from "constructs";
import { LaunchTemplateData } from "../../constructs/launch-template-data";

export interface BatchConstructProps {
  /**
   * VPC to run resources in.
   */
  readonly vpc: IVpc;
  /**
   * The subnets of a vpc to use for batch compute environments.
   */
  readonly subnets: SubnetSelection;
  /**
   * Parameters determined by the context.
   */
  readonly contextParameters: ContextAppParameters;
  /**
   * Request Spot capacity to be created
   */
  readonly createSpotBatch: boolean;
  /**
   * Request On-Demand capacity to be created
   */
  readonly createOnDemandBatch: boolean;
  /**
   * AMI used for compute
   */
  readonly computeEnvImage?: IMachineImage;
}

export class BatchConstruct extends Construct {
  public readonly batchSpot: Batch;
  public readonly batchOnDemand: Batch;

  constructor(scope: Construct, id: string, props: BatchConstructProps) {
    super(scope, id);

    const { vpc, contextParameters, createSpotBatch, createOnDemandBatch, subnets, computeEnvImage } = props;
    const { artifactBucketName, outputBucketName, readBucketArns = [], readWriteBucketArns = [] } = contextParameters;
    if (createSpotBatch) {
      this.batchSpot = this.renderBatch("TaskBatchSpot", vpc, subnets, contextParameters, ComputeResourceType.SPOT, computeEnvImage);
    }
    if (createOnDemandBatch) {
      this.batchOnDemand = this.renderBatch("TaskBatch", vpc, subnets, contextParameters, ComputeResourceType.ON_DEMAND, computeEnvImage);
    }

    const artifactBucket = BucketOperations.importBucket(this, "ArtifactBucket", artifactBucketName);
    const outputBucket = BucketOperations.importBucket(this, "OutputBucket", outputBucketName);

    readBucketArns.push(artifactBucket.bucketArn);
    readWriteBucketArns.push(outputBucket.bucketArn);

    const batchRoles = this.getBatchRoles();
    for (const role of batchRoles) {
      BucketOperations.grantBucketAccess(this, role, readBucketArns, true);
      BucketOperations.grantBucketAccess(this, role, readWriteBucketArns, false);
    }
  }

  private renderBatch(
    id: string,
    vpc: IVpc,
    subnets: SubnetSelection,
    appParams: ContextAppParameters,
    computeType?: ComputeResourceType,
    computeEnvImage?: IMachineImage
  ): Batch {
    return new Batch(this, id, {
      vpc,
      computeType,
      subnets,
      computeEnvImage,
      instanceTypes: appParams.instanceTypes,
      maxVCpus: appParams.maxVCpus,
      launchTemplateData: LaunchTemplateData.renderLaunchTemplateData(appParams.engineName),
      awsPolicyNames: ["AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"],
      resourceTags: Stack.of(this).tags.tagValues(),
      usePublicSubnets: appParams.usePublicSubnets,
    });
  }

  private getBatchRoles(): IRole[] {
    const roles: IRole[] = [];
    if (this.batchOnDemand) {
      roles.push(this.batchOnDemand.role);
    }
    if (this.batchSpot) {
      roles.push(this.batchSpot.role);
    }
    return roles;
  }
}
