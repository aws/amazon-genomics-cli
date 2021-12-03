#!/usr/bin/env node
import { App } from "monocdk";
import "source-map-support/register";
import { getContext, getContextOrDefault } from "../../lib/util";
import { CoreStack } from "../../lib/stacks";
import { Maybe } from "../../lib/types";
import { APP_TAG_KEY, APP_NAME, PRODUCT_NAME, APP_ENV_NAME } from "../../lib/constants";
import { ParameterProps } from "../../lib/stacks/core-stack";
const app = new App();

const account: string = process.env.CDK_DEFAULT_ACCOUNT!;
const region: string = process.env.CDK_DEFAULT_REGION!;

const vpcId = getContextOrDefault<Maybe<string>>(app.node, "VPC_ID");
const bucketName = getContextOrDefault(app.node, `${APP_ENV_NAME}_BUCKET_NAME`, `${APP_NAME}-${account}-${region}`);
const createNewBucket = getContextOrDefault(app.node, `CREATE_${APP_ENV_NAME}_BUCKET`, "true").toLowerCase() == "true";
const ecrImageNames: { descriptionName: string; name: string }[] = [
  { name: "wes", descriptionName: "WES" },
  { name: "cromwell", descriptionName: "Cromwell" },
  { name: "nextflow", descriptionName: "Nextflow" },
  { name: "miniwdl", descriptionName: "MiniWDL" },
];

const generateEcrImages = (): ParameterProps[] => {
  const engineParameters: ParameterProps[] = [];

  for (const ecrImage of ecrImageNames) {
    engineParameters.push(
      {
        name: `${ecrImage.name}/ecr-repo/account`,
        value: getContext(app.node, `ECR_${ecrImage.name.toUpperCase()}_ACCOUNT_ID`),
        description: `Account ID of ECR that contains the ${ecrImage.descriptionName} docker image`,
      },
      {
        name: `${ecrImage.name}/ecr-repo/region`,
        value: getContext(app.node, `ECR_${ecrImage.name.toUpperCase()}_REGION`),
        description: `Region of ECR that contains the ${ecrImage.descriptionName} docker image`,
      },
      {
        name: `${ecrImage.name}/ecr-repo/tag`,
        value: getContext(app.node, `ECR_${ecrImage.name.toUpperCase()}_TAG`),
        description: `Docker tag for the ${ecrImage.descriptionName} image`,
      },
      {
        name: `${ecrImage.name}/ecr-repo/repository`,
        value: getContext(app.node, `ECR_${ecrImage.name.toUpperCase()}_REPOSITORY`),
        description: `ECR repository for the ${ecrImage.descriptionName} image`,
      }
    );
  }

  return engineParameters;
};

new CoreStack(app, `${PRODUCT_NAME}-Core`, {
  vpcId,
  bucketName,
  createNewBucket,
  env: {
    account,
    region,
  },
  tags: {
    [APP_TAG_KEY]: APP_NAME,
  },
  parameters: [
    {
      name: "bucket",
      value: bucketName,
      description: "S3 bucket which contains outputs, intermediate results, and other project-specific data",
    },
    {
      name: "installed-artifacts/s3-root-url",
      value: "s3://healthai-public-assets-us-east-1/batch/1.0.2/artifacts",
      description: "S3 root url for assets",
    },
    ...generateEcrImages(),
  ],
});
