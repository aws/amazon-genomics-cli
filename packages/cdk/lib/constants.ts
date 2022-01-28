export const PRODUCT_NAME = "Agc";
export const APP_NAME = "agc";
export const APP_ENV_NAME = "AGC";
export const APP_TAG_KEY = "application-name";
export const PROJECT_TAG_KEY = `${APP_NAME}-project`;
export const CONTEXT_TAG_KEY = `${APP_NAME}-context`;
export const USER_ID_TAG_KEY = `${APP_NAME}-user-id`;
export const USER_EMAIL_TAG_KEY = `${APP_NAME}-user-email`;
export const AGC_VERSION_KEY = `${APP_NAME}-version`;
export const VPC_PARAMETER_NAME = "vpc";

export const LAUNCH_TEMPLATE = `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="==MYBOUNDARY=="

--==MYBOUNDARY==
Content-Type: text/cloud-config; charset="us-ascii"

packages:
- jq
- grep
- btrfs-progs
- sed
- git
- unzip
- amazon-cloudwatch-agent

write_files:
- permissions: '0644'
  path: /opt/aws/amazon-cloudwatch-agent/etc/config.json
  content: |
    {
      "agent": {
        "logfile": "/opt/aws/amazon-cloudwatch-agent/logs/amazon-cloudwatch-agent.log"
      },
      "logs": {
        "logs_collected": {
          "files": {
            "collect_list": [
              {
                "file_path": "/opt/aws/amazon-cloudwatch-agent/logs/amazon-cloudwatch-agent.log",
                "log_group_name": "/aws/ecs/container-instance/${APP_NAME}",
                "log_stream_name": "/aws/ecs/container-instance/${APP_NAME}/{instance_id}/amazon-cloudwatch-agent.log"
              },
              {
                "file_path": "/var/log/cloud-init.log",
                "log_group_name": "/aws/ecs/container-instance/${APP_NAME}",
                "log_stream_name": "/aws/ecs/container-instance/${APP_NAME}/{instance_id}/cloud-init.log"
              },
              {
                "file_path": "/var/log/cloud-init-output.log",
                "log_group_name": "/aws/ecs/container-instance/${APP_NAME}",
                "log_stream_name": "/aws/ecs/container-instance/${APP_NAME}/{instance_id}/cloud-init-output.log"
              },
              {
                "file_path": "/var/log/ecs/ecs-init.log",
                "log_group_name": "/aws/ecs/container-instance/${APP_NAME}",
                "log_stream_name": "/aws/ecs/container-instance/${APP_NAME}/{instance_id}/ecs-init.log"
              },
              {
                "file_path": "/var/log/ecs/ecs-agent.log",
                "log_group_name": "/aws/ecs/container-instance/${APP_NAME}",
                "log_stream_name": "/aws/ecs/container-instance/${APP_NAME}/{instance_id}/ecs-agent.log"
              },
              {
                "file_path": "/var/log/ecs/ecs-volume-plugin.log",
                "log_group_name": "/aws/ecs/container-instance/${APP_NAME}",
                "log_stream_name": "/aws/ecs/container-instance/${APP_NAME}/{instance_id}/ecs-volume-plugin.log"
              }
            ]
          }
        }
      }
    }

runcmd:

# start the amazon-cloudwatch-agent
- /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:/opt/aws/amazon-cloudwatch-agent/etc/config.json
- /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -m ec2 -a status | jq -r '.status' | grep -iw "running" || sleep 5 && /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:/opt/aws/amazon-cloudwatch-agent/etc/config.json
- /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -m ec2 -a status | jq -r '.status' | grep -iw "running" || sleep 10 && /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:/opt/aws/amazon-cloudwatch-agent/etc/config.json
- /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -m ec2 -a status | jq -r '.status' | grep -iw "running" || shutdown -P now

# install aws-cli v2 and copy the static binary in an easy to find location for bind-mounts into containers
- mkdir -p /opt/aws-cli/bin
- curl -s "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip" && unzip -q /tmp/awscliv2.zip -d /tmp && /tmp/aws/install -b /usr/bin && cp -a -f $(dirname $(find /usr/local/aws-cli -name 'aws' -type f))/. /opt/aws-cli/bin/
- command -v aws || sleep 5 | curl -s "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip" && unzip -q /tmp/awscliv2.zip -d /tmp && /tmp/aws/install -b /usr/bin && cp -a -f $(dirname $(find /usr/local/aws-cli -name 'aws' -type f))/. /opt/aws-cli/bin/
- command -v aws || sleep 10 | curl -s "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip" && unzip -q /tmp/awscliv2.zip -d /tmp && /tmp/aws/install -b /usr/bin && cp -a -f $(dirname $(find /usr/local/aws-cli -name 'aws' -type f))/. /opt/aws-cli/bin/
- command -v aws || shutdown -P now

# set environment variables for provisioning
- export ARTIFACTS_NAMESPACE=${APP_NAME}
- export INSTALLED_ARTIFACTS_S3_ROOT_URL=$(aws ssm get-parameter --name /\${ARTIFACTS_NAMESPACE}/_common/installed-artifacts/s3-root-url --query 'Parameter.Value' --output text)

# enable ecs spot instance draining
- echo ECS_ENABLE_SPOT_INSTANCE_DRAINING=true >> /etc/ecs/ecs.config

# pull docker images only if missing
- echo ECS_IMAGE_PULL_BEHAVIOR=prefer-cached >> /etc/ecs/ecs.config

# Setup ecs additions
- cd /opt
- aws s3 sync \${INSTALLED_ARTIFACTS_S3_ROOT_URL}/ecs-additions/ ./ecs-additions && chmod a+x /opt/ecs-additions/provision.sh  
- test -f ./ecs-additions/fetch_and_run.sh || sleep 5 || aws s3 sync \${INSTALLED_ARTIFACTS_S3_ROOT_URL}/ecs-additions/ ./ecs-additions && chmod a+x /opt/ecs-additions/provision.sh    
- test -f ./ecs-additions/fetch_and_run.sh || sleep 10 || aws s3 sync \${INSTALLED_ARTIFACTS_S3_ROOT_URL}/ecs-additions/ ./ecs-additions && chmod a+x /opt/ecs-additions/provision.sh    
- test -f ./ecs-additions/fetch_and_run.sh || shutdown -P now
- /opt/ecs-additions/provision.sh

- echo "successfully initiated"
--==MYBOUNDARY==--
`;
