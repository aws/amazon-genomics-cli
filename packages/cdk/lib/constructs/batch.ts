import { Fn, Names, Stack } from "aws-cdk-lib";
import { ComputeEnvironment, ComputeResourceType, IComputeEnvironment, IJobQueue, JobQueue } from "@aws-cdk/aws-batch-alpha";
import { CfnLaunchTemplate, InstanceType, IVpc, SubnetType } from "aws-cdk-lib/aws-ec2";
import {
  CfnInstanceProfile,
  Grant,
  IGrantable,
  IManagedPolicy,
  IRole,
  ManagedPolicy,
  PolicyDocument,
  PolicyStatement,
  Role,
  ServicePrincipal,
} from "aws-cdk-lib/aws-iam";
import { getInstanceTypesForBatch } from "../util/instance-types";
import { batchArn, ec2Arn } from "../util";
import { APP_NAME, APP_TAG_KEY } from "../../lib/constants";
import { Construct } from "constructs";

export interface ComputeOptions {
  /**
   * The VPC to run the batch in.
   */
  vpc: IVpc;
  /**
   * User data to make available to the instances.
   *
   * @default none
   */
  launchTemplateData?: string;
  /**
   * The type of compute environment.
   *
   * @default ON_DEMAND
   */
  computeType?: ComputeResourceType;
  /**
   * The types of EC2 instances that may be launched in the compute environment.
   *
   * This property is only valid when using a non-Fargate compute type.
   *
   * @default optimal
   */
  instanceTypes?: InstanceType[];

  /**
   * The maximum number of EC2 vCPUs that a compute environment can reach.
   *
   * Each vCPU is equivalent to 1,024 CPU shares.
   *
   * @default aws-batch:{@link ComputeResources#maxvCpus}
   */
  maxVCpus?: number;

  /**
   * The tags to apply to any compute resources
   * @default none
   */
  resourceTags?: { [p: string]: string };

  /**
   * If true, put EC2 instances into public subnets instead of private subnets. This should be used in conjunction
   * with the ENV_PUBLIC_SUBNETS option
   * This will result in significantly lower ongoing costs when no job is running.
   */
  publicSubnets?: boolean;
}

export interface BatchProps extends ComputeOptions {
  /**
   * The names of AWS managed policies to attach to the batch role.
   *
   * The batch role already includes "service-role/AmazonECSTaskExecutionRolePolicy" or
   * "service-role/AmazonEC2ContainerServiceforEC2Role" depending on whether the compute
   * type is Fargate or not.
   *
   * @default - No additional policies are added to the role
   */
  awsPolicyNames?: string[];
}

const defaultComputeType = ComputeResourceType.ON_DEMAND;

export class Batch extends Construct {
  public readonly role: IRole;
  public readonly computeEnvironment: IComputeEnvironment;
  public readonly jobQueue: IJobQueue;

  constructor(scope: Construct, id: string, props: BatchProps) {
    super(scope, id);

    this.role = this.renderRole(props.computeType, props.awsPolicyNames);
    this.computeEnvironment = this.renderComputeEnvironment(props);

    this.jobQueue = new JobQueue(this, "JobQueue", {
      computeEnvironments: [
        {
          order: 1,
          computeEnvironment: this.computeEnvironment,
        },
      ],
    });
  }

  public grantJobAdministration(grantee: IGrantable, jobDefinitionName = "*"): Grant {
    return Grant.addToPrincipal({
      grantee: grantee,
      actions: ["batch:SubmitJob", "batch:TerminateJob"],
      resourceArns: [this.jobQueue.jobQueueArn, batchArn(this, "job-definition", jobDefinitionName)],
    });
  }

  private renderRole(computeType?: ComputeResourceType, awsPolicyNames?: string[]): IRole {
    const awsPolicies = awsPolicyNames?.map((policyName) => ManagedPolicy.fromAwsManagedPolicyName(policyName));
    if (computeType == ComputeResourceType.FARGATE || computeType == ComputeResourceType.FARGATE_SPOT) {
      return this.renderEcsRole(awsPolicies);
    }
    return this.renderEc2Role(awsPolicies);
  }

  private renderEcsRole(managedPolicies?: IManagedPolicy[]): IRole {
    return new Role(this, "BatchRole", {
      assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
      managedPolicies: [...(managedPolicies ?? []), ManagedPolicy.fromAwsManagedPolicyName("service-role/AmazonECSTaskExecutionRolePolicy")],
    });
  }

  private renderEc2Role(managedPolicies?: IManagedPolicy[]): IRole {
    const volumeArn = ec2Arn(this, "volume");

    return new Role(this, "BatchRole", {
      assumedBy: new ServicePrincipal("ec2.amazonaws.com"),
      inlinePolicies: {
        "ebs-autoscaling": new PolicyDocument({
          statements: [
            new PolicyStatement({
              actions: ["ec2:DescribeVolumes", "ec2:CreateVolume", "ec2:CreateTags"],
              resources: [volumeArn],
            }),
            new PolicyStatement({
              actions: ["ec2:AttachVolume", "ec2:ModifyInstanceAttribute"],
              resources: [ec2Arn(this, "instance"), volumeArn],
            }),
            new PolicyStatement({
              actions: ["ec2:DeleteVolume"],
              resources: [volumeArn],
              conditions: {
                StringEquals: {
                  [`aws:ResourceTag/${APP_TAG_KEY}`]: APP_NAME,
                },
              },
            }),
          ],
        }),
      },
      managedPolicies: [...(managedPolicies ?? []), ManagedPolicy.fromAwsManagedPolicyName("service-role/AmazonEC2ContainerServiceforEC2Role")],
    });
  }

  private renderComputeEnvironment(options: ComputeOptions): IComputeEnvironment {
    const computeType = options.computeType || defaultComputeType;
    if (computeType == ComputeResourceType.FARGATE || computeType == ComputeResourceType.FARGATE_SPOT) {
      return new ComputeEnvironment(this, "ComputeEnvironment", {
        computeResources: {
          vpc: options.vpc,
          type: computeType,
          maxvCpus: options.maxVCpus,
          vpcSubnets: {
            // Even if we use public subnets, CDK will assign security groups that don't allow inbound
            subnetType: options.publicSubnets ? SubnetType.PUBLIC : SubnetType.PRIVATE_WITH_NAT,
          },
        },
      });
    }

    const launchTemplate = options.launchTemplateData
      ? new CfnLaunchTemplate(this, "LaunchTemplate", {
          launchTemplateName: Names.uniqueId(this),
          launchTemplateData: {
            userData: Fn.base64(options.launchTemplateData),
          },
        })
      : undefined;

    const instanceProfile = new CfnInstanceProfile(this, "ComputeProfile", {
      roles: [this.role.roleName],
    });
    return new ComputeEnvironment(this, "ComputeEnvironment", {
      computeResources: {
        vpc: options.vpc,
        type: computeType,
        maxvCpus: options.maxVCpus,
        instanceRole: instanceProfile.attrArn,
        instanceTypes: getInstanceTypesForBatch(options.instanceTypes, computeType, Stack.of(this).region),
        launchTemplate: launchTemplate && {
          launchTemplateName: launchTemplate.launchTemplateName!,
        },
        computeResourcesTags: options.resourceTags,
      },
    });
  }
}
