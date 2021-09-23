// TODO: Use official fargate support once https://github.com/aws/aws-cdk/pull/13591 is merged
export enum ComputeType {
  /**
   * Resources will be EC2 On-Demand resources.
   */
  ON_DEMAND = "EC2",
  /**
   * Resources will be EC2 SpotFleet resources.
   *
   */
  SPOT = "SPOT",
  /**
   * Resources will be Fargate resources.
   */
  FARGATE = "FARGATE",

  /**
   * Resources will be Fargate spot resources.
   */
  FARGATE_SPOT = "FARGATE_SPOT",
}
