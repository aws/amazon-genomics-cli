import { Repository } from "aws-cdk-lib/aws-ecr";
import { CloudMapOptions, ContainerImage, LogDriver, TaskDefinition } from "aws-cdk-lib/aws-ecs";
import { StringParameter } from "aws-cdk-lib/aws-ssm";
import { Maybe, ServiceContainer } from "../types";
import { Arn, Stack } from "aws-cdk-lib";
import { Construct, Node } from "constructs";
import { APP_NAME } from "../constants";
import { SecureService } from "../constructs";
import { Protocol } from "aws-cdk-lib/aws-elasticloadbalancingv2";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { IRole } from "aws-cdk-lib/aws-iam";
import { LogConfiguration, LogDriver as BatchLogDriver } from "@aws-cdk/aws-batch-alpha";
import { ILogGroup } from "aws-cdk-lib/aws-logs";

export const getContext = (node: Node, key: string): string => {
  const context = getContextOrDefault(node, key, undefined);
  if (!context) {
    throw Error(`Context cannot be null for key '${key}'`);
  }
  return context;
};

export const getContextOrDefault = <T extends Maybe<string>>(node: Node, key: string, defaultValue?: T): T => {
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
  const engineName = designation.toUpperCase();
  const accountId = getContext(scope.node, `ECR_${engineName}_ACCOUNT_ID`);
  const region = getContext(scope.node, `ECR_${engineName}_REGION`);
  const tag = getContext(scope.node, `ECR_${engineName}_TAG`);
  const repositoryName = getContext(scope.node, `ECR_${engineName}_REPOSITORY`);
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

export function ec2Arn(scope: Construct, resource: string, resourcePrefix = "*"): string {
  return Arn.format({ resource: `${resource}/${resourcePrefix}`, service: "ec2" }, Stack.of(scope));
}
