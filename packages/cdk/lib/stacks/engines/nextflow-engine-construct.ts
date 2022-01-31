import { Aws } from "aws-cdk-lib";
import { NextflowEngine } from "../../constructs/engines/nextflow/nextflow-engine";
import { EngineOptions } from "../../types";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { ApiProxy } from "../../constructs";
import { EngineOutputs, EngineConstruct } from "./engine-construct";
import { ILogGroup } from "aws-cdk-lib/aws-logs";
import { IJobQueue } from "@aws-cdk/aws-batch-alpha";
import { NextflowEngineRole } from "../../roles/nextflow-engine-role";
import { NextflowAdapterRole } from "../../roles/nextflow-adapter-role";
import { Construct } from "constructs";

export interface NextflowEngineConstructProps extends EngineOptions {
  /**
   * AWS Batch JobQueue to use for running workflows.
   */
  readonly jobQueue: IJobQueue;
  /**
   * AWS Batch JobQueue to use for running workflows.
   */
  readonly headQueue: IJobQueue;
}

export class NextflowEngineConstruct extends EngineConstruct {
  public readonly apiProxy: ApiProxy;
  public readonly adapterLogGroup: ILogGroup;
  public readonly nextflowEngine: NextflowEngine;

  constructor(scope: Construct, id: string, props: NextflowEngineConstructProps) {
    super(scope, id);

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
      jobQueueArn: props.jobQueue.jobQueueArn,
      rootDirS3Uri: params.getEngineBucketPath(),
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
      allowedAccountIds: [Aws.ACCOUNT_ID],
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
    return super.renderPythonLambda(this, "NextflowWesAdapterLambda", vpc, role, {
      ENGINE_NAME: "nextflow",
      JOB_QUEUE: jobQueueArn,
      JOB_DEFINITION: jobDefinitionArn,
      ENGINE_LOG_GROUP: engineLogGroupName,
    });
  }
}
