import { Construct } from "constructs";
import {
  AccessLogField,
  AccessLogFormat,
  ApiKeySourceType,
  AuthorizationType,
  ConnectionType,
  EndpointType,
  HttpIntegration,
  Integration,
  LambdaIntegration,
  LogGroupLogDestination,
  MethodLoggingLevel,
  RestApi,
  VpcLink,
} from "aws-cdk-lib/aws-apigateway";
import { INetworkLoadBalancer } from "aws-cdk-lib/aws-elasticloadbalancingv2";
import { AccountPrincipal, PolicyDocument, PolicyStatement } from "aws-cdk-lib/aws-iam";
import { ILogGroup, LogGroup } from "aws-cdk-lib/aws-logs";
import { IFunction } from "aws-cdk-lib/aws-lambda";

export interface ApiProxyProps {
  /**
   * An allowlist of AWS account IDs that can all this API.
   */
  allowedAccountIds: string[];
  /**
   * The load balancer to proxy.
   *
   * Required if lambda is not specified.
   */
  loadBalancer?: INetworkLoadBalancer;
  /**
   * The lambda to proxy.
   *
   * Required if loadBalancer is not specified.
   */
  lambda?: IFunction;
  /**
   * The name of the REST API.
   *
   * @default - ID of the RestApi construct.
   */
  apiName?: string;
}

export class ApiProxy extends Construct {
  public readonly accessLogGroup: ILogGroup;
  public readonly restApi: RestApi;

  constructor(scope: Construct, props: ApiProxyProps) {
    super(scope, "ApiProxy");

    if ((props.lambda && props.loadBalancer) || (!props.lambda && !props.loadBalancer)) {
      throw Error("Either lambda or loadBalancer must be specified, but not both");
    }

    this.accessLogGroup = new LogGroup(this, "AccessLogGroup");
    this.restApi = new RestApi(this, "Resource", {
      restApiName: props.apiName,
      endpointTypes: [EndpointType.PRIVATE],
      description: "API proxy endpoint for a service",
      apiKeySourceType: ApiKeySourceType.HEADER,
      deployOptions: {
        loggingLevel: MethodLoggingLevel.INFO,
        dataTraceEnabled: true,
        accessLogFormat: this.renderAccessLogFormat(),
        accessLogDestination: new LogGroupLogDestination(this.accessLogGroup),
      },
      policy: new PolicyDocument({
        statements: [
          new PolicyStatement({
            actions: ["execute-api:Invoke"],
            resources: ["execute-api:/*/*"],
            principals: props.allowedAccountIds.map((accountId) => new AccountPrincipal(accountId)),
          }),
        ],
      }),
    });

    const apiTarget = props.lambda ? new LambdaIntegration(props.lambda) : this.renderHttpTarget(props.loadBalancer!);
    this.restApi.root.addProxy({
      defaultIntegration: apiTarget,
      defaultMethodOptions: {
        authorizationType: AuthorizationType.IAM,
        requestParameters: { "method.request.path.proxy": true },
      },
    });
  }

  private renderAccessLogFormat(): AccessLogFormat {
    return AccessLogFormat.custom(
      JSON.stringify({
        requestId: AccessLogField.contextRequestId(),
        caller: AccessLogField.contextIdentityCaller(),
        callerAccountId: AccessLogField.contextAccountId(),
        user: AccessLogField.contextIdentityUser(),
        requestTime: AccessLogField.contextRequestTime(),
        httpMethod: AccessLogField.contextHttpMethod(),
        resourcePath: AccessLogField.contextResourcePath(),
        status: AccessLogField.contextStatus(),
        protocol: AccessLogField.contextProtocol(),
        responseLength: AccessLogField.contextResponseLength(),
        message: AccessLogField.contextErrorMessage(),
        validationError: AccessLogField.contextErrorValidationErrorString(),
      })
    );
  }

  private renderHttpTarget(loadBalancer: INetworkLoadBalancer): Integration {
    const vpcLink = new VpcLink(this, "VpcLink", { targets: [loadBalancer] });
    const apiUrl = `http://${loadBalancer.loadBalancerDnsName}/{proxy}`;
    return new HttpIntegration(apiUrl, {
      httpMethod: "ANY",
      options: {
        connectionType: ConnectionType.VPC_LINK,
        vpcLink,
        requestParameters: { "integration.request.path.proxy": "method.request.path.proxy" },
      },
    });
  }
}
