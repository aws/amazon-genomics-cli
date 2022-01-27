import { CfnOutput, RemovalPolicy, Stack, StackProps } from "monocdk";
import { AttributeType, BillingMode, ITable, ProjectionType, Table } from "monocdk/aws-dynamodb";
import { StringParameter, IParameter } from "monocdk/aws-ssm";
import {GatewayVpcEndpointAwsService, InterfaceVpcEndpointService, IVpc, Vpc, VpcProps} from "monocdk/aws-ec2";
import { Bucket, BucketEncryption, IBucket } from "monocdk/aws-s3";
import { Construct } from "constructs";
import { PRODUCT_NAME, APP_NAME, VPC_PARAMETER_NAME } from "../constants";

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

export interface CoreStackProps extends StackProps {
  /**
   * Name of S3 bucket to create or import
   */
  bucketName: string;
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
   * A list of SSM parameters to create with the stack.
   *
   * @default none
   */
  parameters?: ParameterProps[];
  /**
   * If true, use public subnets instead of private subnets, and do not create the VPC endpoints and the NAT gateway.
   * This will result in significantly lower ongoing costs when no job is running.
   *
   * @default false
   */
  publicSubnets?: boolean;
}

const parameterPrefix = `/${APP_NAME}/_common/`;

export class CoreStack extends Stack {
  public readonly vpc: IVpc;
  public readonly table: ITable;
  public readonly bucket: IBucket;

  constructor(scope: Construct, id: string, props: CoreStackProps) {
    super(scope, id, props);

    this.vpc = this.renderVpc(props.vpcId, props.publicSubnets);
    this.table = this.renderTable();
    this.bucket = this.renderBucket(props.bucketName, props.createNewBucket);

    this.addParameter({ name: VPC_PARAMETER_NAME, value: this.vpc.vpcId, description: `VPC ID for ${PRODUCT_NAME}` });
    props.parameters?.forEach((parameterProps) => this.addParameter(parameterProps));

    new CfnOutput(this, "TableName", { value: this.table.tableName });
  }

  private renderVpc(vpcId?: string, publicSubnets?: boolean): IVpc {
    if (vpcId) {
      return Vpc.fromLookup(this, "Vpc", { vpcId });
    }

    const vpc = new Vpc(this, "Vpc", {
      gatewayEndpoints: publicSubnets ? {} : {
          S3Endpoint: { service: GatewayVpcEndpointAwsService.S3 },
        },
      subnetConfiguration: publicSubnets ? Vpc.DEFAULT_SUBNETS_NO_NAT : Vpc.DEFAULT_SUBNETS
    });

    if (!publicSubnets) {
      const subnetSelection = {subnets: vpc.privateSubnets, onePerAz: true};
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
    }

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
}
