import * as cdk from '@aws-cdk/core';
import { ManagedPolicy, PolicyDocument } from '@aws-cdk/aws-iam';
import * as stmt from './policy-statements';
export class AgcPermissionsStack extends cdk.Stack {
  adminPolicy: ManagedPolicy;
  userPolicyPart1: ManagedPolicy;
  userPolicyPart2: ManagedPolicy;

  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    let agcAdminPolicy = new ManagedPolicy(this, 'agc-admin-policy', {
      description: "managed policy for amazon genomics cli admins"
    })

    let agcUserPolicyPart1 = new ManagedPolicy(this, 'agc-user-policy-part1', {
      description: "managed policy part 1 for amazon genomics cli users"
    });

    let agcUserPolicyPart2 = new ManagedPolicy(this, 'agc-user-policy-part2', {
      description: "managed policy part 2 for amazon genomics cli users"
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

    agcUserPolicyPart1.addStatements(
      // poweruser + iam permissions is sufficient
      ...perms.iam(),
      ...perms.sts(),
      ...perms.ec2(),
      ...perms.s3Create(),
      ...perms.s3Destroy(),
      ...perms.s3Read(),
      ...perms.s3Write(),
      ...perms.dynamodbRead(),
      ...perms.dynamodbWrite(),
      ...perms.ssmCreate(),
      ...perms.ssmRead(),
      ...perms.ssmDestroy(),
    );

    agcUserPolicyPart2.addStatements(
    // splitting user policy due to quota limit for PolicySize
        ...perms.cloudformationUser(),
        ...perms.batch(),
        ...perms.ecr(),
        ...perms.ecs(),
        ...perms.elb(),
        ...perms.apigw(),
        ...perms.efs(),
        ...perms.cloudmap(),
        ...perms.logs(),
        ...perms.route53(),
    )

    this.adminPolicy = agcAdminPolicy;
    this.userPolicyPart1 = agcUserPolicyPart1;
    this.userPolicyPart2 = agcUserPolicyPart2;

  }
}
