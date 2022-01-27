import { Arn, ArnFormat } from "aws-cdk-lib";
import { Bucket, IBucket } from "aws-cdk-lib/aws-s3";
import { IRole } from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";

export class BucketOperations {
  private static readonly importedBuckets: Record<string, IBucket> = {};

  public static grantBucketAccess(scope: Construct, role: IRole, bucketArns: string[], readOnly?: boolean): void {
    bucketArns.forEach((bucketArn) => {
      const arnComponents = Arn.split(bucketArn, ArnFormat.SLASH_RESOURCE_NAME);
      const bucketName = arnComponents.resource;
      const bucketPrefix = arnComponents.resourceName;
      const bucket = this.importBucket(scope, `${bucketName}Bucket`, bucketName);
      if (readOnly) {
        bucket.grantRead(role, bucketPrefix);
      } else {
        bucket.grantReadWrite(role, bucketPrefix);
      }
    });
  }

  public static importBucket(scope: Construct, bucketId: string, bucketName: string): IBucket {
    if (!this.importedBuckets[bucketId]) {
      this.importedBuckets[bucketId] = Bucket.fromBucketName(scope, bucketId, bucketName);
    }
    return this.importedBuckets[bucketId];
  }
}
