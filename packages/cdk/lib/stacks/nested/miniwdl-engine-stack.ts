import { NestedStackProps } from "monocdk";
import { Construct } from "constructs";
import { IRole, PolicyDocument, PolicyStatement, Role, ServicePrincipal } from "monocdk/aws-iam";
import { Bucket } from "monocdk/aws-s3";
import { ApiProxy, Batch } from "../../constructs";
import { LogGroup } from "monocdk/aws-logs";
import { EngineOutputs, NestedEngineStack } from "./nested-engine-stack";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { MiniWdlEngine } from "../../constructs/engines/miniwdl/miniwdl-engine";
import { InstanceType, IVpc } from "monocdk/aws-ec2";
import { LAUNCH_TEMPLATE } from "../../constants";
import { ComputeResourceType } from "monocdk/aws-batch";
import { BucketOperations } from "../../../common/BucketOperations";
import { ContextAppParameters } from "../../env";
import { HeadJobBatchPolicy } from "../../roles/policies/head-job-batch-policy";
import { renderServiceWithContainer } from "../../util";
import { BatchPolicies } from "../../roles/policies/batch-policies";

export interface MiniWdlEngineStackProps extends NestedStackProps {
  /**
   * VPC to run resources in.
   */
  readonly vpc: IVpc;
  /**
   * Parameters determined by the context.
   */
  readonly contextParameters: ContextAppParameters;
}

export class MiniWdlEngineStack extends NestedEngineStack {
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly miniwdlEngine: MiniWdlEngine;
  private readonly batchHead: Batch;
  private readonly batchWorkers: Batch;

  constructor(scope: Construct, id: string, props: MiniWdlEngineStackProps) {
    super(scope, id, props);

    const { vpc, contextParameters } = props;
    const params = props.contextParameters;

    this.batchHead = this.renderBatch("HeadBatch", vpc, contextParameters.instanceTypes, ComputeResourceType.FARGATE);
    const workerComputeType = contextParameters.requestSpotInstances ? ComputeResourceType.SPOT : ComputeResourceType.ON_DEMAND;
    this.batchWorkers = this.renderBatch("TaskBatch", vpc, contextParameters.instanceTypes, workerComputeType);

    this.batchHead.role.attachInlinePolicy(new HeadJobBatchPolicy(this, "HeadJobBatchPolicy"));
    this.batchHead.role.addToPrincipalPolicy(
      new PolicyStatement({
        actions: ["batch:TagResource"],
        resources: ["*"],
      })
    );

    this.miniwdlEngine = new MiniWdlEngine(this, "MiniWdlEngine", {
      vpc: props.vpc,
      outputBucketName: params.outputBucketName,
      engineBatch: this.batchHead,
      workerBatch: this.batchWorkers,
    });

    const adapterRole = new Role(this, "MiniWdlAdapterRole", {
      assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
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

    this.batchHead.grantJobAdministration(adapterRole);
    this.batchWorkers.grantJobAdministration(this.batchHead.role);

    this.grantS3Permissions(contextParameters);

    const adapterContainer = params.getAdapterContainer();
    adapterContainer.environment!["JOB_DEFINITION"] = this.miniwdlEngine.headJobDefinition.jobDefinitionArn;
    adapterContainer.environment!["JOB_QUEUE"] = this.batchHead.jobQueue.jobQueueArn;
    this.adapterLogGroup = new LogGroup(this, "AdapterLogGroup");
    const adapter = renderServiceWithContainer(this, "Adapter", adapterContainer, props.vpc, adapterRole, this.adapterLogGroup);

    this.apiProxy = new ApiProxy(this, {
      apiName: `${params.projectName}${params.userId}${params.contextName}MiniWdlApiProxy`,
      loadBalancer: adapter.loadBalancer,
      allowedAccountIds: [this.account],
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
    const { artifactBucketName, outputBucketName, readBucketArns = [], readWriteBucketArns = [] } = contextParameters;

    const outputBucket = Bucket.fromBucketName(this, "OutputBucket", outputBucketName);
    const artifactBucket = Bucket.fromBucketName(this, "ArtifactBucket", artifactBucketName);

    readBucketArns.push(artifactBucket.bucketArn);
    readWriteBucketArns.push(outputBucket.bucketArn);

    const batchRoles = this.getBatchRoles();
    for (const role of batchRoles) {
      BucketOperations.grantBucketAccess(this, role, readBucketArns, true);
      BucketOperations.grantBucketAccess(this, role, readWriteBucketArns, false);
    }
  }

  private renderBatch(id: string, vpc: IVpc, instanceTypes?: InstanceType[], computeType?: ComputeResourceType): Batch {
    return new Batch(this, id, {
      vpc,
      instanceTypes,
      computeType,
      launchTemplateData: LAUNCH_TEMPLATE,
      awsPolicyNames: ["AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"],
      resourceTags: this.nestedStackParent?.tags.tagValues(),
    });
  }

  private getBatchRoles(): IRole[] {
    return [this.batchHead.role, this.batchWorkers.role];
  }
}
