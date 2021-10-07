import { Construct, Fn, Names, Stack } from "monocdk";
import { ComputeEnvironment, ComputeResourceType, IComputeEnvironment, IJobQueue, JobQueue } from "monocdk/aws-batch";
import { CfnLaunchTemplate, InstanceType, IVpc } from "monocdk/aws-ec2";
import { CfnInstanceProfile, IManagedPolicy, IRole, ManagedPolicy, PolicyDocument, PolicyStatement, Role, ServicePrincipal } from "monocdk/aws-iam";
import { getInstanceTypesForBatch } from "../util/instance-types";

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
   * The tags to apply to any compute resources
   * @default none
   */
  resourceTags?: { [p: string]: string };
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
    return new Role(this, "BatchRole", {
      assumedBy: new ServicePrincipal("ec2.amazonaws.com"),
      inlinePolicies: {
        "ebs-autoscaling": new PolicyDocument({
          statements: [
            new PolicyStatement({
              actions: ["ec2:CreateTags", "ec2:DescribeVolumes", "ec2:CreateVolume", "ec2:AttachVolume", "ec2:DeleteVolume", "ec2:ModifyInstanceAttribute"],
              resources: ["*"],
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
