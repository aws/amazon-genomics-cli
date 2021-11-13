import { ILogGroup, LogGroup } from "monocdk/aws-logs";
import { Construct } from "monocdk";
import { IVpc } from "monocdk/aws-ec2";

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
