import { JobDefinition, PlatformCapabilities } from "@aws-cdk/aws-batch-alpha";
import { FargatePlatformVersion } from "aws-cdk-lib/aws-ecs";
import { Construct } from "constructs";
import { createEcrImage } from "../../../util";
import { Batch } from "../../batch";
import { Engine, EngineProps } from "../engine";
import { EngineJobDefinition } from "../engine-job-definition";

export interface MiniWdlEngineProps extends EngineProps {
  readonly engineBatch: Batch;
  readonly workerBatch: Batch;
}

const MINIWDL_IMAGE_DESIGNATION = "miniwdl";

export class MiniWdlEngine extends Engine {
  readonly headJobDefinition: JobDefinition;
  private readonly engineMemoryMiB = 4096;
  private readonly volumeName = "efs";

  constructor(scope: Construct, id: string, props: MiniWdlEngineProps) {
    super(scope, id);

    const { vpc, rootDirS3Uri, engineBatch, workerBatch } = props;
    const fileSystem = this.createFileSystem(vpc);
    const accessPoint = this.createAccessPoint(fileSystem);

    fileSystem.connections.allowDefaultPortFromAnyIpv4();
    fileSystem.grant(engineBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");
    fileSystem.grant(workerBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");

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
          MINIWDL__AWS__FS: fileSystem.fileSystemId,
          MINIWDL__AWS__FSAP: accessPoint.accessPointId,
          MINIWDL__AWS__TASK_QUEUE: workerBatch.jobQueue.jobQueueArn,
          MINIWDL_S3_OUTPUT_URI: rootDirS3Uri,
        },
        volumes: [this.toVolume(fileSystem, accessPoint, this.volumeName)],
        mountPoints: [this.toMountPoint("/mnt/efs", this.volumeName)],
      },
    });
  }
}
