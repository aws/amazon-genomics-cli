# Adding New Engine Base Images
This document specifies how to create images that run on AWS Batch for new genomics engines.

## Dockerfile
The Dockerfile specifies the instructions for docker to build the image. There is an area of this file marked with ##### MODIFY ####### and #### END MODIFY ######. 
There are 3 things that must be completed in this section. 
1. Provide installation instructions for the engine as well as any dependencies. You can refer to the other engines for examples. 
1. Update the script name that is to be copied. This will be the same name as the script in the next section. 
1. Update the `$PATH` with the new engine.

## \<engine-name\>.aws.sh
This is the entrypoint script which will execute the workflows. There are again ##### MODIFY ####### and #### END MODIFY ###### marked areas to customize this entroypoint script 
to work with the new engine.

## buildspec.yml
This file is used to build the image to the amazon-genomics-cli AWS ECR repository for distribution. Update all variables marked with `<>`

## THIRD-PARTY
This file must be updated with the third-party open-source license for the new engine being added. 