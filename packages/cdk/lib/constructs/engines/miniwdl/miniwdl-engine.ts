import { Construct, RemovalPolicy } from "monocdk";
import { JobDefinition, PlatformCapabilities } from "monocdk/aws-batch";
import { IVpc } from "monocdk/aws-ec2";
import { AccessPoint, FileSystem, PerformanceMode } from "monocdk/aws-efs";
import { FargatePlatformVersion } from "monocdk/aws-ecs";
import { Batch } from "../../batch";
import { Engine, EngineProps } from "../engine";
import { EngineJobDefinition } from "../engine-job-definition";
import { createEcrImage } from "../../../util";

export interface MiniWdlEngineProps extends EngineProps {
  readonly engineBatch: Batch;
  readonly workerBatch: Batch;
}

const MINIWDL_IMAGE_DESIGNATION = "miniwdl";

export class MiniWdlEngine extends Engine {
  readonly headJobDefinition: JobDefinition;
  private readonly volumeName = "efs";

  constructor(scope: Construct, id: string, props: MiniWdlEngineProps) {
    super(scope, id);

    const { vpc, outputBucketName, engineBatch, workerBatch } = props;
    const fileSystem = this.createFileSystem(vpc);
    const accessPoint = this.createAccessPoint(fileSystem);

    fileSystem.connections.allowDefaultPortFromAnyIpv4();
    fileSystem.grant(engineBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");
    fileSystem.grant(workerBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");

    this.headJobDefinition = new EngineJobDefinition(this, "MiniwdlHeadJobDef", {
      logGroup: this.logGroup,
      platformCapabilities: [PlatformCapabilities.FARGATE],
      container: {
        jobRole: engineBatch.role,
        executionRole: engineBatch.role,
        image: createEcrImage(this, MINIWDL_IMAGE_DESIGNATION),
        platformVersion: FargatePlatformVersion.VERSION1_4,
        environment: {
          MINIWDL__AWS__FS: fileSystem.fileSystemId,
          MINIWDL__AWS__FSAP: accessPoint.accessPointId,
          MINIWDL__AWS__TASK_QUEUE: workerBatch.jobQueue.jobQueueArn,
          MINIWDL_S3_OUTPUT_URI: `s3://${outputBucketName}/miniwdl`,
        },
        volumes: [this.toVolume(fileSystem, accessPoint)],
        mountPoints: [this.toMountPoint("/mnt/efs")],
      },
    });
  }

  private toMountPoint(containerPath: string) {
    return {
      sourceVolume: this.volumeName,
      containerPath: containerPath,
      readOnly: false,
    };
  }

  private toVolume(fileSystem: FileSystem, accessPoint: AccessPoint) {
    return {
      name: this.volumeName,
      efsVolumeConfiguration: {
        fileSystemId: fileSystem.fileSystemId,
        transitEncryption: "ENABLED",
        authorizationConfig: {
          accessPointId: accessPoint.accessPointId,
          iam: "ENABLED",
        },
      },
    };
  }

  private createAccessPoint(fileSystem: FileSystem) {
    return new AccessPoint(this, "AccessPoint", {
      fileSystem: fileSystem,
      posixUser: {
        uid: "0",
        gid: "0",
      },
    });
  }

  private createFileSystem(vpc: IVpc) {
    return new FileSystem(this, "FileSystem", {
      vpc: vpc,
      encrypted: true,
      performanceMode: PerformanceMode.MAX_IO,
      removalPolicy: RemovalPolicy.DESTROY,
    });
  }
}
