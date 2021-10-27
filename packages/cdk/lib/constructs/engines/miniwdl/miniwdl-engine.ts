import { Construct, RemovalPolicy } from "monocdk";
import { CfnJobDefinition, JobDefinition, PlatformCapabilities } from "monocdk/aws-batch";
import { IVpc } from "monocdk/aws-ec2";
import { createEcrImage, renderBatchLogConfiguration } from "../../../util";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { LogGroup } from "monocdk/aws-logs";
import { AccessPoint, FileSystem } from "monocdk/aws-efs";
import { FargatePlatformVersion } from "monocdk/aws-ecs";
import { Batch } from "../../batch";

export interface MiniWdlEngineProps {
  readonly vpc: IVpc;
  readonly outputBucketName: string;
  readonly engineBatch: Batch;
  readonly workerBatch: Batch;
}

const MINIWDL_IMAGE_DESIGNATION = "miniwdl";

export class MiniWdlEngine extends Construct {
  readonly headJobDefinition: JobDefinition;
  readonly logGroup: ILogGroup;

  constructor(scope: Construct, id: string, props: MiniWdlEngineProps) {
    super(scope, id);

    const { vpc, outputBucketName, engineBatch, workerBatch } = props;

    this.logGroup = new LogGroup(this, "EngineLogGroup");

    const fileSystem = new FileSystem(this, "FileSystem", {
      vpc: vpc,
      encrypted: true,
      removalPolicy: RemovalPolicy.DESTROY,
    });

    const accessPoint = new AccessPoint(this, "AccessPoint", {
      fileSystem: fileSystem,
    });

    fileSystem.connections.allowDefaultPortFromAnyIpv4();
    fileSystem.grant(engineBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");
    fileSystem.grant(workerBatch.role, "elasticfilesystem:DescribeMountTargets", "elasticfilesystem:DescribeFileSystems");

    const volumeName = "efs";
    this.headJobDefinition = new JobDefinition(this, "MiniwdlHeadJobDef", {
      platformCapabilities: [PlatformCapabilities.FARGATE],
      container: {
        logConfiguration: renderBatchLogConfiguration(this, this.logGroup),
        jobRole: engineBatch.role,
        executionRole: engineBatch.role,
        image: createEcrImage(this, MINIWDL_IMAGE_DESIGNATION),
        platformVersion: FargatePlatformVersion.VERSION1_4,
        command: [],
        environment: {
          MINIWDL__AWS__FS: fileSystem.fileSystemId,
          MINIWDL__AWS__FSAP: accessPoint.accessPointId,
          MINIWDL__AWS__TASK_QUEUE: workerBatch.jobQueue.jobQueueArn,
          MINIWDL_S3_OUTPUT_URI: `s3://${outputBucketName}/miniwdl`,
        },
        volumes: [
          {
            name: volumeName,
            efsVolumeConfiguration: {
              fileSystemId: fileSystem.fileSystemId,
              transitEncryption: "ENABLED",
              authorizationConfig: {
                accessPointId: accessPoint.accessPointId,
                iam: "ENABLED",
              },
            },
          },
        ],
        mountPoints: [
          {
            sourceVolume: volumeName,
            containerPath: "/mnt/efs",
            readOnly: false,
          },
        ],
      },
    });

    const cfnJobDef = this.headJobDefinition.node.defaultChild as CfnJobDefinition;

    //Removing old method for specifying resources. Using newer ResourceRequirements.
    cfnJobDef.addPropertyDeletionOverride("ContainerProperties.Vcpus");
    cfnJobDef.addPropertyDeletionOverride("ContainerProperties.Memory");
    cfnJobDef.addPropertyOverride("ContainerProperties.ResourceRequirements", [
      { Type: "VCPU", Value: "1" },
      { Type: "MEMORY", Value: "2048" },
    ]);
  }
}
