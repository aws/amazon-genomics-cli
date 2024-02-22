import { EcsEc2ContainerDefinition, EcsJobDefinition, EcsJobDefinitionProps } from "aws-cdk-lib/aws-batch";
import { Construct } from "constructs";
import { ILogGroup } from "aws-cdk-lib/aws-logs";
import { Size } from "aws-cdk-lib";

interface EngineJobDefinitionProps extends EcsJobDefinitionProps {
  readonly logGroup: ILogGroup;
}

export class EngineJobDefinition extends EcsJobDefinition {
  constructor(scope: Construct, id: string, props: EngineJobDefinitionProps) {
    super(scope, id, {
      ...props,
      container: new EcsEc2ContainerDefinition(scope, "containerDefn", {
        ...props.container,
        cpu: props.container.cpu || 1,
        memory: props.container.memory || Size.mebibytes(2048),
        environment: {
          AWS_METADATA_SERVICE_TIMEOUT: "10",
          AWS_METADATA_SERVICE_NUM_ATTEMPTS: "10",
          ...props.container.environment,
        },
      }),
    });
  }
}
