export interface ImageConfiguration {
  /**
   * Parameter designation to look up additional repository information.
   */
  designation: string;
}

export interface ServiceContainer {
  /**
   * Configuration needed to retrieve the image used to deploy the service.
   */
  imageConfig: ImageConfiguration;
  /**
   * The name of the service.
   */
  serviceName: string;
  /**
   * The port number on the container that is bound to the service's host port.
   *
   * This default is set in the underlying NetworkLoadBalancedFargateService construct.
   * @default 80
   */
  containerPort?: number;
  /**
   * The number of cpu units used by the task.
   *
   * This default is set in the underlying FargateTaskDefinition construct.
   * @default 256
   */
  cpu?: number;
  /**
   * The amount (in MiB) of memory used by the task.
   *
   * This default is set in the underlying FargateTaskDefinition construct.
   * @default 512
   */
  memoryLimitMiB?: number;
  /**
   * Path to hit for the service's health check.
   *
   * @default "/ga4gh/wes/v1/service-info"
   */
  healthCheckPath?: string;
  /**
   * The environment variables to pass to the container.
   *
   * @default - No environment variables.
   */
  environment?: {
    [key: string]: string;
  };
}
