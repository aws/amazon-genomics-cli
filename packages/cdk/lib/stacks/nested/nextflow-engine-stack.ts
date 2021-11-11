import { NestedStackProps } from "monocdk";
import { Construct } from "constructs";
import { NextflowEngine } from "../../constructs/engines/nextflow/nextflow-engine";
import { renderPythonLambda } from "../../util";
import { EngineOptions } from "../../types";
import { Bucket } from "monocdk/aws-s3";
import { ApiProxy } from "../../constructs";
import { EngineOutputs, NestedEngineStack } from "./nested-engine-stack";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { IJobQueue } from "monocdk/aws-batch";
import { NextflowEngineRole } from "../../roles/nextflow-engine-role";
import { NextflowAdapterRole } from "../../roles/nextflow-adapter-role";
import { wesAdapterSourcePath } from "../../constants";

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
      batchJobPolicyArns: [props.jobQueue.jobQueueArn],
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
      batchJobPolicyArns: [this.nextflowEngine.headJobDefinition.jobDefinitionArn, props.headQueue.jobQueueArn],
      readOnlyBucketArns: [],
      readWriteBucketArns: [outputBucket.bucketArn],
    });

    const engineLogGroup = this.nextflowEngine.logGroup;
    engineLogGroup.grant(engineRole, "logs:StartQuery");
    engineLogGroup.grant(adapterRole, "logs:StartQuery");

    const lambda = this.renderAdapterLambda({
      vpc: props.vpc,
      role: adapterRole,
      jobQueueArn: props.headQueue.jobQueueArn,
      jobDefinitionArn: this.nextflowEngine.headJobDefinition.jobDefinitionArn,
      engineLogGroupName: engineLogGroup.logGroupName,
    });
    this.adapterLogGroup = lambda.logGroup;

    this.apiProxy = new ApiProxy(this, {
      apiName: `${params.projectName}${params.userId}${params.contextName}NextflowApiProxy`,
      lambda,
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

  private renderAdapterLambda({ vpc, role, jobQueueArn, jobDefinitionArn, engineLogGroupName }) {
    return renderPythonLambda(this, "NextflowWesAdapterLambda", vpc, role, wesAdapterSourcePath, {
      ENGINE_NAME: "nextflow",
      JOB_QUEUE: jobQueueArn,
      JOB_DEFINITION: jobDefinitionArn,
      ENGINE_LOG_GROUP: engineLogGroupName,
    });
  }
}
