#!/bin/bash
set -e
set -x

INSTALL_VERSION=dist_release
echo "INSTALL_VERSION = $INSTALL_VERSION"
FILESYSTEM=btrfs
echo "FILESYSTEM = $FILESYSTEM"
INITIAL_SIZE=200
echo "INITIAL_SIZE = $INITIAL_SIZE"

USAGE=$(cat <<EOF
Retrieve and install Amazon EBS Autoscale

    $0 [options]

Options

    -i, --install-version       [dist_release] | release | latest | develop
            Version of Amazon EBS Autoscale to install.
            
                "dist_release" uses 'aws s3 cp' to retrieve a tarball from an S3 bucket.
                    requires setting --artifact-root-url to an S3 URL.

                "release" uses 'curl' or 'aws s3 cp' to retrieve a tarball from a publicly 
                    accessible location - i.e. an upstream distribution.
                    requires setting --artifact-root-url to either an S3 or HTTP URL.
                
                "latest" uses 'curl' to retrieve the latest released version of 
                    amazon-ebs-autoscale from GitHub
                
                "develop" uses 'git' to clone the source code of amazon-ebs-autoscale
                    from GitHub.
    
    -a, --artifact-root-url     s3://... | https:// ...
            Root URL where amazon-ebs-autoscale tarballs can be retrieved.
            Required if --install-version is "dist_release" or "release".
    
    -f, --file-system           [btrfs] | lvm.ext4
            File system to use
    
    -h, --help
            Print help and exit

EOF
)

PARAMS=""
while (( "$#" )); do
    case "$1" in
        -i|--install-version)
            INSTALL_VERSION=$2
            shift 2
            ;;
        -a|--artifact-root-url)
            ARTIFACT_ROOT_URL=$2
            shift 2
            ;;
        -f|--file-system)
            FILESYSTEM=$2
            shift 2
            ;;
        -s|--initial-size)
            INITIAL_SIZE=$2
            shift 2
            ;;
        -h|--help)
            echo "$USAGE"
            exit
            ;;
        --) # end parsing
            shift
            break
            ;;
        -*|--*=)
            error "unsupported argument $1"
            ;;
        *) # positional arguments
            PARAMS="$PARAMS $1"
            shift
            ;;
    esac
done

eval set -- "$PARAMS"


function develop() {
    echo "retrieving the current development version of amazon-ebs-autoscale WARNING may not be fully tested or stable"
    git clone https://github.com/awslabs/amazon-ebs-autoscale.git
    cd /opt/amazon-ebs-autoscale
    git checkout master
}

function latest() {
    # recommended if you want instances to stay up to date with upstream updates
    echo "retrieving the latest released version of amazon-ebs-autoscale"
    local ebs_autoscale_version
    ebs_autoscale_version=$(curl --silent --fail --retry 3 --retry-connrefused "https://api.github.com/repos/awslabs/amazon-ebs-autoscale/releases/latest" | jq -r .tag_name)
    echo "ebs_autoscale_version = $ebs_autoscale_version"
    curl --silent --fail --retry 3 --retry-connrefused  -L \
        "https://github.com/awslabs/amazon-ebs-autoscale/archive/${ebs_autoscale_version}.tar.gz" \
        -o ./amazon-ebs-autoscale.tar.gz 

    tar -xzvf ./amazon-ebs-autoscale.tar.gz
    mv ./amazon-ebs-autoscale*/ ./amazon-ebs-autoscale
    echo "$ebs_autoscale_version" > ./amazon-ebs-autoscale/VERSION
}

function s3CopyWithRetry() {
    local s3_path=$1
    # destination must be the path to a file and not just the directory you want the file in
    local destination=$2
    echo "copying $s3_path to $destination"

    for i in {1..5};
    do
        if [[ $s3_path =~ s3://([^/]+)/(.+) ]]; then
            bucket="${BASH_REMATCH[1]}"
            key="${BASH_REMATCH[2]}"
            content_length=$(aws s3api head-object --bucket "$bucket" --key "$key" --query 'ContentLength')
        else
            echo "$s3_path is not an S3 path with a bucket and key. aborting"
            exit 1
        fi
        
        aws s3 cp --no-progress "$s3_path" "$destination" &&
        [[ $(LC_ALL=C ls -dn -- "$destination" | awk '{print $5; exit}') -eq "$content_length" ]] &&
        break || echo "attempt $i to copy $s3_path failed";

        if [ "$i" -eq 5 ]; then
            echo "failed to copy $s3_path after $i attempts. aborting"
            exit 2
        fi
        sleep $((7 * i))
    done
    echo "copied $s3_path to $destination with return code $?"
}

function release() {
    echo "retrieving the version of amazon-ebs-autoscale from the latest upstream distribution release of aws-genomics-workflows"

    if [[ ! $ARTIFACT_ROOT_URL ]]; then
        echo "missing required argument: --artifact-root-url"
        exit 1
    fi
    echo "ARTIFACT_ROOT_URL = $ARTIFACT_ROOT_URL"

    if [[ "$ARTIFACT_ROOT_URL" =~ ^http.* ]]; then
        echo "obtaining amazon-ebs-autoscale.tgz using http protocol"
        curl --silent --fail --retry 3 --retry-connrefused -LO $ARTIFACT_ROOT_URL/amazon-ebs-autoscale.tgz
    elif [[ "$ARTIFACT_ROOT_URL" =~ ^s3.* ]]; then
        echo "obtaining amazon-ebs-autoscale.tgz using s3 protocol"
        s3CopyWithRetry "$ARTIFACT_ROOT_URL"/amazon-ebs-autoscale.tgz ./amazon-ebs-autoscale.tgz
    else
        echo "unrecognized protocol in $ARTIFACT_ROOT_URL for release()"
        exit 1
    fi

    echo "unpacking amazon-ebs-autoscale.tgz"
    tar -xzvf amazon-ebs-autoscale.tgz
}

function dist_release() {
    echo "retrieve the release version of amazon-ebs-autoscale installed as an artifact with the AGC Core stack."
    # recommended for a fully self-contained deployment

    if [[ ! $ARTIFACT_ROOT_URL ]]; then
        echo "missing required argument: --artifact-root-url"
        exit 1
    fi
    echo "ARTIFACT_ROOT_URL = $ARTIFACT_ROOT_URL"

    if [[ "$ARTIFACT_ROOT_URL" =~ ^s3.* ]]; then
        echo "obtaining amazon-ebs-autoscale.tgz using s3 protocol"
        s3CopyWithRetry $ARTIFACT_ROOT_URL/amazon-ebs-autoscale.tgz ./amazon-ebs-autoscale.tgz
    else
        echo "unrecognized protocol in $ARTIFACT_ROOT_URL for dist_release()"
        exit 1
    fi

    echo "unpacking amazon-ebs-autoscale.tgz"
    tar -xzvf amazon-ebs-autoscale.tgz
}

function install() {
    local filesystem=${1:-btrfs}
    echo "ebs autoscale filesystem = $filesystem"
    local docker_storage_driver

    case $filesystem in
        btrfs)
            docker_storage_driver=$filesystem
            ;;
        lvm.ext4)
            docker_storage_driver=overlay2
            ;;
        *)
            echo "Unsupported filesystem - $filesystem"
            exit 1
    esac
    echo "docker_storage_driver = $docker_storage_driver"

    local docker_storage_options="DOCKER_STORAGE_OPTIONS=\"--storage-driver $docker_storage_driver\""
    
    cp -au /var/lib/docker /var/lib/docker.bk
    rm -rf /var/lib/docker/*

    echo "installing ebs autoscale"
    sh /opt/amazon-ebs-autoscale/install.sh -d /dev/xvdba -f "$filesystem" -m /var/lib/docker > /var/log/ebs-autoscale-install.log 2>&1
    echo "/opt/amazon-ebs-autoscale/install.sh exited with return code $?"

    awk -v docker_storage_options="$docker_storage_options" \
        '{ sub(/DOCKER_STORAGE_OPTIONS=.*/, docker_storage_options); print }' \
        /etc/sysconfig/docker-storage \
        > /opt/amazon-ebs-autoscale/docker-storage
    mv -f /opt/amazon-ebs-autoscale/docker-storage /etc/sysconfig/docker-storage

    cp -au /var/lib/docker.bk/* /var/lib/docker

}

mkdir -p /opt
cd /opt
$INSTALL_VERSION

install "$FILESYSTEM"