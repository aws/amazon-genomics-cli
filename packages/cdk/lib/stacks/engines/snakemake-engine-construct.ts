import { Aws, Stack } from "aws-cdk-lib";
import { SnakemakeEngine } from "../../constructs/engines/snakemake/snakemake-engine";
import { EngineOptions } from "../../types";
import { ApiProxy, Batch } from "../../constructs";
import { EngineOutputs, EngineConstruct } from "./engine-construct";
import { ILogGroup } from "aws-cdk-lib/aws-logs";
import { ComputeResourceType } from "@aws-cdk/aws-batch-alpha";
import { ENGINE_SNAKEMAKE } from "../../constants";
import { Construct } from "constructs";
import { IRole, ManagedPolicy, PolicyDocument, PolicyStatement, Role, ServicePrincipal } from "aws-cdk-lib/aws-iam";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { ContextAppParameters } from "../../env";
import { HeadJobBatchPolicy } from "../../roles/policies/head-job-batch-policy";
import { BatchPolicies } from "../../roles/policies/batch-policies";
import { Bucket, IBucket } from "aws-cdk-lib/aws-s3";
import { BucketOperations } from "../../common/BucketOperations";
import { LaunchTemplateData } from "../../constructs/launch-template-data";

export class SnakemakeEngineConstruct extends EngineConstruct {
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly snakemakeEngine: SnakemakeEngine;
  private readonly batchHead: Batch;
  private readonly batchWorkers: Batch;
  private readonly outputBucket: IBucket;

  constructor(scope: Construct, id: string, props: EngineOptions) {
    super(scope, id);

    const { vpc, contextParameters } = props;
    const params = props.contextParameters;

    this.batchHead = this.renderBatch("HeadBatch", vpc, contextParameters, ComputeResourceType.FARGATE);
    const workerComputeType = contextParameters.requestSpotInstances ? ComputeResourceType.SPOT : ComputeResourceType.ON_DEMAND;
    this.batchWorkers = this.renderBatch("TaskBatch", vpc, contextParameters, workerComputeType);

    this.snakemakeEngine = new SnakemakeEngine(this, "SnakemakeEngine", {
      vpc: props.vpc,
      engineBatch: this.batchHead,
      workerBatch: this.batchWorkers,
      rootDirS3Uri: params.getEngineBucketPath(),
    });

    this.batchHead.role.addToPrincipalPolicy(
      new PolicyStatement({
        actions: ["elasticfilesystem:DescribeAccessPoints"],
        resources: [this.snakemakeEngine.fileSystem.fileSystemArn, this.snakemakeEngine.fsap.accessPointArn],
      })
    );
    this.batchHead.role.attachInlinePolicy(new HeadJobBatchPolicy(this, "HeadJobBatchPolicy"));
    this.batchHead.role.addToPrincipalPolicy(
      new PolicyStatement({
        actions: ["batch:TagResource"],
        resources: ["*"],
      })
    );

    this.batchHead.role.addToPrincipalPolicy(
      new PolicyStatement({
        actions: ["iam:PassRole"],
        resources: [this.batchHead.role.roleArn],
      })
    );

    const adapterRole = new Role(this, "SnakemakeAdapterRole", {
      assumedBy: new ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaVPCAccessExecutionRole")],
      inlinePolicies: {
        SnakemakeAdapterPolicy: new PolicyDocument({
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

    this.grantS3Permissions(contextParameters, [this.batchHead.role, this.batchWorkers.role]);

    const lambda = this.renderAdapterLambda({
      vpc: props.vpc,
      role: adapterRole,
      jobQueueArn: this.batchHead.jobQueue.jobQueueArn,
      jobDefinitionArn: this.snakemakeEngine.headJobDefinition.jobDefinitionArn,
      workflowRoleArn: this.batchHead.role.roleArn,
      taskQueueArn: this.batchWorkers.jobQueue.jobQueueArn,
      fsapId: this.snakemakeEngine.fsap.accessPointId,
      outputBucket: params.getEngineBucketPath(),
    });
    this.adapterLogGroup = lambda.logGroup;

    this.apiProxy = new ApiProxy(this, {
      apiName: `${params.projectName}${params.userId}${params.contextName}SnakemakeApiProxy`,
      lambda,
      allowedAccountIds: [Aws.ACCOUNT_ID],
    });
  }

  protected getOutputs(): EngineOutputs {
    return {
      accessLogGroup: this.apiProxy.accessLogGroup,
      adapterLogGroup: this.adapterLogGroup,
      engineLogGroup: this.snakemakeEngine.logGroup,
      wesUrl: this.apiProxy.restApi.url,
    };
  }

  private grantS3Permissions(contextParameters: ContextAppParameters, batchRoles: IRole[]) {
    const { artifactBucketName, readBucketArns = [], readWriteBucketArns = [] } = contextParameters;

    const artifactBucket = Bucket.fromBucketName(this, "ArtifactBucket", artifactBucketName);

    readBucketArns.push(artifactBucket.bucketArn);
    readWriteBucketArns.push(this.outputBucket.bucketArn);

    batchRoles.forEach((role) => {
      BucketOperations.grantBucketAccess(this, role, readBucketArns, true);
      BucketOperations.grantBucketAccess(this, role, readWriteBucketArns, false);
    });
  }

  private renderBatch(id: string, vpc: IVpc, appParams: ContextAppParameters, computeType?: ComputeResourceType): Batch {
    return new Batch(this, id, {
      vpc,
      computeType,
      instanceTypes: appParams.instanceTypes,
      maxVCpus: appParams.maxVCpus,
      launchTemplateData: LaunchTemplateData.renderLaunchTemplateData(ENGINE_SNAKEMAKE),
      awsPolicyNames: ["AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"],
      resourceTags: Stack.of(this).tags.tagValues(),
    });
  }

  private renderAdapterLambda({ vpc, role, jobQueueArn, jobDefinitionArn, taskQueueArn, workflowRoleArn, fsapId, outputBucket }) {
    return super.renderPythonLambda(this, "SnakemakeWesAdapterLambda", vpc, role, {
      ENGINE_NAME: "snakemake",
      JOB_QUEUE: jobQueueArn,
      JOB_DEFINITION: jobDefinitionArn,
      TASK_QUEUE: taskQueueArn,
      WORKFLOW_ROLE: workflowRoleArn,
      FSAP_ID: fsapId,
      OUTPUT_DIR_S3_URI: outputBucket,
      TIME: Date.now().toString(),
    });
  }
}
