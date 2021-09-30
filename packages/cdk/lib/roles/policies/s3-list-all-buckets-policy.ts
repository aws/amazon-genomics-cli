import * as iam from "monocdk/aws-iam";
export class S3ListAllBucketsPolicy extends iam.PolicyDocument {
  constructor() {
    super({
      assignSids: true,
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["s3:ListAllMyBuckets"],
          resources: ["*"],
        }),
      ],
    });
  }
}
