import { NestedStackProps } from "monocdk";
import { Construct } from "constructs";
import { NextflowEngine } from "../../constructs/engines/nextflow/nextflow-engine";
import { renderServiceWithContainer } from "../../util";
import { EngineOptions } from "../../types";
import { Bucket } from "monocdk/aws-s3";
import { ApiProxy } from "../../constructs";
import { LogGroup } from "monocdk/aws-logs";
import { EngineOutputs, NestedEngineStack } from "./nested-engine-stack";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { IJobQueue } from "monocdk/aws-batch";
import { NextflowEngineRole } from "../../roles/nextflow-engine-role";
import { NextflowAdapterRole } from "../../roles/nextflow-adapter-role";

export interface NextflowEngineStackProps extends EngineOptions, NestedStackProps {
  /**
   * AWS Batch JobQueue to use for running workflows.
   */
  readonly headQueue: IJobQueue;
}

export class NextflowEngineStack extends NestedEngineStack {
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly nextflowEngine: NextflowEngine;

  constructor(scope: Construct, id: string, props: NextflowEngineStackProps) {
    super(scope, id, props);

    const params = props.contextParameters;
    const outputBucket = Bucket.fromBucketName(this, "OutputBucket", params.outputBucketName);
    const artifactBucket = Bucket.fromBucketName(this, "ArtifactBucket", params.artifactBucketName);

    const engineRole = new NextflowEngineRole(this, "NextflowEngineRole", {
      readOnlyBucketArns: (params.readBucketArns ?? []).concat(artifactBucket.bucketArn),
      readWriteBucketArns: (params.readWriteBucketArns ?? []).concat(outputBucket.bucketArn),
      policies: props.policyOptions,
    });

    this.nextflowEngine = new NextflowEngine(this, "NextflowEngine", {
      vpc: props.vpc,
      outputBucketName: params.outputBucketName,
      jobQueueArn: props.jobQueue.jobQueueArn,
      rootDir: params.getEngineBucketPath(),
      taskRole: engineRole,
    });

    const adapterRole = new NextflowAdapterRole(this, "NextflowAdapterRole", {
      headJobDefinitionArn: this.nextflowEngine.headJobDefinition.jobDefinitionArn,
      jobQueueArn: props.jobQueue.jobQueueArn,
      readOnlyBucketArns: [],
      readWriteBucketArns: [outputBucket.bucketArn],
    });

    const engineLogGroup = this.nextflowEngine.logGroup;
    const adapterContainer = params.getAdapterContainer();
    adapterContainer.environment!["JOB_DEFINITION"] = this.nextflowEngine.headJobDefinition.jobDefinitionArn;
    adapterContainer.environment!["JOB_QUEUE"] = props.headQueue.jobQueueArn;
    adapterContainer.environment!["ENGINE_LOG_GROUP"] = engineLogGroup.logGroupName;

    this.adapterLogGroup = new LogGroup(this, "AdapterLogGroup");
    const adapter = renderServiceWithContainer(this, "Adapter", adapterContainer, props.vpc, adapterRole, this.adapterLogGroup);
    engineLogGroup.grant(engineRole, "logs:StartQuery");

    this.apiProxy = new ApiProxy(this, {
      apiName: `${params.projectName}${params.userId}${params.contextName}NextflowApiProxy`,
      loadBalancer: adapter.loadBalancer,
      allowedAccountIds: [this.account],
    });
  }

  protected getOutputs(): EngineOutputs {
    return {
      accessLogGroup: this.apiProxy.accessLogGroup,
      adapterLogGroup: this.adapterLogGroup,
      engineLogGroup: this.nextflowEngine.logGroup,
      wesUrl: this.apiProxy.restApi.url,
    };
  }
}
