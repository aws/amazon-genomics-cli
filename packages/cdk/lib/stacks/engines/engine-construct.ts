import { CfnOutput, Duration, Stack, Fn } from "aws-cdk-lib";
import { ILogGroup } from "aws-cdk-lib/aws-logs";
import { Construct } from "constructs";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { IRole } from "aws-cdk-lib/aws-iam";
import { PythonFunction } from "@aws-cdk/aws-lambda-python-alpha";
import { WES_BUCKET_NAME, WES_KEY_PARAMATER_NAME } from "../../constants";
import { Code, Function, Runtime } from "aws-cdk-lib/aws-lambda";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { getCommonParameter } from "../../util";

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

  public renderPythonLambda(scope: Construct, id: string, vpc: IVpc, role: IRole, environment: Record<string, string>): PythonFunction {
    const account: string = process.env.CDK_DEFAULT_ACCOUNT!;
    const region: string = process.env.CDK_DEFAULT_REGION!;
    return new Function(scope, id, {
      vpc,
      code: Code.fromBucket(Bucket.fromBucketName(scope, "WesAdapter", Fn.importValue(WES_BUCKET_NAME)), getCommonParameter(this, WES_KEY_PARAMATER_NAME)),
      handler: "index.handler",
      runtime: Runtime.PYTHON_3_9,
      environment,
      role,
      timeout: Duration.seconds(60),
    });
  }

  protected abstract getOutputs(): EngineOutputs;
}
