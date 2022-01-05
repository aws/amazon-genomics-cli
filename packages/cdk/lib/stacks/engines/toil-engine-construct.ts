import { Aws } from "aws-cdk-lib";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { FargateTaskDefinition, LogDriver } from "aws-cdk-lib/aws-ecs";
import { ApiProxy, SecureService } from "../../constructs";
import { IRole } from "aws-cdk-lib/aws-iam";
import { createEcrImage, renderServiceWithTaskDefinition } from "../../util";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { EngineOptions, ServiceContainer } from "../../types";
import { LogGroup, ILogGroup } from "aws-cdk-lib/aws-logs";
import { EngineOutputs, EngineConstruct } from "./engine-construct";
import { ToilJobRole } from "../../roles/toil-job-role";
import { ToilEngineRole } from "../../roles/toil-engine-role";
import { IJobQueue } from "@aws-cdk/aws-batch-alpha";
import { Construct } from "constructs";

export interface ToilEngineConstructProps extends EngineOptions {
  /**
   * AWS Batch JobQueue to use for running workflows.
   */
  readonly jobQueue: IJobQueue;
}

export class ToilEngineConstruct extends EngineConstruct {
  public readonly engine: SecureService;
  public readonly adapterRole: IRole;
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly engineLogGroup: ILogGroup;
  public readonly engineRole: IRole;
  public readonly jobRole: IRole;

  constructor(scope: Construct, id: string, props: ToilEngineConstructProps) {
    super(scope, id);
    const params = props.contextParameters;
    this.engineLogGroup = new LogGroup(this, "EngineLogGroup");
    const artifactBucket = Bucket.fromBucketName(this, "ArtifactBucket", params.artifactBucketName);
    const outputBucket = Bucket.fromBucketName(this, "OutputBucket", params.outputBucketName);

    // Make a role for the jobs to run as
    this.jobRole = new ToilJobRole(this, "ToilJobRole", {
      readOnlyBucketArns: (params.readBucketArns ?? []).concat(artifactBucket.bucketArn),
      readWriteBucketArns: (params.readWriteBucketArns ?? []).concat(outputBucket.bucketArn),
      policies: props.policyOptions,
    });

    // Make a role for the server to run as
    this.engineRole = new ToilEngineRole(this, "ToilEngineRole", {
      jobQueueArn: props.jobQueue.jobQueueArn,
      jobRoleArn: this.jobRole.roleArn,
      readOnlyBucketArns: (params.readBucketArns ?? []).concat(artifactBucket.bucketArn),
      readWriteBucketArns: (params.readWriteBucketArns ?? []).concat(outputBucket.bucketArn),
      policies: props.policyOptions,
    });

    // Make the container and pass it the ARN of the role to use for individual jobs.
    const engineContainer = params.getEngineContainer(props.jobQueue.jobQueueArn, {
      TOIL_AWS_BATCH_JOB_ROLE_ARN: this.jobRole.roleArn,
    });

    // TODO: Move log group creation into service construct and make it a property
    this.engine = this.getEngineServiceDefinition(props.vpc, engineContainer, this.engineLogGroup);
    // This is unused because we have no adapter, but a log group is required.
    this.adapterLogGroup = new LogGroup(this, "AdapterLogGroup");

    // We don't use an adapter, so put the access-controlling proxy right in
    // front of the engine load balancer.
    this.apiProxy = new ApiProxy(this, {
      apiName: `${params.projectName}${params.contextName}${engineContainer.serviceName}ApiProxy`,
      loadBalancer: this.engine.loadBalancer,
      allowedAccountIds: [Aws.ACCOUNT_ID],
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
    const definition = new FargateTaskDefinition(this, "EngineTaskDef", {
      taskRole: this.engineRole,
      cpu: serviceContainer.cpu,
      memoryLimitMiB: serviceContainer.memoryLimitMiB,
    });

    definition.addContainer(serviceContainer.serviceName, {
      cpu: serviceContainer.cpu,
      memoryLimitMiB: serviceContainer.memoryLimitMiB,
      environment: serviceContainer.environment,
      containerName: serviceContainer.serviceName,
      image: createEcrImage(this, serviceContainer.imageConfig.designation),
      logging: LogDriver.awsLogs({ logGroup, streamPrefix: id }),
      portMappings: serviceContainer.containerPort ? [{ containerPort: serviceContainer.containerPort }] : [],
    });

    const engine = renderServiceWithTaskDefinition(this, id, serviceContainer, definition, vpc);
    return engine;
  }
}
