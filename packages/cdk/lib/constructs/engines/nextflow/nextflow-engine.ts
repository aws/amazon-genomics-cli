import { Construct } from "constructs";
import { JobDefinition } from "@aws-cdk/aws-batch-alpha";
import { IRole } from "aws-cdk-lib/aws-iam";
import { createEcrImage } from "../../../util";
import { EngineJobDefinition } from "../engine-job-definition";
import { Engine, EngineProps } from "../engine";

export interface NextflowEngineProps extends EngineProps {
  readonly jobQueueArn: string;
  readonly taskRole: IRole;
}

const NEXTFLOW_IMAGE_DESIGNATION = "nextflow";

export class NextflowEngine extends Engine {
  readonly headJobDefinition: JobDefinition;

  constructor(scope: Construct, id: string, props: NextflowEngineProps) {
    super(scope, id);

    this.headJobDefinition = new EngineJobDefinition(this, "NexflowHeadJobDef", {
      logGroup: this.logGroup,
      container: {
        jobRole: props.taskRole,
        image: createEcrImage(this, NEXTFLOW_IMAGE_DESIGNATION),
        command: [],
        environment: {
          NF_JOB_QUEUE: props.jobQueueArn,
          NF_WORKDIR: `${props.rootDirS3Uri}/runs`,
          NF_LOGSDIR: `${props.rootDirS3Uri}/logs`,
        },
        volumes: [],
      },
    });
  }
}
