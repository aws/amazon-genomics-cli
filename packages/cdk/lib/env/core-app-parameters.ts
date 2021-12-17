import { APP_ENV_NAME, APP_NAME } from "../constants";
import { Node } from "constructs";
import { getEnvBoolOrDefault, getEnvStringOrDefault } from "./index";

export class CoreAppParameters {
  /**
   * VPC to run resources in.
   */
  public readonly vpcId?: string;
  /**
   * Name of the application bucket.
   */
  public readonly bucketName: string;
  /**
   * Whether to create a new application bucket.
   */
  public readonly createNewBucket: boolean;

  constructor(node: Node, account: string, region: string) {
    this.vpcId = getEnvStringOrDefault(node, "VPC_ID");
    this.bucketName = getEnvStringOrDefault(node, `${APP_ENV_NAME}_BUCKET_NAME`, `${APP_NAME}-${account}-${region}`)!;
    this.createNewBucket = getEnvBoolOrDefault(node, `CREATE_${APP_ENV_NAME}_BUCKET`, true)!;
  }
}
