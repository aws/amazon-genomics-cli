#!/usr/bin/env node
import { App } from "monocdk";
import "source-map-support/register";
import { getContext, getContextOrDefault } from "../../lib/util";
import { CoreStack } from "../../lib/stacks";
import { Maybe } from "../../lib/types";
import { APP_TAG_KEY, APP_NAME, PRODUCT_NAME, APP_ENV_NAME } from "../../lib/constants";
const app = new App();

const account: string = process.env.CDK_DEFAULT_ACCOUNT!;
const region: string = process.env.CDK_DEFAULT_REGION!;

const vpcId = getContextOrDefault<Maybe<string>>(app.node, "VPC_ID");
const bucketName = getContextOrDefault(app.node, `${APP_ENV_NAME}_BUCKET_NAME`, `${APP_NAME}-${account}-${region}`);
const createNewBucket = getContextOrDefault(app.node, `CREATE_${APP_ENV_NAME}_BUCKET`, "true").toLowerCase() == "true";

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
    {
      name: "wes/ecr-repo/account",
      value: getContext(app.node, "ECR_WES_ACCOUNT_ID"),
      description: "Account ID of ECR that contains the WES docker image",
    },
    {
      name: "wes/ecr-repo/region",
      value: getContext(app.node, "ECR_WES_REGION"),
      description: "Region of ECR that contains the WES docker image",
    },
    {
      name: "wes/ecr-repo/tag",
      value: getContext(app.node, "ECR_WES_TAG"),
      description: "Docker tag for the WES image",
    },
    {
      name: "wes/ecr-repo/repository",
      value: getContext(app.node, "ECR_WES_REPOSITORY"),
      description: "ECR repository for the WES image",
    },
    {
      name: "cromwell/ecr-repo/account",
      value: getContext(app.node, "ECR_CROMWELL_ACCOUNT_ID"),
      description: "Account ID of ECR that contains the Cromwell docker image",
    },
    {
      name: "cromwell/ecr-repo/region",
      value: getContext(app.node, "ECR_CROMWELL_REGION"),
      description: "Region of ECR that contains the Cromwell docker image",
    },
    {
      name: "cromwell/ecr-repo/tag",
      value: getContext(app.node, "ECR_CROMWELL_TAG"),
      description: "Docker tag for the Cromwell image",
    },
    {
      name: "cromwell/ecr-repo/repository",
      value: getContext(app.node, "ECR_CROMWELL_REPOSITORY"),
      description: "ECR repository for the Cromwell image",
    },
    {
      name: "nextflow/ecr-repo/account",
      value: getContext(app.node, "ECR_NEXTFLOW_ACCOUNT_ID"),
      description: "Account ID of ECR that contains the Nextflow docker image",
    },
    {
      name: "nextflow/ecr-repo/region",
      value: getContext(app.node, "ECR_NEXTFLOW_REGION"),
      description: "Region of ECR that contains the Nextflow docker image",
    },
    {
      name: "nextflow/ecr-repo/tag",
      value: getContext(app.node, "ECR_NEXTFLOW_TAG"),
      description: "Docker tag for the Nextflow image",
    },
    {
      name: "nextflow/ecr-repo/repository",
      value: getContext(app.node, "ECR_NEXTFLOW_REPOSITORY"),
      description: "ECR repository for the Nextflow image",
    },
    {
      name: "miniwdl/ecr-repo/account",
      value: getContext(app.node, "ECR_MINIWDL_ACCOUNT_ID"),
      description: "Account ID of ECR that contains the MiniWDL docker image",
    },
    {
      name: "miniwdl/ecr-repo/region",
      value: getContext(app.node, "ECR_MINIWDL_REGION"),
      description: "Region of ECR that contains the MiniWDL docker image",
    },
    {
      name: "miniwdl/ecr-repo/tag",
      value: getContext(app.node, "ECR_MINIWDL_TAG"),
      description: "Docker tag for the MiniWDL image",
    },
    {
      name: "miniwdl/ecr-repo/repository",
      value: getContext(app.node, "ECR_MINIWDL_REPOSITORY"),
      description: "ECR repository for the MiniWDL image",
    },
  ],
});
