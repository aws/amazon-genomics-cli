import { Arn, NestedStack, NestedStackProps } from "monocdk";
import { InstanceType, IVpc } from "monocdk/aws-ec2";
import { Construct } from "constructs";
import { LAUNCH_TEMPLATE } from "../../constants";
import { Bucket, IBucket } from "monocdk/aws-s3";
import { Batch, ComputeType } from "../../constructs";
import { ContextAppParameters } from "../../env";

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
  private readonly importedBuckets: Record<string, IBucket> = {};

  constructor(scope: Construct, id: string, props: BatchStackProps) {
    super(scope, id, props);

    const { vpc, contextParameters } = props;

    this.batchWorkers = this.batchHead = this.renderBatch("TaskBatch", vpc, contextParameters.instanceTypes, ComputeType.ON_DEMAND);
    if (contextParameters.requestSpotInstances) {
      this.batchWorkers = this.renderBatch("TaskBatchSpot", vpc, contextParameters.instanceTypes, ComputeType.SPOT);
    }

    this.grantBucketAccess(contextParameters.readBucketArns ?? [], true);
    this.grantBucketAccess(contextParameters.readWriteBucketArns ?? []);
    const artifactBucket = this.importBucket("ArtifactBucket", contextParameters.artifactBucketName);
    artifactBucket.grantRead(this.batchWorkers.role);
    artifactBucket.grantRead(this.batchHead.role);
    const outputBucket = this.importBucket("OutputBucket", contextParameters.outputBucketName);
    outputBucket.grantReadWrite(this.batchWorkers.role);
    outputBucket.grantReadWrite(this.batchHead.role);
  }

  private grantBucketAccess(bucketArns: string[], readOnly?: boolean): void {
    bucketArns.forEach((bucketArn) => {
      const arnComponents = Arn.parse(bucketArn);
      const bucketName = arnComponents.resource;
      const bucketPrefix = arnComponents.resourceName;
      const bucket = this.importBucket(`${bucketName}Bucket`, bucketName);
      if (readOnly) {
        bucket.grantRead(this.batchWorkers.role, bucketPrefix);
        bucket.grantRead(this.batchHead.role, bucketPrefix);
      } else {
        bucket.grantReadWrite(this.batchWorkers.role, bucketPrefix);
        bucket.grantReadWrite(this.batchHead.role, bucketPrefix);
      }
    });
  }

  private importBucket(bucketId: string, bucketName: string): IBucket {
    if (!this.importedBuckets[bucketId]) {
      this.importedBuckets[bucketId] = Bucket.fromBucketName(this, bucketId, bucketName);
    }
    return this.importedBuckets[bucketId];
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
