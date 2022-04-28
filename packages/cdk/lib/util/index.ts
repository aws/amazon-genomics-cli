import { Repository } from "aws-cdk-lib/aws-ecr";
import { ContainerImage, TaskDefinition } from "aws-cdk-lib/aws-ecs";
import { StringParameter } from "aws-cdk-lib/aws-ssm";
import { Maybe, ServiceContainer } from "../types";
import { Arn, CfnParameter, Fn, Stack } from "aws-cdk-lib";
import { Construct, Node } from "constructs";
import { APP_NAME } from "../constants";
import { SecureService } from "../constructs";
import { Protocol } from "aws-cdk-lib/aws-elasticloadbalancingv2";
import { IVpc, Subnet, SubnetSelection } from "aws-cdk-lib/aws-ec2";
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

/**
 * Obtains the content of a ParameterStore StringList parameter as a {@code string[]}. Because of an issue with the
 * way CloudFormation resolves these parameters to only a comma separated string at deploy time we need to use a
 * different approach which unfortunately requires us to know the length of the array a priori.
 * @param scope the {@code Construct} that will resolve and hold the resulting string[]
 * @param keySuffix the key of the StringList parameter
 * @param lengthSuffix the key for the parameter which will hold the length of the array
 */
export const getCommonParameterList = (scope: Construct, keySuffix: string, lengthSuffix: string): string[] => {
  const arrayLength = Number(StringParameter.valueFromLookup(scope, `/${APP_NAME}/_common/${lengthSuffix}`));
  const cfnParameter = new CfnParameter(scope, `ListParam${keySuffix}`, {
    type: "AWS::SSM::Parameter::Value<List<String>>",
    default: `/${APP_NAME}/_common/${keySuffix}`,
  });
  const list = cfnParameter.valueAsList;
  const subnetIds: string[] = [];
  for (let i = 0; i < arrayLength; i++) {
    subnetIds.push(Fn.select(i, list));
  }
  return subnetIds;
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

export const renderServiceWithTaskDefinition = (
  scope: Construct,
  id: string,
  serviceContainer: ServiceContainer,
  taskDefinition: TaskDefinition,
  vpc: IVpc,
  subnets: SubnetSelection
): SecureService => {
  return new SecureService(scope, id, {
    vpc,
    taskSubnets: subnets,
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

export function subnetSelectionFromIds(scope: Construct, subnetIds: string[]): SubnetSelection {
  const subnets = subnetIds.map((id, index) => {
    return Subnet.fromSubnetId(scope, `ContextSubnet${index}`, id);
  });

  return { subnets };
}
