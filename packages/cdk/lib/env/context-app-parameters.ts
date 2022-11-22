import { getEnvNumber, getEnvBoolOrDefault, getEnvString, getEnvStringListOrDefault, getEnvStringOrDefault } from "./";
import { InstanceType } from "aws-cdk-lib/aws-ec2";
import { Node } from "constructs";
import { ServiceContainer } from "../types";

const oneCpuUnit = 1024;
const oneGBinMiB = 1024;

export class ContextAppParameters {
  /**
   * Name of the project.
   */
  public readonly projectName: string;
  /**
   * Name of the context.
   */
  public readonly contextName: string;
  /**
   * The user's ID.
   */
  public readonly userId: string;
  /**
   * The user's email.
   */
  public readonly userEmail: string;

  /**
   * Bucket used to store outputs.
   */
  public readonly outputBucketName: string;
  /**
   * Bucket that stores artifacts.
   */
  public readonly artifactBucketName: string;
  /**
   * A list of ARNs that batch will access for workflow reads.
   */
  public readonly readBucketArns?: string[];
  /**
   * A list of ARNs that batch will access for workflow reads and writes.
   */
  public readonly readWriteBucketArns?: string[];

  /**
   * Name of the engine to run.
   */
  public readonly engineName: string;
  /**
   * Workflow language supported by the engine.
   */
  public readonly engineType: string;
  /**
   * Name of the filesystem type to use (e.g. EFS, S3).
   */
  public readonly filesystemType?: string;
  /**
   * Amount of provisioned IOPS to use.
   */
  public readonly fsProvisionedThroughput?: number;
  /**
   * Name of the engine ECR image.
   */
  public readonly engineDesignation: string;
  /**
   * Health check path for the engine.
   */
  public readonly engineHealthCheckPath: string;
  /**
   * Whether to enable workflow call caching for the engine.
   */
  public readonly callCachingEnabled: boolean;

  /**
   * Name of the WES adapter.
   */
  public readonly adapterName: string;
  /**
   * Name of the WES adapter ECR image.
   */
  public readonly adapterDesignation: string;

  /**
   * The maximum number of Amazon EC2 vCPUs that an environment can reach.
   */
  public readonly maxVCpus?: number;
  /**
   * Property to specify if the compute environment uses On-Demand or Spot compute resources.
   */
  public readonly requestSpotInstances: boolean;
  /**
   * The types of EC2 instances that may be launched in the compute environment.
   */
  public readonly instanceTypes?: InstanceType[];
  /**
   * If true, put EC2 instances into public subnets instead of private subnets.
   * This allows you to obtain significantly lower ongoing costs if used in conjunction with the usePublicSubnets option
   * for the associated account/core stack, which is enabled using `agc account activate --usePublicSubnets`.
   * Note that this option risks security vulnerabilities if security groups are manually modified.
   *
   * @default false
   */
  public readonly usePublicSubnets?: boolean;
  /**
   * AGC version being deployed.
   */
  public readonly agcVersion: string;

  /**
   * Map of custom tags to be applied to all the infrastructure in the context.
   */
  public readonly customTags: { [key: string]: string };

  constructor(node: Node) {
    const instanceTypeStrings = getEnvStringListOrDefault(node, "BATCH_COMPUTE_INSTANCE_TYPES");

    this.projectName = getEnvString(node, "PROJECT");
    this.contextName = getEnvString(node, "CONTEXT");
    this.userId = getEnvString(node, "USER_ID");
    this.userEmail = getEnvString(node, "USER_EMAIL");

    this.outputBucketName = getEnvString(node, "OUTPUT_BUCKET");
    this.artifactBucketName = getEnvString(node, "ARTIFACT_BUCKET");
    this.readBucketArns = getEnvStringListOrDefault(node, "READ_BUCKET_ARNS");
    this.readWriteBucketArns = getEnvStringListOrDefault(node, "READ_WRITE_BUCKET_ARNS");

    this.engineName = getEnvString(node, "ENGINE_NAME");
    this.filesystemType = getEnvStringOrDefault(node, "FILESYSTEM_TYPE", this.getDefaultFilesystem());
    this.fsProvisionedThroughput = getEnvNumber(node, "FS_PROVISIONED_THROUGHPUT");
    this.engineDesignation = getEnvString(node, "ENGINE_DESIGNATION");
    this.engineHealthCheckPath = getEnvStringOrDefault(node, "ENGINE_HEALTH_CHECK_PATH", "/engine/v1/status")!;
    this.callCachingEnabled = getEnvBoolOrDefault(node, "CALL_CACHING_ENABLED", true)!;

    this.adapterName = getEnvStringOrDefault(node, "ADAPTER_NAME", "wesAdapter")!;
    this.adapterDesignation = getEnvStringOrDefault(node, "ADAPTER_DESIGNATION", "wes")!;

    this.maxVCpus = getEnvNumber(node, "MAX_V_CPUS");
    this.requestSpotInstances = getEnvBoolOrDefault(node, "REQUEST_SPOT_INSTANCES", false)!;
    this.instanceTypes = instanceTypeStrings ? instanceTypeStrings.map((instanceType) => new InstanceType(instanceType.trim())) : undefined;

    this.usePublicSubnets = getEnvBoolOrDefault(node, "PUBLIC_SUBNETS", false);
    this.agcVersion = getEnvString(node, "AGC_VERSION");

    const tagsJson = getEnvStringOrDefault(node, "CUSTOM_TAGS");
    if (tagsJson != null) {
      this.customTags = JSON.parse(tagsJson);
    } else {
      this.customTags = {};
    }

    this.engineType = this.getEngineType();
  }

  public getContextBucketPath(): string {
    return `s3://${this.outputBucketName}/project/${this.projectName}/userid/${this.userId}/context/${this.contextName}`;
  }

  public getEngineBucketPath(): string {
    return `${this.getContextBucketPath()}/${this.engineName}-execution`;
  }

  /**
   * This function defines the container that server-based engines (like Toil
   * or Cromwell) will run their servers in. It is going to run on Fargate.
   */
  public getEngineContainer(jobQueueArn: string, additionalEnvVars?: { [key: string]: string }): ServiceContainer {
    return {
      serviceName: this.engineName,
      imageConfig: { designation: this.engineDesignation },
      containerPort: 8000,
      cpu: this.callCachingEnabled ? oneCpuUnit * 2 : oneCpuUnit / 2,
      memoryLimitMiB: this.callCachingEnabled ? oneGBinMiB * 16 : oneGBinMiB * 2,
      healthCheckPath: this.engineHealthCheckPath,
      environment: {
        S3BUCKET: this.outputBucketName,
        ROOT_DIR: this.getEngineBucketPath(),
        JOB_QUEUE_ARN: jobQueueArn,
        ...additionalEnvVars,
      },
    };
  }

  public getAdapterContainer(additionalEnvVars?: { [key: string]: string }): ServiceContainer {
    return {
      serviceName: this.adapterName,
      imageConfig: { designation: this.adapterDesignation },
      cpu: oneCpuUnit / 2,
      memoryLimitMiB: oneGBinMiB * 4,
      environment: {
        PROJECT_NAME: this.projectName,
        CONTEXT_NAME: this.contextName,
        USER_ID: this.userId,
        ENGINE_NAME: this.engineName,
        ...additionalEnvVars,
      },
    };
  }

  public getDefaultFilesystem(): string {
    let defFilesystem: string;
    switch (this.engineName) {
      case "cromwell":
        defFilesystem = "S3";
        break;
      case "nextflow":
        defFilesystem = "S3";
        break;
      case "miniwdl":
        defFilesystem = "EFS";
        break;
      case "snakemake":
        defFilesystem = "EFS";
        break;
      case "toil":
        defFilesystem = "S3";
        break;
      default:
        throw Error(`Engine '${this.engineName}' is not supported`);
    }
    return defFilesystem;
  }

  public getEngineType(): string {
    let engineType: string;
    switch (this.engineName.toLowerCase()) {
      case "cromwell":
        engineType = "wdl";
        break;
      case "nextflow":
        engineType = "nextflow";
        break;
      case "miniwdl":
        engineType = "wdl";
        break;
      case "snakemake":
        engineType = "snakemake";
        break;
      case "toil":
        engineType = "cwl";
        break;
      default:
        engineType = "(unknown)";
        break;
    }
    return engineType;
  }
}
