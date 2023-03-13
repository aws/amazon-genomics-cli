import { Construct } from "constructs";
import { JobDefinition, PlatformCapabilities } from "@aws-cdk/aws-batch-alpha";
import { createEcrImage } from "../../../util";
import { EngineJobDefinition } from "../engine-job-definition";
import { Engine, EngineProps } from "../engine";
import { Batch } from "../../batch";
import { FargatePlatformVersion } from "aws-cdk-lib/aws-ecs";
import { AccessPoint, FileSystem } from "aws-cdk-lib/aws-efs";
import { Size } from "aws-cdk-lib";
import { SecurityGroup } from "aws-cdk-lib/aws-ec2";

export interface SnakemakeEngineProps extends EngineProps {
  readonly engineBatch: Batch;
  readonly workerBatch: Batch;
  readonly iops?: Size;
}

const SNAKEMAKE_IMAGE_DESIGNATION = "snakemake";

export class SnakemakeEngine extends Engine {
  readonly headJobDefinition: JobDefinition;
  private readonly volumeName = "efs";
  private readonly engineMemoryMiB = 4096;
  public readonly fsap: AccessPoint;
  public readonly fileSystem: FileSystem;
  private readonly securityGroup: SecurityGroup;

  constructor(scope: Construct, id: string, props: SnakemakeEngineProps) {
    super(scope, id);

    const { vpc, subnets, iops, engineBatch, workerBatch } = props;
    if (iops?.toMebibytes() == 0 || iops == undefined) {
      this.fileSystem = this.createFileSystemDefaultThroughput(vpc, subnets);
    } else {
      this.fileSystem = this.createFileSystemIOPS(vpc, subnets, iops);
    }
    this.fsap = this.createAccessPoint(this.fileSystem);
    this.fileSystem.connections.allowDefaultPortFrom(engineBatch.computeEnvironment.connections);
    this.fileSystem.connections.allowDefaultPortFrom(workerBatch.computeEnvironment.connections);
    this.fileSystem.grant(engineBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");
    this.fileSystem.grant(workerBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");
    this.headJobDefinition = new EngineJobDefinition(this, "SnakemakeHeadJobDef", {
      logGroup: this.logGroup,
      platformCapabilities: [PlatformCapabilities.FARGATE],
      container: {
        memoryLimitMiB: this.engineMemoryMiB,
        jobRole: engineBatch.role,
        executionRole: engineBatch.role,
        image: createEcrImage(this, SNAKEMAKE_IMAGE_DESIGNATION),
        platformVersion: FargatePlatformVersion.VERSION1_4,
        command: [],
        environment: {
          SM__AWS__FS: this.fileSystem.fileSystemId,
          SM__AWS__FSAP: this.fsap.accessPointId,
          SM__AWS__TASK_QUEUE: workerBatch.jobQueue.jobQueueArn,
          SM_S3_OUTPUT_URI: props.rootDirS3Uri,
        },
        volumes: [this.toVolume(this.fileSystem, this.fsap, this.volumeName)],
        mountPoints: [this.toMountPoint("/mnt/efs", this.volumeName)],
      },
    });
  }
}
