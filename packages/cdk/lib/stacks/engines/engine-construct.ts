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

  public outputToParent(): void {
    const outputs = this.getOutputs();
    new CfnOutput(Stack.of(this), "AccessLogGroupName", { value: outputs.accessLogGroup.logGroupName });
    new CfnOutput(Stack.of(this), "AdapterLogGroupName", { value: outputs.adapterLogGroup.logGroupName });
    new CfnOutput(Stack.of(this), "EngineLogGroupName", { value: outputs.engineLogGroup.logGroupName });
    new CfnOutput(Stack.of(this), "WesUrl", { value: outputs.wesUrl });
  }

  protected abstract getOutputs(): EngineOutputs;
}
