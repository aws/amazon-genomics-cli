import { RoleProps } from "monocdk/aws-iam";
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
   * VPC to run resources in.
   */
  readonly vpc: IVpc;
  /**
   * Parameters determined by the context.
   */
  readonly contextParameters: ContextAppParameters;
}
