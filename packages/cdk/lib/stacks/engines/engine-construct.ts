import { CfnOutput, Stack, Construct } from "monocdk";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";

export interface EngineOutputs {
  accessLogGroup: ILogGroup;
  adapterLogGroup: ILogGroup;
  engineLogGroup: ILogGroup;
  wesUrl: string;
}

export abstract class EngineConstruct extends Construct {
  protected constructor(scope: Construct, id: string) {
    super(scope, id);
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
