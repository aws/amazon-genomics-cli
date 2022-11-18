import { TaggedResourceTypes } from "./types/tagged-resource-types";

export const PRODUCT_NAME = "Agc";
export const APP_NAME = "agc";
export const APP_ENV_NAME = "AGC";
export const APP_TAG_KEY = "application-name";
export const PROJECT_TAG_KEY = `${APP_NAME}-project`;
export const CONTEXT_TAG_KEY = `${APP_NAME}-context`;
export const USER_ID_TAG_KEY = `${APP_NAME}-user-id`;
export const USER_EMAIL_TAG_KEY = `${APP_NAME}-user-email`;
export const ENGINE_TAG_KEY = `${APP_NAME}-engine`;
export const ENGINE_TYPE_TAG_KEY = `${APP_NAME}-engine-type`;
export const AGC_VERSION_KEY = `${APP_NAME}-version`;
export const VPC_PARAMETER_NAME = "vpc";
export const VPC_PARAMETER_ID = "VpcId";
export const VPC_SUBNETS_PARAMETER_NAME = "InfraSubnets";
export const VPC_NUMBER_SUBNETS_PARAMETER_NAME = "NumInfraSubnets";
export const WES_KEY_PARAMETER_NAME = "WesAdapterZipKeyParameter";
export const WES_BUCKET_NAME = "WesAdapterZipBucket";
export const COMPUTE_IMAGE_PARAMETER_NAME = "ComputeEnvImage";

export const ENGINE_CROMWELL = "cromwell";
export const ENGINE_MINIWDL = "miniwdl";
export const ENGINE_NEXTFLOW = "nextflow";
export const ENGINE_SNAKEMAKE = "snakemake";
export const ENGINE_TOIL = "toil";

export const TAGGED_RESOURCE_TYPES: TaggedResourceTypes = ["volume"];
