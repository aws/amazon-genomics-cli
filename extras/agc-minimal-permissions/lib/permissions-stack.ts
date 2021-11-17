import * as cdk from '@aws-cdk/core';
import { ManagedPolicy, PolicyDocument } from '@aws-cdk/aws-iam';
import * as stmt from './policy-statements';
export class AgcPermissionsStack extends cdk.Stack {
  adminPolicy: ManagedPolicy;
  userPolicy: ManagedPolicy;

  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    let agcAdminPolicy = new ManagedPolicy(this, 'agc-admin-policy', {
      description: "managed policy for amazon genomics cli admins"
    })

    let agcUserPolicy = new ManagedPolicy(this, 'agc-user-policy', {
      description: "managed policy for amazon genomics cli users"
    });

    let perms = new stmt.AgcPermissions(this);

    agcAdminPolicy.addStatements(
      // explicit permissions
      ...perms.vpc(),
      ...perms.s3Create(),
      ...perms.s3Destroy(),
      ...perms.s3Read(),
      ...perms.s3Write(),
      ...perms.s3CDK(),
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
    );

    agcUserPolicy.addStatements(
      // poweruser + iam permissions is sufficient
      ...perms.iam(),
        
      ...perms.ec2(),
      ...perms.s3Read(),
      ...perms.s3Write(),
      ...perms.s3CDK(),
      ...perms.dynamodbRead(),
      ...perms.dynamodbWrite(),
      ...perms.ssmRead(),
      ...perms.cloudformationUser(),
      ...perms.batch(),
      ...perms.ecs(),
      ...perms.elb(),
      ...perms.apigw(),
      ...perms.efs(),
      ...perms.cloudmap(),
      ...perms.logs(),
      ...perms.route53(),
    );

    this.adminPolicy = agcAdminPolicy;
    this.userPolicy = agcUserPolicy;

  }
}
