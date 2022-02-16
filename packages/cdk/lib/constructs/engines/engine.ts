import { ILogGroup, LogGroup } from "aws-cdk-lib/aws-logs";
import { Construct } from "constructs";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { AccessPoint, FileSystem, PerformanceMode } from "aws-cdk-lib/aws-efs";
import { RemovalPolicy } from "aws-cdk-lib";
import { MountPoint, Volume } from "aws-cdk-lib/aws-ecs";

export interface EngineProps {
  readonly vpc: IVpc;
  readonly rootDirS3Uri: string;
}

export class Engine extends Construct {
  readonly logGroup: ILogGroup;

  constructor(scope: Construct, id: string) {
    super(scope, id);
    this.logGroup = new LogGroup(this, "EngineLogGroup");
  }

  protected toMountPoint(containerPath: string, volumeName: string): MountPoint {
    return {
      sourceVolume: volumeName,
      containerPath: containerPath,
      readOnly: false,
    };
  }

  protected toVolume(fileSystem: FileSystem, accessPoint: AccessPoint, volumeName: string): Volume {
    return {
      name: volumeName,
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

  protected createAccessPoint(fileSystem: FileSystem): AccessPoint {
    return new AccessPoint(this, "AccessPoint", {
      fileSystem: fileSystem,
      posixUser: {
        uid: "0",
        gid: "0",
      },
    });
  }

  protected createAccessPointcreateFileSystem(vpc: IVpc): FileSystem {
    return new FileSystem(this, "FileSystem", {
      vpc: vpc,
      encrypted: true,
      performanceMode: PerformanceMode.MAX_IO,
      removalPolicy: RemovalPolicy.DESTROY,
    });
  }

  protected createFileSystem(vpc: IVpc): FileSystem {
    return new FileSystem(this, "FileSystem", {
      vpc: vpc,
      encrypted: true,
      performanceMode: PerformanceMode.MAX_IO,
      removalPolicy: RemovalPolicy.DESTROY,
    });
  }
}
