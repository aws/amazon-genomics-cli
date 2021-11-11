import { CfnJobDefinition, JobDefinition, JobDefinitionProps } from "monocdk/aws-batch";
import { Construct } from "monocdk";
import { renderBatchLogConfiguration } from "../../util";
import { ILogGroup } from "monocdk/aws-logs";

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
      },
    });

    const cfnJobDef = this.node.defaultChild as CfnJobDefinition;

    //Removing old method for specifying resources. Using newer ResourceRequirements.
    cfnJobDef.addPropertyDeletionOverride("ContainerProperties.Vcpus");
    cfnJobDef.addPropertyDeletionOverride("ContainerProperties.Memory");
    cfnJobDef.addPropertyOverride("ContainerProperties.ResourceRequirements", [
      { Type: "VCPU", Value: props.container.vcpus || "1" },
      { Type: "MEMORY", Value: props.container.memoryLimitMiB || "2048" },
    ]);
  }
}
