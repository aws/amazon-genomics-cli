import { JobDefinition, JobDefinitionProps } from "@aws-cdk/aws-batch-alpha";
import { Construct } from "constructs";
import { renderBatchLogConfiguration } from "../../util";
import { ILogGroup } from "aws-cdk-lib/aws-logs";

interface EngineJobDefinitionProps extends JobDefinitionProps {
  readonly logGroup: ILogGroup;
}

export class EngineJobDefinition extends JobDefinition {
  constructor(scope: Construct, id: string, props: EngineJobDefinitionProps) {
    super(scope, id, {
      ...props,
      container: {
        ...props.container,
        logConfiguration: renderBatchLogConfiguration(scope, props.logGroup),
        environment: {
          AWS_METADATA_SERVICE_TIMEOUT: "10",
          AWS_METADATA_SERVICE_NUM_ATTEMPTS: "10",
          ...props.container.environment,
        },
        memoryLimitMiB: props.container.memoryLimitMiB || 2048,
        vcpus: props.container.vcpus || 1,
      },
    });
  }
}
