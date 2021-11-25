import { CfnOutput, NestedStack, NestedStackProps, Stack } from "monocdk";
import { Construct } from "constructs";
import { ILogGroup } from "monocdk/aws-logs";

export interface EngineOutputs {
  accessLogGroup: ILogGroup;
  adapterLogGroup: ILogGroup;
  engineLogGroup: ILogGroup;
  wesUrl: string;
}

export abstract class NestedEngineStack extends NestedStack {
  protected constructor(scope: Construct, id: string, props: NestedStackProps) {
    super(scope, id, props);
  }

  public outputToParent(parentStack: Stack): void {
    const outputs = this.getOutputs();
    new CfnOutput(parentStack, "AccessLogGroupName", { value: outputs.accessLogGroup.logGroupName });
    new CfnOutput(parentStack, "AdapterLogGroupName", { value: outputs.adapterLogGroup.logGroupName });
    new CfnOutput(parentStack, "EngineLogGroupName", { value: outputs.engineLogGroup.logGroupName });
    new CfnOutput(parentStack, "WesUrl", { value: outputs.wesUrl });
  }

  protected abstract getOutputs(): EngineOutputs;
}
