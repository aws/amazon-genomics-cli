# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## 1.1.0 (2021-11-11)


### Features

* context destroy --force flag ([#118](https://github.com/aws/amazon-genomics-cli/issues/118)) ([4282093](https://github.com/aws/amazon-genomics-cli/commit/428209311aa247c999816348a972737739b1189f))
* Initial infrastructure for MiniWdl support ([#125](https://github.com/aws/amazon-genomics-cli/issues/125)) ([387393a](https://github.com/aws/amazon-genomics-cli/commit/387393a64593d08ec2016a69382f7d64de37914e))
* Introducing AWS Lambda based WES Adapter for running the workflows ([#155](https://github.com/aws/amazon-genomics-cli/issues/155)) ([bd90f47](https://github.com/aws/amazon-genomics-cli/commit/bd90f47b94ede34c31ea109221225ff3cd65d200))


### Bug Fixes

* Adds a message when new logs aren't shown to the user immediately ([#131](https://github.com/aws/amazon-genomics-cli/issues/131)) ([54349d2](https://github.com/aws/amazon-genomics-cli/commit/54349d2858a837da26e6479c409e4a8445055562))
* Asserts order deterministically ([#153](https://github.com/aws/amazon-genomics-cli/issues/153)) ([c4299e8](https://github.com/aws/amazon-genomics-cli/commit/c4299e86e499edeb4695152f558a88e82bcf2da3))
* correctly link to core app ([#133](https://github.com/aws/amazon-genomics-cli/issues/133)) ([ce103b2](https://github.com/aws/amazon-genomics-cli/commit/ce103b202d50c7a8a40e6d94daca0c4dd5141da7))
* Deregionalize min permissions ([#128](https://github.com/aws/amazon-genomics-cli/issues/128)) ([c4dc96b](https://github.com/aws/amazon-genomics-cli/commit/c4dc96b1641431ed7c20fad348e7d87d2156a4b8))
* Fixes how users interact with the context commands ([#115](https://github.com/aws/amazon-genomics-cli/issues/115)) ([ffb3bb6](https://github.com/aws/amazon-genomics-cli/commit/ffb3bb6fdffeabd09a33288086c5442aa5e14c60))
* improve contrast in docs ([#149](https://github.com/aws/amazon-genomics-cli/issues/149)) ([beb10f4](https://github.com/aws/amazon-genomics-cli/commit/beb10f4b02f9533da13ce0b3579ae2fd55a337aa))
* invalid AWS Health url ([#130](https://github.com/aws/amazon-genomics-cli/issues/130)) ([1aef83b](https://github.com/aws/amazon-genomics-cli/commit/1aef83b682ba276ae5d8720ccaffc97a66bb34cb))
* miniwdl interpolation workaround ([27f4bf5](https://github.com/aws/amazon-genomics-cli/commit/27f4bf571712c6509e6352f4459e452fdd6a1cb1))
* Pass engine endpoint directly the wes adapter ([#122](https://github.com/aws/amazon-genomics-cli/issues/122)) ([81ed484](https://github.com/aws/amazon-genomics-cli/commit/81ed484a94ce195259315826377ece0443b582e1))
* show logs for workflows with more than 100 tasks ([#114](https://github.com/aws/amazon-genomics-cli/issues/114)) ([4e54c3b](https://github.com/aws/amazon-genomics-cli/commit/4e54c3bae5ad8242fb1af0ab171aeb4c5b818923))
* Shows the relevant error if the workflow logs can't be retrieved ([#103](https://github.com/aws/amazon-genomics-cli/issues/103)) ([2051b95](https://github.com/aws/amazon-genomics-cli/commit/2051b9542d07c5f999bd149e2a9f65aefaccba00))
* temporary folder potential leak in some error scenarios. unit test for cdk command execution ([#140](https://github.com/aws/amazon-genomics-cli/issues/140)) ([c50608e](https://github.com/aws/amazon-genomics-cli/commit/c50608e594b528a7bddd33b678da984feabc50b4))
* updates context describe to be consistent with context destroy ([#143](https://github.com/aws/amazon-genomics-cli/issues/143)) ([bb7667b](https://github.com/aws/amazon-genomics-cli/commit/bb7667b44027b3374b8011da11418d6ee0054b79))
* updates how the logs are shown from cloudwatch ([#142](https://github.com/aws/amazon-genomics-cli/issues/142)) ([1502578](https://github.com/aws/amazon-genomics-cli/commit/1502578415c7db4c5a633982301a887bcd393514))
* Use correct context name ([#132](https://github.com/aws/amazon-genomics-cli/issues/132)) ([c1516b6](https://github.com/aws/amazon-genomics-cli/commit/c1516b60b5706d06b30d7516a3aa1d80efd216af))
* use proper go tags for windows build ([#117](https://github.com/aws/amazon-genomics-cli/issues/117)) ([7dfa94a](https://github.com/aws/amazon-genomics-cli/commit/7dfa94a775fdba5193c99d0c697c8013a52a23ce))
* Workflow status now ignores unqueryable stacks ([#138](https://github.com/aws/amazon-genomics-cli/issues/138)) ([cf817a8](https://github.com/aws/amazon-genomics-cli/commit/cf817a882de2160d8e333c17d4eb28508cd886e1))
* workflows from demo-wdl-project should run without errors out of the box ([#108](https://github.com/aws/amazon-genomics-cli/issues/108)) ([6960eac](https://github.com/aws/amazon-genomics-cli/commit/6960eacf236e744d3c5658c5557061ab9cd3d468))

## 1.0.1 (2021-10-01)

### Bug Fixes

* Updated documentation
* Scoped down engine IAM permissions
* Improved error messages

## 1.0.0 (2021-09-26)

### Features

* First release!
