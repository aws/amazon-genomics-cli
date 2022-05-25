import { Aws, Stack } from "aws-cdk-lib";
import { Bucket, IBucket } from "aws-cdk-lib/aws-s3";
import { ApiProxy, Batch } from "../../constructs";
import { EngineConstruct, EngineOutputs } from "./engine-construct";
import { Effect, IRole, ManagedPolicy, PolicyDocument, PolicyStatement, Role, ServicePrincipal } from "aws-cdk-lib/aws-iam";
import { ILogGroup } from "aws-cdk-lib/aws-logs";
import { MiniWdlEngine } from "../../constructs/engines/miniwdl/miniwdl-engine";
import { IMachineImage, IVpc, SubnetSelection } from "aws-cdk-lib/aws-ec2";
import { ENGINE_MINIWDL } from "../../constants";
import { ComputeResourceType } from "@aws-cdk/aws-batch-alpha";
import { BucketOperations } from "../../common/BucketOperations";
import { ContextAppParameters } from "../../env";
import { HeadJobBatchPolicy } from "../../roles/policies/head-job-batch-policy";
import { BatchPolicies } from "../../roles/policies/batch-policies";
import { EngineOptions } from "../../types";
import { Construct } from "constructs";
import { LaunchTemplateData } from "../../constructs/launch-template-data";

export class MiniwdlEngineConstruct extends EngineConstruct {
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly miniwdlEngine: MiniWdlEngine;
  private readonly batchHead: Batch;
  private readonly batchWorkers: Batch;
  private readonly outputBucket: IBucket;

  constructor(scope: Construct, id: string, props: EngineOptions) {
    super(scope, id);

    const { vpc, contextParameters, subnets, computeEnvImage } = props;
    const params = props.contextParameters;
    const rootDirS3Uri = params.getEngineBucketPath();

    this.batchHead = this.renderBatch("HeadBatch", vpc, subnets, contextParameters, ComputeResourceType.FARGATE);
    const workerComputeType = contextParameters.requestSpotInstances ? ComputeResourceType.SPOT : ComputeResourceType.ON_DEMAND;
    this.batchWorkers = this.renderBatch("TaskBatch", vpc, subnets, contextParameters, workerComputeType, computeEnvImage);

    this.batchHead.role.attachInlinePolicy(new HeadJobBatchPolicy(this, "HeadJobBatchPolicy"));
    this.batchHead.role.addToPrincipalPolicy(
      new PolicyStatement({
        actions: ["batch:TagResource"],
        resources: ["*"],
      })
    );
    this.batchHead.role.addToPrincipalPolicy(
      new PolicyStatement({
        effect: Effect.ALLOW,
        actions: ["batch:TerminateJob"],
        resources: ["*"],
        conditions: { "ForAllValues:StringEquals": { "aws:TagKeys": ["AWS_BATCH_PARENT_JOB_ID"] } },
      })
    );

    this.miniwdlEngine = new MiniWdlEngine(this, "MiniWdlEngine", {
      vpc: props.vpc,
      subnets: props.subnets,
      iops: props.iops,
      rootDirS3Uri: rootDirS3Uri,
      engineBatch: this.batchHead,
      workerBatch: this.batchWorkers,
    });

    const adapterRole = new Role(this, "MiniWdlAdapterRole", {
      assumedBy: new ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaVPCAccessExecutionRole")],
      inlinePolicies: {
        MiniwdlAdapterPolicy: new PolicyDocument({
          statements: [
            BatchPolicies.listAndDescribe,
            new PolicyStatement({
              actions: ["tag:GetResources"],
              resources: ["*"],
            }),
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["batch:TerminateJob"],
              resources: ["*"],
            }),
          ],
        }),
      },
    });
    this.outputBucket = Bucket.fromBucketName(this, "OutputBucket", params.outputBucketName);
    this.outputBucket.grantRead(adapterRole);

    this.batchHead.grantJobAdministration(adapterRole);
    this.batchWorkers.grantJobAdministration(this.batchHead.role);

    this.grantS3Permissions(contextParameters);

    const lambda = this.renderAdapterLambda({
      role: adapterRole,
      jobQueueArn: this.batchHead.jobQueue.jobQueueArn,
      jobDefinitionArn: this.miniwdlEngine.headJobDefinition.jobDefinitionArn,
      rootDirS3Uri: rootDirS3Uri,
      vpc: props.contextParameters.usePublicSubnets ? undefined : props.vpc,
      vcpSubnets: props.contextParameters.usePublicSubnets ? undefined : props.subnets,
    });
    this.adapterLogGroup = lambda.logGroup;

    this.apiProxy = new ApiProxy(this, {
      apiName: `${params.projectName}${params.userId}${params.contextName}MiniWdlApiProxy`,
      lambda,
      allowedAccountIds: [Aws.ACCOUNT_ID],
    });
  }

  protected getOutputs(): EngineOutputs {
    return {
      accessLogGroup: this.apiProxy.accessLogGroup,
      adapterLogGroup: this.adapterLogGroup,
      engineLogGroup: this.miniwdlEngine.logGroup,
      wesUrl: this.apiProxy.restApi.url,
    };
  }

  private grantS3Permissions(contextParameters: ContextAppParameters) {
    const { artifactBucketName, readBucketArns = [], readWriteBucketArns = [] } = contextParameters;

    const artifactBucket = Bucket.fromBucketName(this, "ArtifactBucket", artifactBucketName);

    readBucketArns.push(artifactBucket.bucketArn);
    readWriteBucketArns.push(this.outputBucket.bucketArn);

    const batchRoles = this.getBatchRoles();
    for (const role of batchRoles) {
      BucketOperations.grantBucketAccess(this, role, readBucketArns, true);
      BucketOperations.grantBucketAccess(this, role, readWriteBucketArns, false);
    }
  }

  private renderBatch(
    id: string,
    vpc: IVpc,
    subnets: SubnetSelection,
    appParams: ContextAppParameters,
    computeType?: ComputeResourceType,
    computeEnvImage?: IMachineImage
  ): Batch {
    return new Batch(this, id, {
      vpc,
      subnets,
      computeType,
      instanceTypes: appParams.instanceTypes,
      maxVCpus: appParams.maxVCpus,
      launchTemplateData: LaunchTemplateData.renderLaunchTemplateData(ENGINE_MINIWDL),
      awsPolicyNames: ["AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"],
      resourceTags: Stack.of(this).tags.tagValues(),
      workflowOrchestrator: ENGINE_MINIWDL,
      computeEnvImage,
    });
  }

  private getBatchRoles(): IRole[] {
    return [this.batchHead.role, this.batchWorkers.role];
  }

  private renderAdapterLambda({ role, jobQueueArn, jobDefinitionArn, rootDirS3Uri, vpc, vcpSubnets }) {
    return super.renderPythonLambda(
      this,
      "MiniWDLWesAdapterLambda",
      role,
      {
        ENGINE_NAME: ENGINE_MINIWDL,
        JOB_QUEUE: jobQueueArn,
        JOB_DEFINITION: jobDefinitionArn,
        OUTPUT_DIR_S3_URI: rootDirS3Uri,
      },
      vpc,
      vcpSubnets
    );
  }
}
