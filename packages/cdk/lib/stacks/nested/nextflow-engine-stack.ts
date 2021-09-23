import { NestedStackProps } from "monocdk";
import { Construct } from "constructs";
import { NextflowEngine } from "../../constructs/engines/nextflow/nextflow-engine";
import { renderServiceWithContainer } from "../../util";
import { EngineOptions } from "../../types";
import { PolicyStatement, Role, ServicePrincipal } from "monocdk/aws-iam";
import { Bucket } from "monocdk/aws-s3";
import { ApiProxy } from "../../constructs";
import { LogGroup } from "monocdk/aws-logs";
import { EngineOutputs, NestedEngineStack } from "./nested-engine-stack";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { IJobQueue } from "monocdk/aws-batch";

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
    const taskRole = new Role(this, "TaskRole", { assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"), ...props.policyOptions });

    const outputBucket = Bucket.fromBucketName(this, "OutputBucket", params.outputBucketName);
    outputBucket.grantReadWrite(taskRole);

    this.nextflowEngine = new NextflowEngine(this, "NextflowEngine", {
      vpc: props.vpc,
      outputBucketName: params.outputBucketName,
      jobQueueArn: props.jobQueue.jobQueueArn,
      rootDir: params.getEngineBucketPath(),
      taskRole,
    });

    const engineLogGroup = this.nextflowEngine.logGroup;
    const adapterContainer = params.getAdapterContainer();
    adapterContainer.environment!["JOB_DEFINITION"] = this.nextflowEngine.headJobDefinition.jobDefinitionArn;
    adapterContainer.environment!["JOB_QUEUE"] = props.headQueue.jobQueueArn;
    adapterContainer.environment!["ENGINE_LOG_GROUP"] = engineLogGroup.logGroupName;

    this.adapterLogGroup = new LogGroup(this, "AdapterLogGroup");
    const adapter = renderServiceWithContainer(this, "Adapter", adapterContainer, props.vpc, taskRole, this.adapterLogGroup);
    engineLogGroup.grant(taskRole, "logs:StartQuery");
    taskRole.addToPolicy(
      new PolicyStatement({
        actions: ["logs:GetQueryResults", "logs:StopQuery"],
        resources: ["*"],
      })
    );
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
