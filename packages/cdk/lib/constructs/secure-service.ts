import { Construct } from "constructs";
import { Port } from "aws-cdk-lib/aws-ec2";
import { BaseService, FargatePlatformVersion, ICluster } from "aws-cdk-lib/aws-ecs";
import { NetworkLoadBalancedFargateService, NetworkLoadBalancedFargateServiceProps } from "aws-cdk-lib/aws-ecs-patterns";
import { HealthCheck, NetworkLoadBalancer } from "aws-cdk-lib/aws-elasticloadbalancingv2";

export type ServiceOptions = Omit<NetworkLoadBalancedFargateServiceProps, "publicLoadBalancer" | "platformVersion">;

export interface SecureServiceProps extends ServiceOptions {
  /**
   * Configuration for the service's health check.
   *
   * @default - No health check is done.
   */
  healthCheck?: HealthCheck;
}

export class SecureService extends Construct {
  public readonly service: BaseService;
  public readonly cluster: ICluster;
  public readonly loadBalancer: NetworkLoadBalancer;

  private readonly resource: NetworkLoadBalancedFargateService;

  constructor(scope: Construct, id: string, props: SecureServiceProps) {
    super(scope, id);

    this.resource = new NetworkLoadBalancedFargateService(this, "Resource", {
      ...props,
      publicLoadBalancer: false,
      platformVersion: FargatePlatformVersion.VERSION1_4,
    });
    if (props.healthCheck) {
      this.resource.targetGroup.configureHealthCheck(props.healthCheck);
    }

    this.service = this.resource.service;
    this.cluster = this.resource.cluster;
    this.loadBalancer = this.resource.loadBalancer;

    this.loadBalancer.setAttribute("load_balancing.cross_zone.enabled", "true");

    if (props.taskImageOptions) {
      const containerPort = props.taskImageOptions.containerPort;
      this.service.connections.allowFromAnyIpv4(Port.tcp(containerPort ?? 80));
    }

    if (props.taskDefinition) {
      const containerPort = props.taskDefinition?.defaultContainer?.containerPort;
      this.service.connections.allowFromAnyIpv4(Port.tcp(containerPort ?? 80));
    }
  }
}
