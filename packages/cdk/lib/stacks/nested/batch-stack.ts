import { NestedStack, NestedStackProps } from "monocdk";
import { InstanceType, IVpc } from "monocdk/aws-ec2";
import { Construct } from "constructs";
import { LAUNCH_TEMPLATE } from "../../constants";

import { Batch } from "../../constructs";
import { ContextAppParameters } from "../../env";
import { BucketOperations } from "../../../common/BucketOperations";
import { IRole } from "monocdk/aws-iam";
import { ComputeResourceType } from "monocdk/aws-batch";

export interface BatchStackProps extends NestedStackProps {
  /**
   * VPC to run resources in.
   */
  readonly vpc: IVpc;
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
}

export class BatchStack extends NestedStack {
  public readonly batchSpot: Batch;
  public readonly batchOnDemand: Batch;

  constructor(scope: Construct, id: string, props: BatchStackProps) {
    super(scope, id, props);

    const { vpc, contextParameters, createSpotBatch, createOnDemandBatch } = props;

    if (createSpotBatch) {
      this.batchSpot = this.renderBatch("TaskBatchSpot", vpc, contextParameters.instanceTypes, ComputeResourceType.SPOT);
    }
    if (createOnDemandBatch) {
      this.batchOnDemand = this.renderBatch("TaskBatch", vpc, contextParameters.instanceTypes, ComputeResourceType.ON_DEMAND);
    }

    const artifactBucket = BucketOperations.importBucket(this, "ArtifactBucket", contextParameters.artifactBucketName);
    const outputBucket = BucketOperations.importBucket(this, "OutputBucket", contextParameters.outputBucketName);

    const { readBucketArns = [], readWriteBucketArns = [] } = contextParameters;
    readBucketArns.push(artifactBucket.bucketArn);
    readWriteBucketArns.push(outputBucket.bucketArn);

    const batchRoles = this.getBatchRoles();
    for (const role of batchRoles) {
      BucketOperations.grantBucketAccess(this, role, readBucketArns, true);
      BucketOperations.grantBucketAccess(this, role, readWriteBucketArns, false);
    }
  }

  private renderBatch(id: string, vpc: IVpc, instanceTypes?: InstanceType[], computeType?: ComputeResourceType): Batch {
    return new Batch(this, id, {
      vpc,
      instanceTypes,
      computeType,
      launchTemplateData: LAUNCH_TEMPLATE,
      awsPolicyNames: ["AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"],
      resourceTags: this.nestedStackParent?.tags.tagValues(),
    });
  }

  private getBatchRoles(): IRole[] {
    const roles = [];
    if (this.batchOnDemand) {
      roles.push(this.batchOnDemand.role);
    }
    if (this.batchSpot) {
      roles.push(this.batchSpot.role);
    }
    return roles;
  }
}
