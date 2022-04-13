import { CfnOutput, RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { AttributeType, BillingMode, ITable, ProjectionType, Table } from "aws-cdk-lib/aws-dynamodb";
import { IParameter, StringListParameter, StringParameter } from "aws-cdk-lib/aws-ssm";
import { GatewayVpcEndpointAwsService, InterfaceVpcEndpointService, IVpc, SubnetType, Vpc } from "aws-cdk-lib/aws-ec2";
import { Bucket, BucketEncryption, IBucket } from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";
import {
  PRODUCT_NAME,
  APP_NAME,
  VPC_PARAMETER_NAME,
  WES_KEY_PARAMETER_NAME,
  WES_BUCKET_NAME,
  VPC_PARAMETER_ID,
  VPC_SUBNETS_PARAMETER_NAME,
} from "../constants";
import { BucketDeployment, Source } from "aws-cdk-lib/aws-s3-deployment";
import * as path from "path";
import { homedir } from "os";
import { Asset } from "aws-cdk-lib/aws-s3-assets";

export interface ParameterProps {
  /**
   * The name of this parameter.
   *
   * All parameter names are prefixed with "/agc/_common/".
   */
  name: string;
  /**
   * The value stored in this parameter
   */
  value: string;
  /**
   * The description for this parameter
   *
   * @default none
   */
  description?: string;
}

export interface ListParameterProps {
  /**
   * The name of this parameter.
   *
   * All parameter names are prefixed with "/agc/_common/".
   */
  name: string;

  /**
   * The values stored in this parameter
   */
  value: string[];

  /**
   * The description for this parameter
   *
   * @default none
   */
  description?: string;
}

export interface CoreStackProps extends StackProps {
  /**
   * Name of S3 bucket to create or import
   */
  bucketName: string;
  /**
   * Key used to determine uniqueness of assets.
   */
  idempotencyKey: string;
  /**
   * Whether the bucket should be created or imported using bucketName
   *
   * @default true
   */
  createNewBucket?: boolean;
  /**
   * The name of the VPC the service should use
   *
   * @default - A new VPC is created
   */
  vpcId?: string;

  /**
   * A list of private subnet ids from within the VPC specified by vpcId to use for infrastructure deployment
   * @default - All private subnets are used
   */
  subnetIds?: string[];

  /**
   * A list of SSM parameters to create with the stack.
   *
   * @default none
   */
  parameters?: ParameterProps[];
  /**
   * If true, spawn a VPC with no NAT gateways or VPC endpoints (ie, no private subnets).
   * This **must** be used in conjunction with the uesPublicSubnets option for any context deployed within this account,
   * which can be specified in the `agc-project.yaml`.
   *
   * Note that this option risks security vulnerabilities if security groups are manually modified.
   *
   * @default false
   */
  usePublicSubnets?: boolean;
}

const parameterPrefix = `/${APP_NAME}/_common/`;

export class CoreStack extends Stack {
  public readonly vpc: IVpc;
  public readonly table: ITable;
  public readonly bucket: IBucket;

  constructor(scope: Construct, id: string, props: CoreStackProps) {
    super(scope, id, props);

    this.vpc = this.renderVpc(props.vpcId, props.usePublicSubnets);
    this.table = this.renderTable();
    this.bucket = this.renderBucket(props.bucketName, props.createNewBucket);

    new BucketDeployment(this, "BatchArtifacts", {
      sources: [Source.asset(path.join(__dirname, "../artifacts"))],
      destinationBucket: this.bucket,
      destinationKeyPrefix: "artifacts",
      prune: false,
      metadata: {
        "idempotency-key": props.idempotencyKey,
      },
    });
    new CfnOutput(this, VPC_PARAMETER_ID, { value: this.vpc.vpcId });

    const asset = new Asset(this, "WesAdapter", {
      path: path.join(homedir(), ".agc", "wes", "wes_adapter.zip"),
    });

    new CfnOutput(this, WES_BUCKET_NAME, { value: asset.s3BucketName, exportName: WES_BUCKET_NAME });
    // key cannot be a cfn output as it would block updates if some contexts are already deployed
    this.addParameter({ name: WES_KEY_PARAMETER_NAME, value: asset.s3ObjectKey, description: "The s3 key for the wes_adapter zip file" });

    this.addParameter({ name: VPC_PARAMETER_NAME, value: this.vpc.vpcId, description: `VPC ID for ${PRODUCT_NAME}` });
    props.parameters?.forEach((parameterProps) => this.addParameter(parameterProps));

    if (props.subnetIds != null && props.subnetIds.length > 0) {
      this.addParameterList({
        name: VPC_SUBNETS_PARAMETER_NAME,
        value: props.subnetIds,
        description: `Subnets of ${props.vpcId} to use for ${PRODUCT_NAME} infrastructure`,
      });
    }

    new CfnOutput(this, "TableName", { value: this.table.tableName });
  }

  private renderVpc(vpcId?: string, publicSubnets?: boolean): IVpc {
    if (vpcId) {
      return Vpc.fromLookup(this, "Vpc", { vpcId });
    } else if (publicSubnets) {
      return new Vpc(this, "Vpc", {
        subnetConfiguration: [
          {
            subnetType: SubnetType.PUBLIC,
            name: "Public",
          },
        ],
      });
    }

    const vpc = new Vpc(this, "Vpc", {
      gatewayEndpoints: {
        S3Endpoint: { service: GatewayVpcEndpointAwsService.S3 },
        DynamoDBEndpoint: { service: GatewayVpcEndpointAwsService.DYNAMODB },
      },
      subnetConfiguration: Vpc.DEFAULT_SUBNETS,
    });

    const subnetSelection = { subnets: vpc.privateSubnets, onePerAz: true };
    vpc.addInterfaceEndpoint(`${PRODUCT_NAME}LogsEndpoint`, {
      service: new InterfaceVpcEndpointService(`com.amazonaws.${this.region}.logs`),
      subnets: subnetSelection,
      open: true,
    });
    vpc.addInterfaceEndpoint(`${PRODUCT_NAME}EcrDkrEndpoint`, {
      service: new InterfaceVpcEndpointService(`com.amazonaws.${this.region}.ecr.dkr`),
      subnets: subnetSelection,
      open: true,
    });
    vpc.addInterfaceEndpoint(`${PRODUCT_NAME}EcrApiEndpoint`, {
      service: new InterfaceVpcEndpointService(`com.amazonaws.${this.region}.ecr.api`),
      subnets: subnetSelection,
      open: true,
    });
    vpc.addInterfaceEndpoint(`${PRODUCT_NAME}EcsAgentEndpoint`, {
      service: new InterfaceVpcEndpointService(`com.amazonaws.${this.region}.ecs-agent`),
      subnets: subnetSelection,
      open: true,
    });
    vpc.addInterfaceEndpoint(`${PRODUCT_NAME}EcsTelemEndpoint`, {
      service: new InterfaceVpcEndpointService(`com.amazonaws.${this.region}.ecs-telemetry`),
      subnets: subnetSelection,
      open: true,
    });
    vpc.addInterfaceEndpoint(`${PRODUCT_NAME}EcsEndpoint`, {
      service: new InterfaceVpcEndpointService(`com.amazonaws.${this.region}.ecs`),
      subnets: subnetSelection,
      open: true,
    });

    return vpc;
  }

  private renderTable(): ITable {
    const table = new Table(this, "Table", {
      tableName: PRODUCT_NAME,
      partitionKey: {
        name: "PK",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "SK",
        type: AttributeType.STRING,
      },
      timeToLiveAttribute: "expiry",
      billingMode: BillingMode.PAY_PER_REQUEST,
      removalPolicy: RemovalPolicy.DESTROY,
    });

    table.addGlobalSecondaryIndex({
      indexName: "gsi1",
      partitionKey: { name: "GSI1_PK", type: AttributeType.STRING },
      sortKey: { name: "GSI1_SK", type: AttributeType.STRING },
      projectionType: ProjectionType.ALL,
    });

    table.addLocalSecondaryIndex({
      indexName: "lsi1",
      sortKey: { name: "LSI1_SK", type: AttributeType.STRING },
      projectionType: ProjectionType.ALL,
    });

    table.addLocalSecondaryIndex({
      indexName: "lsi2",
      sortKey: { name: "LSI2_SK", type: AttributeType.STRING },
      projectionType: ProjectionType.ALL,
    });

    table.addLocalSecondaryIndex({
      indexName: "lsi3",
      sortKey: { name: "LSI3_SK", type: AttributeType.STRING },
      projectionType: ProjectionType.ALL,
    });

    return table;
  }

  private renderBucket(bucketName: string, createNew?: boolean): IBucket {
    if (createNew ?? true) {
      return new Bucket(this, "Bucket", {
        bucketName: bucketName,
        encryption: BucketEncryption.KMS_MANAGED,
        enforceSSL: true,
      });
    }
    return Bucket.fromBucketName(this, "Bucket", bucketName);
  }

  private addParameter(props: ParameterProps): IParameter {
    return new StringParameter(this, props.name, {
      parameterName: `${parameterPrefix}${props.name}`,
      stringValue: props.value,
      description: props.description,
    });
  }

  private addParameterList(props: ListParameterProps): IParameter {
    return new StringListParameter(this, props.name, {
      parameterName: `${parameterPrefix}${props.name}`,
      stringListValue: props.value,
      description: props.description,
    });
  }
}
