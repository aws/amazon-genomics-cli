import { PolicyStatement, Effect } from '@aws-cdk/aws-iam';
import { Arn, ArnComponents, Stack } from '@aws-cdk/core';

function action(svc: string, action: string) : string {
    return [svc, action].join(":")
}


function actions(svc: string, ...actions: string[]) : string[] {
    return actions.map((x) => action(svc, x))
}


export class AgcPermissions {
    stack: Stack;
    constructor(stack: Stack) {
        this.stack = stack;
    }

    private arn(components: ArnComponents) : string {
        return Arn.format(components, this.stack);
    }

    vpc() : PolicyStatement[] {
        const svc = "ec2";
        return [
            new PolicyStatement({
                
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "CreateVpc",
                    "DeleteVpc",
                    "CreateNatGateway",
                    "DeleteNatGateway",
                    "CreateSecurityGroup",
                    "DeleteSecurityGroup",
                    "CreateInternetGateway",
                    "DeleteInternetGateway",
                    "CreateRoute",
                    "DeleteRoute",
                    "CreateSubnet",
                    "DeleteRouteTable",
                    "AuthorizeSecurityGroupIngress",
                    "AuthorizeSecurityGroupEgress",
                    "RevokeSecurityGroupIngress",
                    "RevokeSecurityGroupEgress",
                    "CreateVpcEndpoint",
                    "DeleteVpcEndpoints",
                    "AllocateAddress",
                    "AssociateAddress",
                    "ReleaseAddress",
                ),
                resources: [
                  this.arn({service: svc, region: "*", resource: "vpc", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "natgateway", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "security-group", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "internet-gateway", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "subnet", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "route-table", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "vpc-endpoint", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "ipv4pool-ec2", resourceName: "*"}),
                  this.arn({service: svc, region: "*", resource: "elastic-ip", resourceName: "*"}),
                ],
              }),
            new PolicyStatement({
                
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "DeleteSubnet",
                    "CreateRouteTable",
                    "DeleteTags",
                    "CreateTags",
                    "ModifyVpcAttribute",
                    "AttachInternetGateway",

                    "DescribeInternetGateways",
                    "DescribeVpcs",
                    "DescribeAvailabilityZones",
                    "DescribeSecurityGroups",
                    "DescribeAccountAttributes",
                    "DescribeSubnets",
                    "DescribeRouteTables",
                    "DescribeVpcEndpointServices",
                    "DetachInternetGateway",
                    "DescribeVpcendpoints",
                    "DescribeAddresses",
                    "DescribeNatGateways",

                    "ModifySubnetAttribute",
                    "DisassociateRouteTable",
                    "AssociateRouteTable",
                ),
                resources: [
                  "*"
                ],
              }),
        ]
    }

    ec2() : PolicyStatement[] {
        const svc = "ec2";
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "CreateSecurityGroup",
                    "DeleteSecurityGroup",
                    "AuthorizeSecurityGroupIngress",
                    "AuthorizeSecurityGroupEgress",
                    "RevokeSecurityGroupIngress",
                    "RevokeSecurityGroupEgress",

                    "*LaunchTemplate*",
                    
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "vpc", resourceName: "*"}),
                    this.arn({service: svc, region: "*", resource: "security-group", resourceName: "*"}),
                    this.arn({service: svc, region: "*", resource: "launch-template", resourceName: "*"}),
                ]
            }),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "DescribeVpcs",  // required for private hosted zones
                    "DescribeRegions",  // required for private hosted zones
                    "DescribeSubnets",
                    "DescribeRouteTables",
                    "DescribeVpnGateways",
                    "DescribeSecurityGroups",
                    "CreateTags",
                    "CreateVpcEndpointServiceConfiguration",
                    "DeleteVpcEndpointServiceConfigurations",
                    "DescribeVpcEndpointServiceConfigurations",
                    "ModifyVpcEndpointServicePermissions"
                ),
                resources: [
                    "*",
                ]
            })
        ]
    }

    s3Create() : PolicyStatement[] {
        const svc = "s3";
        return [
            new PolicyStatement({
                
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "CreateBucket",
                    "PutBucketPolicy",
                    "PutBucketTagging",
                    "PutEncryptionConfiguration",
                ),
                resources: [
                    this.arn({service: svc, region: "", account: "", resource: "agc-*"})
                ],
            })
        ]
    }
    
    s3Destroy() : PolicyStatement[] {
        const svc = "s3";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "DeleteBucket",
                "DeleteBucketPolicy",
            ),
            resources: [
              this.arn({service: svc, region: "", account: "", resource: "agc-*"})
            ],
          })]
    }
    
    s3Read() : PolicyStatement[] {
        const svc = "s3";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "GetBucketPolicy",
                "GetBucketTagging",
                "GetEncryptionConfiguration",
                "ListBucket",
                "GetObject",
                "GetObjectTagging"
            ),
            resources: [
              this.arn({service: svc, region: "", account: "", resource: "agc-*"})
            ],
          })]
    }
    
    s3Write() : PolicyStatement[] {
        const svc = "s3";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "PutObject",
                "PutObjectTagging",
                "DeleteObjectTagging",
                "DeleteObject",
            ),
            resources: [
              this.arn({service: svc, region: "", account: "", resource: "agc-*"})
            ],
          })]
    }

    s3CDK() : PolicyStatement[] {
        const svc = "s3";
        return [new PolicyStatement({
            effect: Effect.ALLOW,
            actions: actions(svc,
                "ListBucket",
                "GetObject",
                "GetBucketLocation",
                "PutObject",
                "DeleteObject",
            ),
            resources: [
              this.arn({service: svc, region: "", account: "", resource: "cdktoolkit-*"}),
            ],
        })]
    }

    dynamodbCreate() : PolicyStatement[] {
        const svc = "dynamodb";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "CreateTable",
                "UpdateTable",
                "BatchWriteItem",
                "TagResource",
            ),
            resources: [
              this.arn({service: svc, region: "*", resource: "table", resourceName: "Agc"})
            ],
          })]
          
    }
    
    dynamodbDestroy() : PolicyStatement[] {
        const svc = "dynamodb";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "DeleteTable",
                "UntagResource",
            ),
            resources: [
              this.arn({service: svc, region: "*", resource: "table", resourceName: "Agc"})
            ],
          })]
    }
    
    dynamodbRead() : PolicyStatement[] {
        const svc = "dynamodb";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "DescribeTable",
                "DescribeTimeToLive",
                "ListTables", // *
                "ListTagsOfResource",
                "BatchGetItem",
                "GetItem",
                "Scan",
                "Query",
            ),
            resources: [
              this.arn({service: svc, region: "*", resource: "table", resourceName: "Agc*"})
            ],
          })]
    }
    
    dynamodbWrite() : PolicyStatement[] {
        const svc = "dynamodb";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "UpdateTimeToLive",
                "BatchWriteItem",
                "DeleteItem",
                "PutItem",
                "UpdateItem",
            ),
            resources: [
              this.arn({service: svc, region: "*", resource: "table", resourceName: "Agc"})
            ],
          })]
    }

    ssmCreate() : PolicyStatement[] {
        const svc = "ssm";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "LabelParameterVersion",
                "PutParameter",
                "AddTagsToResource",
            ),
            resources: [
              this.arn({service: svc, region: "*", resource: "parameter", resourceName: "agc/*"})
            ]
          })]
    }
    
    ssmDestroy() : PolicyStatement[] {
        const svc = "ssm";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
              "DeleteParameter",
              "DeleteParameters",
              "RemoveTagsFromResource",
            ),
            resources: [
              this.arn({service: svc, region: "*", resource: "parameter", resourceName: "agc/*"})
            ]
          })]
    }
    
    ssmRead() : PolicyStatement[] {
        const svc = "ssm";
        return [new PolicyStatement({
            
            effect: Effect.ALLOW,
            actions: actions(svc,
                "Describe*",
                //"DescribeParameters",
                "Get*",
                //"GetParameter",
                //"GetParameters",
                //"GetParametersByPath",
                "ListTagsForResource",
            ),
            resources: [
              this.arn({service: svc, region: "*", resource: "parameter", resourceName: "agc/*"})
            ]
          })]
    }

    iam() : PolicyStatement[] {
        const svc = "iam";
        return [
            new PolicyStatement({
                
                effect: Effect.ALLOW, 
                actions: actions(svc,
                    "*Role",
                    // "GetRole",
                    // "CreateRole",
                    // "DeleteRole",
                    // "PassRole",
                    "ListRoleTags",
                    // "TagRole",
                    // "UntagRole",
                    "*RolePolicy",
                    // "GetRolePolicy",
                    // "PutRolePolicy",
                    // "DeleteRolePolicy",
                    // "AttachRolePolicy",
                    // "DetachRolePolicy",
                    "ListRolePolicies",
                    "ListAttachedRolePolicies",
                    "ListInstanceProfilesForRole",

                    "*InstanceProfile",
                    // "GetInstanceProfile",
                    // "CreateInstanceProfile",
                    // "DeleteInstanceProfile",
                    // "ListInstanceProfileTags",
                    // "TagInstanceProfile",
                    // "UntagInstanceProfile",
                    // "AddRoleToInstanceProfile",
                    // "RemoveRoleFromInstanceProfile",
                ),
                resources: [
                    this.arn({service: svc, region: "", resource: "role", resourceName: "Agc-*"}),
                    this.arn({service: svc, region: "", resource: "instance-profile", resourceName: "Agc-*"})
                ]
            }),
            new PolicyStatement({
                
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "PassRole",
                ),
                resources: [
                    this.arn({service: svc, region: "", resource: "role", resourceName: "aws-service-role/*"})
                ]
            }),
        ]
    }

    private cloudformationCommon() {
        const svc = "cloudformation";
    
        return new PolicyStatement({
            effect: Effect.ALLOW,
            actions: actions(svc, 
                "Describe*",
                // "DescribeStacks",
                // "DescribeStackEvents",
    
                // "DescribeStackResource",
                // "DescribeStackResourceDrifts",
                // "DescribeStackResources",
    
                "List*",
                // "ListStacks",
                // "ListChangeSets",
                // "ListExports",
                // "ListStackResources",
                
                "Detect*",
                // "DetectStackDrift",
                // "DetectStackResourceDrift",
                
                "*Stack",
                //"CreateStack",
                //"DeleteStack",
                //"UpdateStack",
                
                "GetTemplate",
    
                "*ChangeSet",
                // "CreateChangeSet",
                // "DescribeChangeSet",
                // "ExecuteChangeSet",
                // "DeleteChangeSet",
                
                "TagResource",
                "UntagResource",
            ),
            resources: [
                this.arn({service: "cloudformation", region: "*", resource: "stack", resourceName: "CDKToolkit*"}),
            ]
        })
    }
    
    cloudformationAdmin() : PolicyStatement[] {
        let stmt = this.cloudformationCommon();
    
        stmt.addResources(
            this.arn({service: "cloudformation", region: "*", resource: "stack", resourceName: "Agc-*"}) // allow on all AGC related stacks
        );
    
        return [
            stmt,
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions("cloudformation",
                    "ListStacks"
                ),
                resources: [
                    this.arn({service: "cloudformation", region: "*", resource: "stack", resourceName: "*"})
                ]
            })
        ]
    }
    
    cloudformationUser() : PolicyStatement[] {
        let stmt = this.cloudformationCommon();
    
        stmt.addResources(
            this.arn({service: "cloudformation", region: "*", resource: "stack", resourceName: "Agc-*-*"}) // allow only on non Agc-Core stacks
        );
        
        return [
            stmt,
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions("cloudformation",
                    "ListStacks"
                ),
                resources: [
                    this.arn({service: "cloudformation", region: "*", resource: "stack", resourceName: "*"})
                ]
            })
        ]
    }

    batch() : PolicyStatement[] {
        const svc = "batch";
        return [
            new PolicyStatement({
                
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "*ComputeEnvironment",
                    // "CreateComputeEnvironment",
                    // "DeleteComputeEnvironment",
                    "*JobQueue",
                    // "CreateJobQueue",
                    // "DeleteJobQueue",
                    "*Tag*",
                    //"ListTagsForResource",
                    //"TagResource",
                    //"UntagResource",
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "compute-environment", resourceName: "TaskBatch*"}),
                    this.arn({service: svc, region: "*", resource: "job-queue", resourceName: "TaskBatch*"}),
                ]
            }),
            new PolicyStatement({
                
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "*Job",
                    // "CancelJob",
                    // "SubmitJob",
                    // "TerminateJob",
                    "*JobDefinition",
                    // "DeregisterJobDefinition",
                    // "RegisterJobDefinition",
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "*"}),
                ]
            }),
            new PolicyStatement({
                
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "Describe*",
                    // "DescribeComputeEnvironments", //* 
                    // "DescribeJobDefinitions", //*
                    // "DescribeJobQueues", //*
                    // "DescribeJobs", //*
                    "List*",
                    // "ListJobs", //*
                    "*Tag*",
                    //"ListTagsForResource",
                    //"TagResource",
                    //"UntagResource",
                ),
                resources: [
                    "*",
                ]
            }),
        ]
    }
    
    ecs() : PolicyStatement[] {
        const svc = "ecs"
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "*Cluster",
                    //"CreateCluster",
                    //"DeleteCluster",
                    //"UpdateCluster",

                    "*Service",
                    //"CreateService",
                    //"DeleteService",
                    //"UpdateService",
                
                    "Describe*",
                    //"DescribeClusters",
                    //"DescribeServices",
                    //"DescribeTaskDefinition",
                    //"DescribeTasks",
        
                    "List*",
                    //"ListClusters",
                    //"ListServices",
                    //"ListTaskDefinitions",
                    //"ListTasks",
                    
                    "*Task",
                    //"RunTask",
                    //"StartTask",
                    //"StopTask",
                    
                    "ListTagsForResource",
                    "TagResource",
                    "UntagResource",
        
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "*", resourceName: "Agc*"}),
                    this.arn({service: svc, region: "*", resource: "service", resourceName: "wesAdapter"}),
                    this.arn({service: svc, region: "*", resource: "task", resourceName: "*"}),
                ]
            }),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "CreateCluster",
                    "RegisterTaskDefinition",
                    "DeregisterTaskDefinition",
                    "DescribeTaskDefinition",
                ),
                resources: [
                    "*"
                ]
            })
        ]
    }
    
    // elb
    elb() : PolicyStatement[] {
        const svc = "elasticloadbalancing";
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    
                    "Create*",
                    //"CreateListener",
                    //"CreateLoadBalancer",
                    //"CreateTargetGroup",
                    
                    "Delete*",
                    //"DeleteListener",
                    //"DeleteLoadBalancer",
                    //"DeleteTargetGroup",
                    
                    "ModifyLoadBalancerAttributes",
                    
                    "*Tags",
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "loadbalancer/net", resourceName: "Agc-*"}),
                    this.arn({service: svc, region: "*", resource: "listener/net", resourceName: "Agc-*"}),
                    this.arn({service: svc, region: "*", resource: "targetgroup", resourceName: "Agc-*"}),
                ]
            }),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "Describe*",
                    // "DescribeListeners", //*
                    // "DescribeLoadBalancers", //*
                    // "DescribeTags",
                    // "DescribeTargetGroups", //*
                ),
                resources: [
                    "*"
                ]
            })
    ]
    }
    
    
    // apigw
    // these permissions replicate the aws managed polices for apigateway
    // AmazonAPIGatewayInvokeFullAccess
    // AmazonAPIGatewayAdministrator
    apigw() : PolicyStatement[] {
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions("execute-api",
                    "Invoke",
                    "ManageConnections",
                ),
                resources: [
                    this.arn({service: "execute-api", region: "*", resource: "*"}),
                ]
            }),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions("apigateway", "*"),
                resources: [
                    this.arn({service: "apigateway", region: "*", account: "", resource: "*"}),
                ]
            })
        ]
    }
    
    
    // efs
    efs() : PolicyStatement[] {
        const svc = "elasticfilesystem";
        return [new PolicyStatement({
            effect: Effect.ALLOW,
            actions: actions(svc,
                "*FileSystem",
                // "CreateFileSystem",
                // "DeleteFilesystem",
                // "UpdateFileSystem",
                
                "*MountTarget",
                // "CreateMountTarget",
                // "DeleteMountTarget",
                
                "Describe*",
                //"DescribeFilesystems",
                //"DescribeMountTargets",
                //"DescribeTags",
                
                
                "*Tag*",
                //"CreateTags",
                //"DeleteTags",
                //"ListTagsForResource",
                //"TagResource",
                //"UntagResource",
    
            ),
            resources: [
                this.arn({service: svc, region: "*", resource: "file-system", resourceName: "*"}),
            ]
        })]
    }
    
    
    
    // service-discovery (aka cloudmap)
    cloudmap() : PolicyStatement[]{
        const svc = "servicediscovery";
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "*Service",
                    //"CreateService",
                    //"DeleteService",
                    
                    "*Namespace",
                    //"CreatePrivateDnsNamespace",
                    //"DeleteNamespace",
                    
                    //"Get*",
                    //"GetService",
                    //"GetNamespace",
                    //"GetOperation",
        
                    "List*",
                    //"ListNamespaces",
                    //"ListServices",
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "namespace", resourceName: "*"}),
                    this.arn({service: svc, region: "*", resource: "service", resourceName: "*"}),
                ]
            }),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "CreatePrivateDnsNamespace",
                    "GetOperation",
                    "*Tag*",
                    //"ListTagsForResouce",
                    //"TagResource",
                    //"UntagResource",
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "*", resourceName: "*"}),
                ]
            })
        ]
    }
    
    
    // logs
    logs() : PolicyStatement[] {
        const svc = "logs";
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "*LogGroup",
                    //"CreateLogGroup",
                    //"DeleteLogGroup",
                    //"ListTagsForLogGroup",
                    //"TagLogGroup",
                    //"UntagLogGroup",
                    
                    "*LogStream",
                    //"CreateLogStream",
                    //"DeleteLogStream",
                    
                    "*RetentionPolicy",
                    //"DeleteRetentionPolicy",
                    //"PutRetentionPolicy",

                    "Describe*",
                    //"DescribeLogGroups",
                    //"DescribeLogStreams",
        
                    "Filter*",
                    "Get*",
                    //"GetLogEvents",
                    //"GetLogGroupFields",
                    //"GetLogRecord",
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "log-group:agc*"}),
                    this.arn({service: svc, region: "*", resource: "log-group:Agc*"}),
                    this.arn({service: svc, region: "*", resource: "log-group:/aws/batch/job*"}),
                ]
            }),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "DescribeLogGroups",
                ),
                resources: [
                    this.arn({service: svc, region: "*", resource: "log-group:*"}),
                ]
            })
        ]
    }

    ecr() : PolicyStatement[] {
        const svc = "ecr";
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "ListImages",
                ),
                resources: [
                    this.arn({service: svc, region: "*", account: "*", resource: "*"})
                ]
            })
        ]
    }

    route53() :  PolicyStatement[] {
        const svc = "route53";
        return [
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: actions(svc,
                    "CreateHostedZone",
                    "ListHostedZonesByName",
                ),
                resources: [
                    "*",
                ]
            })
        ]
    }

    deactivate() : PolicyStatement[] {
        return [
            // most of this is for `context destroy` to enable `account deactivate --force`
            ...this.iam(),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: [
                    ...actions("apigateway",
                        "GET",
                        "DELETE"
                    ),
                    ...actions("elasticfilesystem",
                        "DeleteFileSystem",
                        "DeleteMountTarget",
                    ),
                    ...actions("ecs",
                        "Delete*",
                    ),
                    ...actions("ec2",
                        "*LaunchTemplate*",
                    ),
                    ...actions("elasticloadbalancing",
                        "Delete*",
                    ),
                    ...actions("servicediscovery",
                        "DeleteNamespace",
                        "DeleteService",
                        "GetOperation",
                    ),
                    ...actions("batch",
                        "Update*",
                        "Delete*",
                    ),
                ],
                resources: [
                    this.arn({service: "apigateway", account: "", resource: "/restapis*"}),
                    this.arn({service: "apigateway", account: "", resource: "/vpclinks*"}),
                    this.arn({service: "elasticfilesystem", resource: "file-system", resourceName: "*"}),
                    this.arn({service: "ecs", resource: "cluster", resourceName: "Agc*"}),
                    this.arn({service: "ecs", resource: "service", resourceName: "Agc*"}),
                    this.arn({service: "ec2", resource: "launch-template", resourceName: "*"}),
                    this.arn({service: "elasticloadbalancing", resource: "*", resourceName: "Agc*"}),
                    this.arn({service: "servicediscovery", resource: "*"}),
                    this.arn({service: "batch", resource: "job-queue", resourceName: "TaskBatch*"}),
                    this.arn({service: "batch", resource: "compute-environment", resourceName: "TaskBatch*"}),
                ]
            }),
            new PolicyStatement({
                effect: Effect.ALLOW,
                actions: [
                    ...actions("elasticfilesystem",
                        "Describe*",
                    ),
                    ...actions("ecs",
                        "Describe*",
                        "DeregisterTaskDefinition",
                    ),
                    ...actions("elasticloadbalancing",
                        "Describe*",
                    ),
                    ...actions("batch",
                        "Describe*",
                    ),
                    ...actions("ec2",
                        "DeleteVpcEndpointServiceConfigurations",
                        "DescribeVpcEndpointServiceConfigurations",
                    ),
                ],
                resources: [
                    "*"
                ]
            })
        ]
    }
}



