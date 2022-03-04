import * as cdk from '@aws-cdk/core';
import { ManagedPolicy, PolicyDocument } from '@aws-cdk/aws-iam';
import * as stmt from './policy-statements';
export class AgcPermissionsStack extends cdk.Stack {
  adminPolicy: ManagedPolicy;
  userPolicyCDK: ManagedPolicy;
  userPolicy: ManagedPolicy;

  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    let agcAdminPolicy = new ManagedPolicy(this, 'agc-admin-policy', {
      description: "managed policy for amazon genomics cli admins"
    })

    let agcUserPolicyCDK = new ManagedPolicy(this, 'agc-user-policy-cdk', {
      description: "managed policy for amazon genomics cli users to run cdk"
    });

    let agcUserPolicy = new ManagedPolicy(this, 'agc-user-policy', {
      description: "managed policy part 2 for amazon genomics cli users to run agc"
    });

    let perms = new stmt.AgcPermissions(this);

    agcAdminPolicy.addStatements(
      // explicit permissions
      ...perms.s3Create(),
      ...perms.s3Destroy(),
      ...perms.s3Read(),
      ...perms.s3Write(),
      ...perms.dynamodbCreate(),
      ...perms.dynamodbRead(),
      ...perms.dynamodbWrite(),
      ...perms.dynamodbDestroy(),
      ...perms.ssmCreate(),
      ...perms.ssmRead(),
      ...perms.ssmDestroy(),
      ...perms.cloudformationAdmin(),
      ...perms.ecr(),
      ...perms.deactivate(),
      ...perms.sts(),
      ...perms.iam(),
    );

    agcUserPolicyCDK.addStatements(
      ...perms.iam(),
      ...perms.sts(),
      ...perms.ec2(),
      ...perms.s3Create(),
      ...perms.s3Destroy(),
      ...perms.s3Write(),
      ...perms.ssmCreate(),
      ...perms.ssmDestroy(),
      ...perms.ecs(),
      ...perms.elb(),
      ...perms.apigw(),
      ...perms.route53(),
      ...perms.cloudformationUser(),
    );

    agcUserPolicy.addStatements(
      ...perms.dynamodbRead(),
      ...perms.dynamodbWrite(),
      ...perms.s3Read(),
      ...perms.ssmRead(),
      ...perms.batch(),
      ...perms.ecr(),
      ...perms.efs(),
      ...perms.cloudmap(),
      ...perms.logs(),
    )

    this.adminPolicy = agcAdminPolicy;
    this.userPolicyCDK = agcUserPolicyCDK;
    this.userPolicy = agcUserPolicy;

  }
}
