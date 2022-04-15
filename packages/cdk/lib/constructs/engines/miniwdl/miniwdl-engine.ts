import { JobDefinition, PlatformCapabilities } from "@aws-cdk/aws-batch-alpha";
import { FargatePlatformVersion } from "aws-cdk-lib/aws-ecs";
import { Construct } from "constructs";
import { createEcrImage } from "../../../util";
import { Batch } from "../../batch";
import { Engine, EngineProps } from "../engine";
import { EngineJobDefinition } from "../engine-job-definition";
import { Size } from "aws-cdk-lib";
import { FileSystem } from "aws-cdk-lib/aws-efs";

export interface MiniWdlEngineProps extends EngineProps {
  readonly engineBatch: Batch;
  readonly workerBatch: Batch;
  readonly iops?: Size;
}

const MINIWDL_IMAGE_DESIGNATION = "miniwdl";

export class MiniWdlEngine extends Engine {
  readonly headJobDefinition: JobDefinition;
  private readonly engineMemoryMiB = 4096;
  private readonly volumeName = "efs";
  readonly fileSystem: FileSystem;

  constructor(scope: Construct, id: string, props: MiniWdlEngineProps) {
    super(scope, id);

    const { vpc, subnets, iops, rootDirS3Uri, engineBatch, workerBatch } = props;
    if (iops?.toMebibytes() == 0 || iops == undefined) {
      this.fileSystem = this.createFileSystemDefaultThroughput(vpc, subnets);
    } else {
      this.fileSystem = this.createFileSystemIOPS(vpc, subnets, iops);
    }
    const accessPoint = this.createAccessPoint(this.fileSystem);

    this.fileSystem.connections.allowDefaultPortFromAnyIpv4();
    this.fileSystem.grant(engineBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");
    this.fileSystem.grant(workerBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");

    this.headJobDefinition = new EngineJobDefinition(this, "MiniwdlHeadJobDef", {
      logGroup: this.logGroup,
      platformCapabilities: [PlatformCapabilities.FARGATE],
      container: {
        memoryLimitMiB: this.engineMemoryMiB,
        jobRole: engineBatch.role,
        executionRole: engineBatch.role,
        image: createEcrImage(this, MINIWDL_IMAGE_DESIGNATION),
        platformVersion: FargatePlatformVersion.VERSION1_4,
        environment: {
          MINIWDL__AWS__FS: this.fileSystem.fileSystemId,
          MINIWDL__AWS__FSAP: accessPoint.accessPointId,
          MINIWDL__AWS__TASK_QUEUE: workerBatch.jobQueue.jobQueueArn,
          MINIWDL_S3_OUTPUT_URI: rootDirS3Uri,
        },
        volumes: [this.toVolume(this.fileSystem, accessPoint, this.volumeName)],
        mountPoints: [this.toMountPoint("/mnt/efs", this.volumeName)],
      },
    });
  }
}
