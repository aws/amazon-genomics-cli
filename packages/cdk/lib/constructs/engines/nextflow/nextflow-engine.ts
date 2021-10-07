import { Construct } from "monocdk";
import { CfnJobDefinition, JobDefinition } from "monocdk/aws-batch";
import { IVpc } from "monocdk/aws-ec2";
import { IRole } from "monocdk/aws-iam";
import { createEcrImage, renderBatchLogConfiguration } from "../../../util";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { LogGroup } from "monocdk/aws-logs";

export interface NextflowEngineProps {
  readonly vpc: IVpc;
  readonly outputBucketName: string;
  readonly jobQueueArn: string;
  readonly taskRole: IRole;
  readonly rootDir: string;
}

const NEXTFLOW_IMAGE_DESIGNATION = "nextflow";

export class NextflowEngine extends Construct {
  readonly headJobDefinition: JobDefinition;
  readonly logGroup: ILogGroup;

  constructor(scope: Construct, id: string, props: NextflowEngineProps) {
    super(scope, id);

    this.logGroup = new LogGroup(this, "EngineLogGroup");
    this.headJobDefinition = new JobDefinition(this, "NexflowHeadJobDef", {
      container: {
        logConfiguration: renderBatchLogConfiguration(this, this.logGroup),
        jobRole: props.taskRole,
        image: createEcrImage(this, NEXTFLOW_IMAGE_DESIGNATION),
        command: [],
        environment: {
          NF_JOB_QUEUE: props.jobQueueArn,
          NF_WORKDIR: `${props.rootDir}/runs`,
          NF_LOGSDIR: `${props.rootDir}/logs`,
          AWS_METADATA_SERVICE_TIMEOUT: "10",
          AWS_METADATA_SERVICE_NUM_ATTEMPTS: "10",
        },
        volumes: [],
      },
    });

    const cfnJobDef = this.headJobDefinition.node.defaultChild as CfnJobDefinition;

    //Removing old method for specifying resources. Using newer ResourceRequirements.
    cfnJobDef.addPropertyDeletionOverride("ContainerProperties.Vcpus");
    cfnJobDef.addPropertyDeletionOverride("ContainerProperties.Memory");
    cfnJobDef.addPropertyOverride("ContainerProperties.ResourceRequirements", [
      { Type: "VCPU", Value: "1" },
      { Type: "MEMORY", Value: "2048" },
    ]);
  }
}
