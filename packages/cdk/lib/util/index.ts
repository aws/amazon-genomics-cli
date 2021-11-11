import { Repository } from "monocdk/aws-ecr";
import { CloudMapOptions, ContainerImage, LogDriver, TaskDefinition } from "monocdk/aws-ecs";
import { StringParameter } from "monocdk/aws-ssm";
import { Maybe, ServiceContainer } from "../types";
import { Arn, Construct, ConstructNode, Stack } from "monocdk";
import { APP_NAME } from "../constants";
import { SecureService } from "../constructs";
import { Protocol } from "monocdk/aws-elasticloadbalancingv2";
import { IVpc } from "monocdk/aws-ec2";
import { IRole } from "monocdk/aws-iam";
import { LogConfiguration, LogDriver as BatchLogDriver } from "monocdk/aws-batch";
import { ILogGroup } from "monocdk/lib/aws-logs/lib/log-group";
import { PythonFunction } from "monocdk/aws-lambda-python";
import { Runtime } from "monocdk/aws-lambda";
import { Duration } from "monocdk";

export const getContext = (node: ConstructNode, key: string): string => {
  const context = getContextOrDefault(node, key, undefined);
  if (!context) {
    throw Error(`Context cannot be null for key '${key}'`);
  }
  return context;
};

export const getContextOrDefault = <T extends Maybe<string>>(node: ConstructNode, key: string, defaultValue?: T): T => {
  const value = node.tryGetContext(key);
  return !value || value == "" ? defaultValue : value;
};

export const getCommonParameter = (scope: Construct, keySuffix: string): string => {
  return StringParameter.valueFromLookup(scope, `/${APP_NAME}/_common/${keySuffix}`);
};

export const getProjectParameter = (scope: Construct, project: string, keySuffix: string): string => {
  return StringParameter.valueFromLookup(scope, `/${APP_NAME}/${project}/${keySuffix}`);
};

export const createEcrImage = (scope: Construct, designation: string): ContainerImage => {
  const propertyPrefix = `${designation}/ecr-repo`;
  const accountId = getCommonParameter(scope, `${propertyPrefix}/account`);
  const region = getCommonParameter(scope, `${propertyPrefix}/region`);
  const tag = getCommonParameter(scope, `${propertyPrefix}/tag`);
  const repositoryName = getCommonParameter(scope, `${propertyPrefix}/repository`);
  const ecrArn = `arn:aws:ecr:${region}:${accountId}:repository/${repositoryName}`;
  const repository = Repository.fromRepositoryAttributes(scope, repositoryName, {
    repositoryName,
    repositoryArn: ecrArn,
  });
  return ContainerImage.fromEcrRepository(repository, tag);
};

const defaultHealthCheckPath = "/ga4gh/wes/v1/service-info";

export const renderServiceWithContainer = (
  scope: Construct,
  id: string,
  serviceContainer: ServiceContainer,
  vpc: IVpc,
  taskRole: IRole,
  logGroup: ILogGroup,
  cloudMapOptions?: CloudMapOptions
): SecureService => {
  return new SecureService(scope, id, {
    vpc,
    serviceName: serviceContainer.serviceName,
    cpu: serviceContainer.cpu,
    memoryLimitMiB: serviceContainer.memoryLimitMiB,
    cloudMapOptions,
    healthCheck: {
      path: serviceContainer.healthCheckPath ?? defaultHealthCheckPath,
      protocol: Protocol.HTTP,
    },
    taskImageOptions: {
      taskRole,
      image: createEcrImage(scope, serviceContainer.imageConfig.designation),
      environment: serviceContainer.environment,
      containerPort: serviceContainer.containerPort,
      logDriver: LogDriver.awsLogs({ logGroup, streamPrefix: id }),
    },
  });
};

export const renderServiceWithTaskDefinition = (
  scope: Construct,
  id: string,
  serviceContainer: ServiceContainer,
  taskDefinition: TaskDefinition,
  vpc: IVpc
): SecureService => {
  return new SecureService(scope, id, {
    vpc,
    serviceName: serviceContainer.serviceName,
    taskDefinition: taskDefinition,
    healthCheck: {
      path: serviceContainer.healthCheckPath ?? defaultHealthCheckPath,
      protocol: Protocol.HTTP,
    },
  });
};

export function renderBatchLogConfiguration(scope: Construct, logGroup: ILogGroup): LogConfiguration {
  return {
    logDriver: BatchLogDriver.AWSLOGS,
    options: {
      "awslogs-group": logGroup.logGroupName,
    },
  };
}

export function batchArn(scope: Construct, resource: string, resourcePrefix = "*"): string {
  return Arn.format({ resource: `${resource}/${resourcePrefix}`, service: "batch" }, Stack.of(scope));
}
export const renderPythonLambda = (
  scope: Construct,
  id: string,
  vpc: IVpc,
  role: IRole,
  codePath: string,
  environment: Record<string, string>
): PythonFunction => {
  return new PythonFunction(scope, id, {
    vpc,
    entry: codePath,
    runtime: Runtime.PYTHON_3_9,
    environment,
    role,
    timeout: Duration.seconds(30),
  });
};
