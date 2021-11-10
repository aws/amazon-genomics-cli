import { NestedStackProps, RemovalPolicy } from "monocdk";
import { IVpc } from "monocdk/aws-ec2";
import { FargateTaskDefinition, LogDriver } from "monocdk/aws-ecs";
import { Construct } from "constructs";
import { ApiProxy, SecureService } from "../../constructs";
import { IRole } from "monocdk/aws-iam";
import { createEcrImage, renderPythonLambda, renderServiceWithTaskDefinition } from "../../util";
import { Bucket } from "monocdk/aws-s3";
import { FileSystem } from "monocdk/aws-efs";
import { EngineOptions, ServiceContainer } from "../../types";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { LogGroup } from "monocdk/aws-logs";
import { EngineOutputs, NestedEngineStack } from "./nested-engine-stack";
import { CromwellEngineRole } from "../../roles/cromwell-engine-role";
import { CromwellAdapterRole } from "../../roles/cromwell-adapter-role";
import { wesAdapterSourcePath } from "../../constants";

export interface CromwellEngineStackProps extends EngineOptions, NestedStackProps {}

export class CromwellEngineStack extends NestedEngineStack {
  public readonly engine: SecureService;
  public readonly adapterRole: IRole;
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly engineLogGroup: ILogGroup;
  public readonly engineRole: IRole;

  constructor(scope: Construct, id: string, props: CromwellEngineStackProps) {
    super(scope, id, props);
    const params = props.contextParameters;
    this.engineLogGroup = new LogGroup(this, "EngineLogGroup");
    const engineContainer = params.getEngineContainer(props.jobQueue.jobQueueArn);
    const artifactBucket = Bucket.fromBucketName(this, "ArtifactBucket", params.artifactBucketName);
    const outputBucket = Bucket.fromBucketName(this, "OutputBucket", params.outputBucketName);

    this.engineRole = new CromwellEngineRole(this, "CromwellEngineRole", {
      jobQueueArn: props.jobQueue.jobQueueArn,
      readOnlyBucketArns: (params.readBucketArns ?? []).concat(artifactBucket.bucketArn),
      readWriteBucketArns: (params.readWriteBucketArns ?? []).concat(outputBucket.bucketArn),
      policies: props.policyOptions,
    });
    this.adapterRole = new CromwellAdapterRole(this, "CromwellAdapterRole", {
      readOnlyBucketArns: [],
      readWriteBucketArns: [outputBucket.bucketArn],
    });

    // TODO: Move log group creation into service construct and make it a property
    this.engine = this.getEngineServiceDefinition(props.vpc, engineContainer, this.engineLogGroup);
    this.adapterLogGroup = new LogGroup(this, "AdapterLogGroup");

    const lambda = this.renderAdapterLambda({
      vpc: props.vpc,
      role: this.adapterRole,
      engineLogGroupName: this.adapterLogGroup.logGroupName,
      jobQueueArn: props.jobQueue.jobQueueArn,
      projectName: params.projectName,
      contextName: params.contextName,
      userId: params.userId,
      engineEndpoint: this.engine.loadBalancer.loadBalancerDnsName,
    });
    this.adapterLogGroup = lambda.logGroup;

    this.apiProxy = new ApiProxy(this, {
      apiName: `${params.projectName}${params.contextName}${engineContainer.serviceName}ApiProxy`,
      lambda,
      allowedAccountIds: [this.account],
    });
  }

  protected getOutputs(): EngineOutputs {
    return {
      accessLogGroup: this.apiProxy.accessLogGroup,
      adapterLogGroup: this.adapterLogGroup,
      engineLogGroup: this.engineLogGroup,
      wesUrl: this.apiProxy.restApi.url,
    };
  }

  private getEngineServiceDefinition(vpc: IVpc, serviceContainer: ServiceContainer, logGroup: ILogGroup) {
    const id = "Engine";
    const fileSystem = new FileSystem(this, "EngineFileSystem", {
      vpc,
      encrypted: true,
      removalPolicy: RemovalPolicy.DESTROY,
    });
    const definition = new FargateTaskDefinition(this, "EngineTaskDef", {
      taskRole: this.engineRole,
      cpu: serviceContainer.cpu,
      memoryLimitMiB: serviceContainer.memoryLimitMiB,
    });

    const volumeName = "cromwell-executions";
    definition.addVolume({
      name: volumeName,
      efsVolumeConfiguration: {
        fileSystemId: fileSystem.fileSystemId,
      },
    });

    const container = definition.addContainer(serviceContainer.serviceName, {
      cpu: serviceContainer.cpu,
      memoryLimitMiB: serviceContainer.memoryLimitMiB,
      environment: serviceContainer.environment,
      containerName: serviceContainer.serviceName,
      image: createEcrImage(this, serviceContainer.imageConfig.designation),
      logging: LogDriver.awsLogs({ logGroup, streamPrefix: id }),
      portMappings: serviceContainer.containerPort ? [{ containerPort: serviceContainer.containerPort }] : [],
    });

    container.addMountPoints({
      containerPath: "/cromwell-executions",
      readOnly: false,
      sourceVolume: volumeName,
    });

    const engine = renderServiceWithTaskDefinition(this, id, serviceContainer, definition, vpc);
    fileSystem.connections.allowDefaultPortFrom(engine.service);
    return engine;
  }

  private renderAdapterLambda({ vpc, role, jobQueueArn, engineLogGroupName, projectName, contextName, userId, engineEndpoint }) {
    return renderPythonLambda(this, "CromwellWesAdapterLambda", vpc, role, wesAdapterSourcePath, {
      ENGINE_NAME: "cromwell",
      ENGINE_ENDPOINT: engineEndpoint,
      ENGINE_LOG_GROUP: engineLogGroupName,
      JOB_QUEUE: jobQueueArn,
      PROJECT_NAME: projectName,
      CONTEXT_NAME: contextName,
      USER_ID: userId,
    });
  }
}
