import { NestedStack, NestedStackProps } from "monocdk";
import { InstanceType, IVpc } from "monocdk/aws-ec2";
import { Construct } from "constructs";
import { LAUNCH_TEMPLATE } from "../../constants";
import { Batch, ComputeType } from "../../constructs";
import { ContextAppParameters } from "../../env";
import { BucketOperations } from "../../../common/BucketOperations";

export interface BatchStackProps extends NestedStackProps {
  /**
   * VPC to run resources in.
   */
  readonly vpc: IVpc;
  /**
   * Parameters determined by the context.
   */
  readonly contextParameters: ContextAppParameters;
}

export class BatchStack extends NestedStack {
  public readonly batchWorkers: Batch;
  public readonly batchHead: Batch;

  constructor(scope: Construct, id: string, props: BatchStackProps) {
    super(scope, id, props);

    const { vpc, contextParameters } = props;

    this.batchWorkers = this.batchHead = this.renderBatch("TaskBatch", vpc, contextParameters.instanceTypes, ComputeType.ON_DEMAND);
    if (contextParameters.requestSpotInstances) {
      this.batchWorkers = this.renderBatch("TaskBatchSpot", vpc, contextParameters.instanceTypes, ComputeType.SPOT);
    }

    const artifactBucket = BucketOperations.importBucket(this, "ArtifactBucket", contextParameters.artifactBucketName);
    const outputBucket = BucketOperations.importBucket(this, "OutputBucket", contextParameters.outputBucketName);

    BucketOperations.grantBucketAccess(this, this.batchWorkers.role, (contextParameters.readBucketArns ?? []).concat(artifactBucket.bucketArn), true);
    BucketOperations.grantBucketAccess(this, this.batchHead.role, (contextParameters.readBucketArns ?? []).concat(artifactBucket.bucketArn), true);

    BucketOperations.grantBucketAccess(this, this.batchWorkers.role, (contextParameters.readWriteBucketArns ?? []).concat(outputBucket.bucketArn));
    BucketOperations.grantBucketAccess(this, this.batchHead.role, (contextParameters.readWriteBucketArns ?? []).concat(outputBucket.bucketArn));

    artifactBucket.grantRead(this.batchWorkers.role);
    artifactBucket.grantRead(this.batchHead.role);
    outputBucket.grantReadWrite(this.batchWorkers.role);
    outputBucket.grantReadWrite(this.batchHead.role);
  }

  private renderBatch(id: string, vpc: IVpc, instanceTypes?: InstanceType[], computeType?: ComputeType): Batch {
    return new Batch(this, id, {
      vpc,
      instanceTypes,
      computeType,
      launchTemplateData: LAUNCH_TEMPLATE,
      awsPolicyNames: ["AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"],
      resourceTags: this.nestedStackParent?.tags.tagValues(),
    });
  }
}
