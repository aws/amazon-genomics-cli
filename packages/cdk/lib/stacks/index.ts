import { InterfaceVpcEndpoint, IVpcEndpoint } from "aws-cdk-lib/aws-ec2";
import { Construct } from "constructs";

export { CoreStack } from "./core-stack";
export { BatchConstruct } from "./engines/batch-construct";
export { CromwellEngineConstruct } from "./engines/cromwell-engine-construct";
export { NextflowEngineConstruct } from "./engines/nextflow-engine-construct";

const SSL_PORT = 443;

export function apiGatewayVpcEndpointFromId(app: Construct, id?: string): IVpcEndpoint[] {
  if (id) {
    return [
      InterfaceVpcEndpoint.fromInterfaceVpcEndpointAttributes(app, "ApiGatewayVpcEndpointLookup", {
        vpcEndpointId: id,
        port: SSL_PORT,
      }),
    ];
  } else {
    return [];
  }
}
