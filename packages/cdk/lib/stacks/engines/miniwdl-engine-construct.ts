import { Stack, Aws } from "aws-cdk-lib";
import { Bucket, IBucket } from "aws-cdk-lib/aws-s3";
import { ApiProxy, Batch } from "../../constructs";
import { EngineOutputs, EngineConstruct } from "./engine-construct";
import { IRole, PolicyDocument, PolicyStatement, Role, ServicePrincipal, ManagedPolicy } from "aws-cdk-lib/aws-iam";
import { ILogGroup } from "aws-cdk-lib/aws-logs";
import { MiniWdlEngine } from "../../constructs/engines/miniwdl/miniwdl-engine";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { LAUNCH_TEMPLATE } from "../../constants";
import { ComputeResourceType } from "@aws-cdk/aws-batch-alpha";
import { BucketOperations } from "../../common/BucketOperations";
import { ContextAppParameters } from "../../env";
import { HeadJobBatchPolicy } from "../../roles/policies/head-job-batch-policy";
import { BatchPolicies } from "../../roles/policies/batch-policies";
import { EngineOptions } from "../../types";
import { Construct } from "constructs";

export class MiniwdlEngineConstruct extends EngineConstruct {
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly miniwdlEngine: MiniWdlEngine;
  private readonly batchHead: Batch;
  private readonly batchWorkers: Batch;
  private readonly outputBucket: IBucket;

  constructor(scope: Construct, id: string, props: EngineOptions) {
    super(scope, id);

    const { vpc, contextParameters } = props;
    const params = props.contextParameters;
    const rootDirS3Uri = params.getEngineBucketPath();

    this.batchHead = this.renderBatch("HeadBatch", vpc, contextParameters, ComputeResourceType.FARGATE);
    const workerComputeType = contextParameters.requestSpotInstances ? ComputeResourceType.SPOT : ComputeResourceType.ON_DEMAND;
    this.batchWorkers = this.renderBatch("TaskBatch", vpc, contextParameters, workerComputeType);

    this.batchHead.role.attachInlinePolicy(new HeadJobBatchPolicy(this, "HeadJobBatchPolicy"));
    this.batchHead.role.addToPrincipalPolicy(
      new PolicyStatement({
        actions: ["batch:TagResource"],
        resources: ["*"],
      })
    );

    this.miniwdlEngine = new MiniWdlEngine(this, "MiniWdlEngine", {
      vpc: props.vpc,
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
      vpc: props.vpc,
      role: adapterRole,
      jobQueueArn: this.batchHead.jobQueue.jobQueueArn,
      jobDefinitionArn: this.miniwdlEngine.headJobDefinition.jobDefinitionArn,
      rootDirS3Uri: rootDirS3Uri,
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

  private renderBatch(id: string, vpc: IVpc, appParams: ContextAppParameters, computeType?: ComputeResourceType): Batch {
    return new Batch(this, id, {
      vpc,
      computeType,
      instanceTypes: appParams.instanceTypes,
      maxVCpus: appParams.maxVCpus,
      launchTemplateData: LAUNCH_TEMPLATE,
      awsPolicyNames: ["AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"],
      resourceTags: Stack.of(this).tags.tagValues(),
    });
  }

  private getBatchRoles(): IRole[] {
    return [this.batchHead.role, this.batchWorkers.role];
  }

  private renderAdapterLambda({ vpc, role, jobQueueArn, jobDefinitionArn, rootDirS3Uri }) {
    return super.renderPythonLambda(this, "MiniWDLWesAdapterLambda", vpc, role, {
      ENGINE_NAME: "miniwdl",
      JOB_QUEUE: jobQueueArn,
      JOB_DEFINITION: jobDefinitionArn,
      OUTPUT_DIR_S3_URI: rootDirS3Uri,
    });
  }
}
