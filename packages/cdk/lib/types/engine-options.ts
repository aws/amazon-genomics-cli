import { RoleProps } from "aws-cdk-lib/aws-iam";
import { IMachineImage, IVpc, SubnetSelection } from "aws-cdk-lib/aws-ec2";
import { ContextAppParameters } from "../env";
import { Size } from "aws-cdk-lib";

export type PolicyOptions = Pick<RoleProps, "inlinePolicies" | "managedPolicies">;

export interface EngineOptions {
  /**
   * Policies to add to the task role for the engine.
   *
   * @default - No policies are added.
   */
  policyOptions: PolicyOptions;
  /**
   * VPC to run resources in.
   */
  readonly vpc: IVpc;
  /**
   * VPC subnets to run resources in
   */
  readonly subnets: SubnetSelection;
  /**
   * Filesystem provisioned throughput to use for EFS.
   */
  readonly iops?: Size;
  /**
   * Parameters determined by the context.
   */
  readonly contextParameters: ContextAppParameters;
  /**
   * The AMI to use for compute environments. Ignored for Fargate environments
   */
  readonly computeEnvImage?: IMachineImage;
}
