import { ILogGroup, LogGroup } from "aws-cdk-lib/aws-logs";
import { Construct } from "constructs";
import { IVpc } from "aws-cdk-lib/aws-ec2";

export interface EngineProps {
  readonly vpc: IVpc;
  readonly rootDirS3Uri: string;
}

export class Engine extends Construct {
  readonly logGroup: ILogGroup;

  constructor(scope: Construct, id: string) {
    super(scope, id);
    this.logGroup = new LogGroup(this, "EngineLogGroup");
  }
}
