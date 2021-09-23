import { RoleProps } from "monocdk/aws-iam";
import { IJobQueue } from "monocdk/aws-batch";
import { IVpc } from "monocdk/aws-ec2";
import { ContextAppParameters } from "../env";

export type PolicyOptions = Pick<RoleProps, "inlinePolicies" | "managedPolicies">;

export interface EngineOptions {
  /**
   * Policies to add to the task role for the engine.
   *
   * @default - No policies are added.
   */
  policyOptions: PolicyOptions;
  /**
   * AWS Batch JobQueue to use for running workflows.
   */
  readonly jobQueue: IJobQueue;
  /**
   * VPC to run resources in.
   */
  readonly vpc: IVpc;
  /**
   * Parameters determined by the context.
   */
  readonly contextParameters: ContextAppParameters;
}
