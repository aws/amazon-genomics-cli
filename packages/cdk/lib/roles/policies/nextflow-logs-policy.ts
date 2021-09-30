import * as iam from "monocdk/aws-iam";
export class NextflowLogsPolicy extends iam.PolicyDocument {
  constructor() {
    super({
      assignSids: true,
      statements: [
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["logs:GetQueryResults", "logs:StopQuery"],
          resources: ["*"],
        }),
      ],
    });
  }
}
