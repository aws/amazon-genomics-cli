# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.6.0](https://github.com/aws/amazon-genomics-cli/compare/v1.5.2...v1.6.0) (2022-12-28)


### Features

* Add ARM support for CLI build ([#569](https://github.com/aws/amazon-genomics-cli/issues/569)) ([f002a3d](https://github.com/aws/amazon-genomics-cli/commit/f002a3d8e1856842a3b8bc915ba3eb1e1365e6d8))
* added verbose logging to `workflow run` ([#480](https://github.com/aws/amazon-genomics-cli/issues/480)) ([5bec60a](https://github.com/aws/amazon-genomics-cli/commit/5bec60ae95c5c3c857b3aa2db212c6e2034e5d61))
* reduce permissions related to APIGateway for recommended user role ([#538](https://github.com/aws/amazon-genomics-cli/issues/538)) ([fd9e230](https://github.com/aws/amazon-genomics-cli/commit/fd9e230e993010836baeed3fe0e3280116f6d03f))
* Support and use Toil 5.8's WES task output ([#501](https://github.com/aws/amazon-genomics-cli/issues/501)) ([a5e52b5](https://github.com/aws/amazon-genomics-cli/commit/a5e52b5b316205290016197eb8163eef3cd864b1))
* support for PRIVATE API Gateway endpoints ([#507](https://github.com/aws/amazon-genomics-cli/issues/507)) ([357d545](https://github.com/aws/amazon-genomics-cli/commit/357d54590b465aa11b59a6d4b13a4f2609f37765))


### Bug Fixes

* bump nextflow version to 22.04.1 ([#450](https://github.com/aws/amazon-genomics-cli/issues/450)) ([ec9c02c](https://github.com/aws/amazon-genomics-cli/commit/ec9c02cd8cdb9d1fcffa31e26844a00ef8260efd)), closes [#387](https://github.com/aws/amazon-genomics-cli/issues/387)
* Install script output name ([def4335](https://github.com/aws/amazon-genomics-cli/commit/def4335c2193150ab06cc8385f69b45810cc39f5))
* update go version to fix bug with `go test -race` in OSX Monterey ([#505](https://github.com/aws/amazon-genomics-cli/issues/505)) ([c918e97](https://github.com/aws/amazon-genomics-cli/commit/c918e97fe48d9a23daaaaef42a5f3a1775dc4c72))
* update nextflow version environment variable ([b8ea514](https://github.com/aws/amazon-genomics-cli/commit/b8ea514213452151ef8b38458503bc044dfa6433))

## [1.5.2](https://github.com/aws/amazon-genomics-cli/compare/v1.5.1...v1.5.2) (2022-10-18)

### Security fixes

* Updates CDK and VM2 versions to include the latest releases to fix the following dependabot issues ([#10](https://github.com/aws/amazon-genomics-cli/security/dependabot/10)) ([#11](https://github.com/aws/amazon-genomics-cli/security/dependabot/11))

## [1.5.1](https://github.com/aws/amazon-genomics-cli/compare/v1.5.0...v1.5.1) (2022-09-08)


### Bug Fixes

* updated CDK packages ([#511](https://github.com/aws/amazon-genomics-cli/pull/511)) ([1140da1](https://github.com/aws/amazon-genomics-cli/commit/1140da166f0f2a9512ef1fb2ce022105d3afa9d8))
* bump nextflow version to 22.04.1 ([#450](https://github.com/aws/amazon-genomics-cli/issues/450)) ([4296ee7](https://github.com/aws/amazon-genomics-cli/commit/4296ee7316620cdf71bf14786433ae41bafff999)), closes [#387](https://github.com/aws/amazon-genomics-cli/issues/387)

## [1.5.0](https://github.com/aws/amazon-genomics-cli/compare/v1.5.0-rc1...v1.5.0) (2022-05-25)


### Bug Fixes

* enforce GO111MODULE=on for the build ([#457](https://github.com/aws/amazon-genomics-cli/issues/457)) ([412eff7](https://github.com/aws/amazon-genomics-cli/commit/412eff7dead0f8624fb2eaa7ad70efeedba63859))

## [1.4.0](https://github.com/aws/amazon-genomics-cli/compare/v1.3.0...v1.4.0) (2022-04-15)


### Features

* add "public subnet" mode for lower ongoing costs ([#283](https://github.com/aws/amazon-genomics-cli/issues/283)) ([d6a67f9](https://github.com/aws/amazon-genomics-cli/commit/d6a67f9ba819e6d2894e2dfce3d948d63439e556))
* Add optionsFile flag ([#352](https://github.com/aws/amazon-genomics-cli/issues/352)) ([884308f](https://github.com/aws/amazon-genomics-cli/commit/884308fc2b7e00d918c57e8a912957076309dde2))
* allow engine logs of specific cromwell workflows using run-id ([#389](https://github.com/aws/amazon-genomics-cli/issues/389)) ([ef8c6a3](https://github.com/aws/amazon-genomics-cli/commit/ef8c6a392c81c49b373766f2dcb9ea87aa5d1474))
* Download and report WES run log files ([#319](https://github.com/aws/amazon-genomics-cli/issues/319)) ([2a8fab0](https://github.com/aws/amazon-genomics-cli/commit/2a8fab00ceedfd72321f5466183edf92f4c4d0fa))
* Engine logs for specific workflow runs ([#379](https://github.com/aws/amazon-genomics-cli/issues/379)) ([8bda674](https://github.com/aws/amazon-genomics-cli/commit/8bda67492b3cf143ffd5a0a8931bbcdf3e04992d))
* ensures vpcId is recalled (if available) on account_activate ([#365](https://github.com/aws/amazon-genomics-cli/issues/365)) ([f056026](https://github.com/aws/amazon-genomics-cli/commit/f0560267fbbffaa17f72a0764950ec86d005b63d))
* Record command run and show in agc workflow status ([#353](https://github.com/aws/amazon-genomics-cli/issues/353)) ([45fac3a](https://github.com/aws/amazon-genomics-cli/commit/45fac3a476b95a94e65f1e8938a5a5750a703b0d))


### Bug Fixes

* Add safety checks in optional manifest fields ([#403](https://github.com/aws/amazon-genomics-cli/issues/403)) ([c084782](https://github.com/aws/amazon-genomics-cli/commit/c08478241cc09d218b3ee1e8a8593172cb2094c3))
* Converts WDL log exit_code from string to int32 per WES spec ([#380](https://github.com/aws/amazon-genomics-cli/issues/380)) ([2b93a2c](https://github.com/aws/amazon-genomics-cli/commit/2b93a2cb843bce2945538d06f487886b7efe6707))
* do not set WorkflowEngineParameters when empty ([#402](https://github.com/aws/amazon-genomics-cli/issues/402)) ([3260d73](https://github.com/aws/amazon-genomics-cli/commit/3260d73e7f43f0556a3adb88086bf58e67c43090))
* ensure workflow engines can and do terminate child tasks when the engine task is stopped by AGC ([#391](https://github.com/aws/amazon-genomics-cli/issues/391)) ([68861e1](https://github.com/aws/amazon-genomics-cli/commit/68861e13dd42b077a8da371a12d39a9c3ef4bc71))
* miniwdl wdl1.1support ([#370](https://github.com/aws/amazon-genomics-cli/issues/370)) ([7bcfd8f](https://github.com/aws/amazon-genomics-cli/commit/7bcfd8f10ee014c407ecbb801a1d4e5ee93c2496))
* nextflow wes adapter can retrieve more than 100 task logs ([#377](https://github.com/aws/amazon-genomics-cli/issues/377)) ([a1c2d30](https://github.com/aws/amazon-genomics-cli/commit/a1c2d30784c6ca5f81b53993edb21bc46aff2aa3))
* remove temp directory prefix from input URI ([#405](https://github.com/aws/amazon-genomics-cli/issues/405)) ([3933d7a](https://github.com/aws/amazon-genomics-cli/commit/3933d7a1b6a332008e9ad352946ff519d47e739c))
* Revert change to WDL version for this workflow ([#378](https://github.com/aws/amazon-genomics-cli/issues/378)) ([75603a8](https://github.com/aws/amazon-genomics-cli/commit/75603a860ca401b7be2d5d04c51248a227283d00))
* updates cromwell exit_code ([#398](https://github.com/aws/amazon-genomics-cli/issues/398)) ([9f6480d](https://github.com/aws/amazon-genomics-cli/commit/9f6480dcbd3bd280c822cce960aec0c6aceb7615))
* updates workflow engine parameters for key, val pairing ([#407](https://github.com/aws/amazon-genomics-cli/issues/407)) ([7a27122](https://github.com/aws/amazon-genomics-cli/commit/7a2712245521d7778614cbf6c08d83008e9471ed))

## [1.3.0](https://github.com/aws/amazon-genomics-cli/compare/v1.1.2...v1.3.0) (2022-03-04)


### Features

* add a global --silent flag ([#219](https://github.com/aws/amazon-genomics-cli/issues/219)) ([eed4189](https://github.com/aws/amazon-genomics-cli/commit/eed4189b2d988d88ef87be63b8f39610be569291))
* add a global --silent flag ([#219](https://github.com/aws/amazon-genomics-cli/issues/219)) ([#285](https://github.com/aws/amazon-genomics-cli/issues/285)) ([388533b](https://github.com/aws/amazon-genomics-cli/commit/388533b40d7ad4901fe2ef20649743ae21887fb1))
* Add AGC version tag to AWS resources ([#243](https://github.com/aws/amazon-genomics-cli/issues/243)) ([fc6cad7](https://github.com/aws/amazon-genomics-cli/commit/fc6cad707eb4c90098d3c9435b23efccf3ac8a72))
* Add JSON as a CLI output format ([#215](https://github.com/aws/amazon-genomics-cli/issues/215)) ([#309](https://github.com/aws/amazon-genomics-cli/issues/309)) ([2d6ad0c](https://github.com/aws/amazon-genomics-cli/commit/2d6ad0c60eab860b9a83077c7988a3dbb3858aa8))
* Add snakemake wes adapter + cdk + example workflow ([#304](https://github.com/aws/amazon-genomics-cli/issues/304)) ([6357b47](https://github.com/aws/amazon-genomics-cli/commit/6357b478522f718b6fcb44a56db142cb3edbb4db))
* Add WES API endpoint information to agc context describe output ([#253](https://github.com/aws/amazon-genomics-cli/issues/253)) ([6145db3](https://github.com/aws/amazon-genomics-cli/commit/6145db387344e8001e132320dbc404fecee38124))
* Added output file creation snakemake + copying efs files to s3 ([#316](https://github.com/aws/amazon-genomics-cli/issues/316)) ([bb4a42f](https://github.com/aws/amazon-genomics-cli/commit/bb4a42f3f9558d2bfe002a5947af913e562ea1cc))
* Bootstrap CDK during Account Activation ([#272](https://github.com/aws/amazon-genomics-cli/issues/272)) ([5228492](https://github.com/aws/amazon-genomics-cli/commit/5228492fba52aa9f6c075a9b27a293ca2298969a))
* Bootstrap CDK during Account Activation ([#272](https://github.com/aws/amazon-genomics-cli/issues/272)) ([219596e](https://github.com/aws/amazon-genomics-cli/commit/219596e8edf1666ca46a065bfa1ccc74848180f5))
* context destroy --force flag ([#118](https://github.com/aws/amazon-genomics-cli/issues/118)) ([4282093](https://github.com/aws/amazon-genomics-cli/commit/428209311aa247c999816348a972737739b1189f))
* don't use ebs auto scale for minidl ([#300](https://github.com/aws/amazon-genomics-cli/issues/300)) ([225d139](https://github.com/aws/amazon-genomics-cli/commit/225d139453fe233e93addfea257be13ab7f56952))
* Enable hosting assets in AWS commercial regions ([#240](https://github.com/aws/amazon-genomics-cli/issues/240)) ([2eb5007](https://github.com/aws/amazon-genomics-cli/commit/2eb5007ed304847f03c75ac16f7a58acc5e8ace1))
* Explain a bit more about S3 errors ([#318](https://github.com/aws/amazon-genomics-cli/issues/318)) ([fa41fa8](https://github.com/aws/amazon-genomics-cli/commit/fa41fa807da8177a0d5c5055a9541ca9a37f5ac7))
* Implement EFS IOPs specification in Project YAML ([#313](https://github.com/aws/amazon-genomics-cli/issues/313)) ([c1af22b](https://github.com/aws/amazon-genomics-cli/commit/c1af22bbf4c90251518a5d204ae1e1ecfa9209e5))
* Improve context stack deployment speed ([#150](https://github.com/aws/amazon-genomics-cli/issues/150)) ([f7d6c7f](https://github.com/aws/amazon-genomics-cli/commit/f7d6c7f4e5e0fb91b3851bf806bc411b7f4526e5))
* Improve inputs JSON handling via --inputsFile flag ([#335](https://github.com/aws/amazon-genomics-cli/issues/335)) ([#336](https://github.com/aws/amazon-genomics-cli/issues/336)) ([2f24c28](https://github.com/aws/amazon-genomics-cli/commit/2f24c28ae0f9674158750d3336911d45dde37c78))
* improve reliability of provision scripts ([#308](https://github.com/aws/amazon-genomics-cli/issues/308)) ([16eff36](https://github.com/aws/amazon-genomics-cli/commit/16eff364d56132712ec5fb67f5c16e9996acc7c2))
* Improved Workflow logs ([#156](https://github.com/aws/amazon-genomics-cli/issues/156)) ([f71231c](https://github.com/aws/amazon-genomics-cli/commit/f71231c556c7ecb92940f3685f55c7078b2d1028))
* Initial infrastructure for MiniWdl support ([#125](https://github.com/aws/amazon-genomics-cli/issues/125)) ([387393a](https://github.com/aws/amazon-genomics-cli/commit/387393a64593d08ec2016a69382f7d64de37914e))
* Introducing AWS Lambda based WES Adapter for running the workflows ([#155](https://github.com/aws/amazon-genomics-cli/issues/155)) ([bd90f47](https://github.com/aws/amazon-genomics-cli/commit/bd90f47b94ede34c31ea109221225ff3cd65d200))
* release extras along with the agc package ([b83edd4](https://github.com/aws/amazon-genomics-cli/commit/b83edd4a72de32db707bfd92eb60f80071eee654))
* release extras along with the agc package ([d4891e7](https://github.com/aws/amazon-genomics-cli/commit/d4891e7d5b1271db179f2365c4913324ba17a307))
* release extras along with the agc package ([b350ca4](https://github.com/aws/amazon-genomics-cli/commit/b350ca419d5518a8ed7d2e13002bc59fc06b86b3))
* stop discarding un-actionable deployment errors ([#270](https://github.com/aws/amazon-genomics-cli/issues/270)) ([2c7fbcb](https://github.com/aws/amazon-genomics-cli/commit/2c7fbcbe321f9e1d86c59455bf5c2c9467ce788c))
* Tag infrastructure with custom tags to prevent automatic deletion ([#254](https://github.com/aws/amazon-genomics-cli/issues/254)) ([#282](https://github.com/aws/amazon-genomics-cli/issues/282)) ([f4d7aaf](https://github.com/aws/amazon-genomics-cli/commit/f4d7aaf902701ac6ba17f827851ba4c98de51741))


### Bug Fixes

* Add AGC assets bucket to project read list ([#247](https://github.com/aws/amazon-genomics-cli/issues/247)) ([2019b6a](https://github.com/aws/amazon-genomics-cli/commit/2019b6ac0226276337ee4e93237c7cfb19c4d28e))
* add build wes adapter step to run dev script and readme ([#307](https://github.com/aws/amazon-genomics-cli/issues/307)) ([466ea69](https://github.com/aws/amazon-genomics-cli/commit/466ea69fa8076dd0db348490ccac81f6b5a7180f))
* Adds a message when new logs aren't shown to the user immediately ([#131](https://github.com/aws/amazon-genomics-cli/issues/131)) ([54349d2](https://github.com/aws/amazon-genomics-cli/commit/54349d2858a837da26e6479c409e4a8445055562))
* Allow home directory path in workflow project config ([#137](https://github.com/aws/amazon-genomics-cli/issues/137)) ([760ed3c](https://github.com/aws/amazon-genomics-cli/commit/760ed3cec26daacc1f7dd5f328a63b119aa98886))
* Asserts order deterministically ([#153](https://github.com/aws/amazon-genomics-cli/issues/153)) ([c4299e8](https://github.com/aws/amazon-genomics-cli/commit/c4299e86e499edeb4695152f558a88e82bcf2da3))
* container validation check with the correct region from user profile ([#310](https://github.com/aws/amazon-genomics-cli/issues/310)) ([c36a1f0](https://github.com/aws/amazon-genomics-cli/commit/c36a1f0ee052aea0cef721227c79a4b1370828fc))
* correctly handle missing return code from Cromwell executions ([#276](https://github.com/aws/amazon-genomics-cli/issues/276)) ([3c61735](https://github.com/aws/amazon-genomics-cli/commit/3c6173525b6e1b991c9ba52258822ad37fb024bc))
* correctly link to core app ([#133](https://github.com/aws/amazon-genomics-cli/issues/133)) ([ce103b2](https://github.com/aws/amazon-genomics-cli/commit/ce103b202d50c7a8a40e6d94daca0c4dd5141da7))
* Delete bootstrap bucket on account deactivate ([#306](https://github.com/aws/amazon-genomics-cli/issues/306)) ([a81710b](https://github.com/aws/amazon-genomics-cli/commit/a81710b1298648cae6a040065d1403d373cb0ba6))
* deploy wes adapter code as cdk asset ([#321](https://github.com/aws/amazon-genomics-cli/issues/321)) ([cbaa38c](https://github.com/aws/amazon-genomics-cli/commit/cbaa38c83115f8f4b10837a2bf175b23a324712a))
* Deregionalize min permissions ([#128](https://github.com/aws/amazon-genomics-cli/issues/128)) ([c4dc96b](https://github.com/aws/amazon-genomics-cli/commit/c4dc96b1641431ed7c20fad348e7d87d2156a4b8))
* Fix project data path prefix and docs ([#315](https://github.com/aws/amazon-genomics-cli/issues/315)) ([69a6bb4](https://github.com/aws/amazon-genomics-cli/commit/69a6bb4faca4b2ec3839905b02bbf79289dc81e4))
* fixed arg validation for format ([#198](https://github.com/aws/amazon-genomics-cli/issues/198)) ([0b68019](https://github.com/aws/amazon-genomics-cli/commit/0b6801942dd2fea22bcf8291833732f0c337d0e9))
* Fixes how inputs are processed in AGC via Manifest ([#277](https://github.com/aws/amazon-genomics-cli/issues/277)) ([be5c684](https://github.com/aws/amazon-genomics-cli/commit/be5c684ca0967baeda7cb8055d24d39b54cd3084))
* Fixes how users interact with the context commands ([#115](https://github.com/aws/amazon-genomics-cli/issues/115)) ([ffb3bb6](https://github.com/aws/amazon-genomics-cli/commit/ffb3bb6fdffeabd09a33288086c5442aa5e14c60))
* Fixes the mechanism for copying the directory contents ([#311](https://github.com/aws/amazon-genomics-cli/issues/311)) ([29f8e9e](https://github.com/aws/amazon-genomics-cli/commit/29f8e9ee012e09ca294547ea3dba225304b36bfa))
* fixing monocdk imports ([#190](https://github.com/aws/amazon-genomics-cli/issues/190)) ([9c94fd8](https://github.com/aws/amazon-genomics-cli/commit/9c94fd87bcd860fe3842536dbcedd2165fdcc80e))
* force use node v14 in cdk package ([d67a18c](https://github.com/aws/amazon-genomics-cli/commit/d67a18cdb31387642b015daeab44ac80a5972c14))
* function for default config ([#201](https://github.com/aws/amazon-genomics-cli/issues/201)) ([e6983ae](https://github.com/aws/amazon-genomics-cli/commit/e6983ae774f9d85cca634d4f9a0c04098ae6a234))
* Handle errors from CDK command ([#245](https://github.com/aws/amazon-genomics-cli/issues/245)) ([99d9378](https://github.com/aws/amazon-genomics-cli/commit/99d9378a3f5ad3f408ef24ac3dae42b46b79bbd5))
* improve contrast in docs ([#149](https://github.com/aws/amazon-genomics-cli/issues/149)) ([beb10f4](https://github.com/aws/amazon-genomics-cli/commit/beb10f4b02f9533da13ce0b3579ae2fd55a337aa))
* Increasing WES adapter lambda timeout ([#180](https://github.com/aws/amazon-genomics-cli/issues/180)) ([20cd77d](https://github.com/aws/amazon-genomics-cli/commit/20cd77dac1def4414f0159509caa8dfe853d62bb))
* invalid AWS Health url ([#130](https://github.com/aws/amazon-genomics-cli/issues/130)) ([1aef83b](https://github.com/aws/amazon-genomics-cli/commit/1aef83b682ba276ae5d8720ccaffc97a66bb34cb))
* miniwdl interpolation workaround ([27f4bf5](https://github.com/aws/amazon-genomics-cli/commit/27f4bf571712c6509e6352f4459e452fdd6a1cb1))
* nest batch artifacts and disable bucket clean up ([#262](https://github.com/aws/amazon-genomics-cli/issues/262)) ([73bb98d](https://github.com/aws/amazon-genomics-cli/commit/73bb98d23a18f210d4bfa231b0f8002bb294529c))
* Nextflow workflow logs unmarshal number ([#286](https://github.com/aws/amazon-genomics-cli/issues/286)) ([1ddf67c](https://github.com/aws/amazon-genomics-cli/commit/1ddf67c50e53e64c70acd003616a0f6e2f3af716))
* Nextflow workflow logs unmarshal number ([#286](https://github.com/aws/amazon-genomics-cli/issues/286)) ([f739f03](https://github.com/aws/amazon-genomics-cli/commit/f739f03b8825084a0f679ceb7b83a3db87b78c2e))
* Pass engine endpoint directly the wes adapter ([#122](https://github.com/aws/amazon-genomics-cli/issues/122)) ([81ed484](https://github.com/aws/amazon-genomics-cli/commit/81ed484a94ce195259315826377ece0443b582e1))
* Pinned running of rnaseq to version 3.4 ([5386a3c](https://github.com/aws/amazon-genomics-cli/commit/5386a3c4afd619989b61ad853d4c23aeecbe154e))
* progress bar doesn't show any progress ([#166](https://github.com/aws/amazon-genomics-cli/issues/166)) ([9dd17ae](https://github.com/aws/amazon-genomics-cli/commit/9dd17aef9e5edee719ed1c78f9d09aa3a8b4f9c8))
* Remove Chdir when looking for project files ([#228](https://github.com/aws/amazon-genomics-cli/issues/228)) ([3564d6f](https://github.com/aws/amazon-genomics-cli/commit/3564d6f5a2433191c3c31456f73d9b65f2cf1a87))
* rename unassigned variable ([#329](https://github.com/aws/amazon-genomics-cli/issues/329)) ([a3c26a8](https://github.com/aws/amazon-genomics-cli/commit/a3c26a805674e717901569123df3c44683b4cdab))
* Respect maxVCpus in miniwdl contexts ([#202](https://github.com/aws/amazon-genomics-cli/issues/202)) ([e4ad17e](https://github.com/aws/amazon-genomics-cli/commit/e4ad17efa9a3724c7dd9bc2f37430c21ebcf9d4e))
* Return timezone information in RunLog ([#174](https://github.com/aws/amazon-genomics-cli/issues/174)) ([118e1e2](https://github.com/aws/amazon-genomics-cli/commit/118e1e24ce8b6c454d641f27129bb04d6a44cf36))
* Revert add a global --silent flag ([#219](https://github.com/aws/amazon-genomics-cli/issues/219)) ([#279](https://github.com/aws/amazon-genomics-cli/issues/279)) ([a6cc990](https://github.com/aws/amazon-genomics-cli/commit/a6cc9904cd2d1899f7f64bf3c024a242c1a33749)), closes [#274](https://github.com/aws/amazon-genomics-cli/issues/274)
* Revert shrinkwrap upgrade ([#298](https://github.com/aws/amazon-genomics-cli/issues/298)) ([9dd6806](https://github.com/aws/amazon-genomics-cli/commit/9dd68060e044e4af487f0f67dc9d3da2a8439e8b))
* Scope down batch role permissions ([#230](https://github.com/aws/amazon-genomics-cli/issues/230)) ([14a5aca](https://github.com/aws/amazon-genomics-cli/commit/14a5acaf3f931cbaf53c0a495034cbff462797b0))
* show logs for workflows with more than 100 tasks ([#114](https://github.com/aws/amazon-genomics-cli/issues/114)) ([4e54c3b](https://github.com/aws/amazon-genomics-cli/commit/4e54c3bae5ad8242fb1af0ab171aeb4c5b818923))
* show warning when user choose ARM based instance type in verbose mode during context deployment ([3967b81](https://github.com/aws/amazon-genomics-cli/commit/3967b81558d7b115159ff24125bf20af6d045a29))
* Shows the relevant error if the workflow logs can't be retrieved ([#103](https://github.com/aws/amazon-genomics-cli/issues/103)) ([2051b95](https://github.com/aws/amazon-genomics-cli/commit/2051b9542d07c5f999bd149e2a9f65aefaccba00))
* switch from a progress bar to a spinner ([#203](https://github.com/aws/amazon-genomics-cli/issues/203)) ([1ae9c8a](https://github.com/aws/amazon-genomics-cli/commit/1ae9c8aae867ac56ca488d9e4cb263ea597c0726))
* temporary folder potential leak in some error scenarios. unit test for cdk command execution ([#140](https://github.com/aws/amazon-genomics-cli/issues/140)) ([c50608e](https://github.com/aws/amazon-genomics-cli/commit/c50608e594b528a7bddd33b678da984feabc50b4))
* Typo in documentation ([#255](https://github.com/aws/amazon-genomics-cli/issues/255)) ([db57703](https://github.com/aws/amazon-genomics-cli/commit/db577031de8af7a64d61d1699b720bf2509dd260))
* unique project name ([#171](https://github.com/aws/amazon-genomics-cli/issues/171)) ([2cb4303](https://github.com/aws/amazon-genomics-cli/commit/2cb4303ab445b5893689f831fcd24309923db18c))
* update go dependency versions ([#284](https://github.com/aws/amazon-genomics-cli/issues/284)) ([55f9bed](https://github.com/aws/amazon-genomics-cli/commit/55f9bedf9170ddebff2bce014ff6bbb16adce004))
* update go dependency versions ([#284](https://github.com/aws/amazon-genomics-cli/issues/284)) ([5925da5](https://github.com/aws/amazon-genomics-cli/commit/5925da50c44075208d9d19252bbbb5e406acc12c))
* updates context describe to be consistent with context destroy ([#143](https://github.com/aws/amazon-genomics-cli/issues/143)) ([bb7667b](https://github.com/aws/amazon-genomics-cli/commit/bb7667b44027b3374b8011da11418d6ee0054b79))
* updates how the logs are shown from cloudwatch ([#142](https://github.com/aws/amazon-genomics-cli/issues/142)) ([1502578](https://github.com/aws/amazon-genomics-cli/commit/1502578415c7db4c5a633982301a887bcd393514))
* Use correct context name ([#132](https://github.com/aws/amazon-genomics-cli/issues/132)) ([c1516b6](https://github.com/aws/amazon-genomics-cli/commit/c1516b60b5706d06b30d7516a3aa1d80efd216af))
* use fixed python 3.9 in wes adapter ([f88a69e](https://github.com/aws/amazon-genomics-cli/commit/f88a69efb530d96ccc8ac35a60bad2d3fdd6d30c))
* use proper assertion in order to avoid sporadic unit test failures ([#167](https://github.com/aws/amazon-genomics-cli/issues/167)) ([a836f20](https://github.com/aws/amazon-genomics-cli/commit/a836f20f6c25ed049524ca7ab97a938bdc9d3a37))
* use proper go tags for windows build ([#117](https://github.com/aws/amazon-genomics-cli/issues/117)) ([7dfa94a](https://github.com/aws/amazon-genomics-cli/commit/7dfa94a775fdba5193c99d0c697c8013a52a23ce))
* Using shrinkwrap and ci to install set NPM dependencies ([#193](https://github.com/aws/amazon-genomics-cli/issues/193)) ([27f456e](https://github.com/aws/amazon-genomics-cli/commit/27f456edbcad31182b281fc0fbb61de72029d400))
* Workflow status now ignores unqueryable stacks ([#138](https://github.com/aws/amazon-genomics-cli/issues/138)) ([cf817a8](https://github.com/aws/amazon-genomics-cli/commit/cf817a882de2160d8e333c17d4eb28508cd886e1))
* workflows from demo-wdl-project should run without errors out of the box ([#108](https://github.com/aws/amazon-genomics-cli/issues/108)) ([6960eac](https://github.com/aws/amazon-genomics-cli/commit/6960eacf236e744d3c5658c5557061ab9cd3d468))

## [1.2.0](https://github.com/aws/amazon-genomics-cli/compare/v1.1.2...v1.2.0) (2022-02-17)


### Features

* add a global --silent flag ([#219](https://github.com/aws/amazon-genomics-cli/issues/219)) ([eed4189](https://github.com/aws/amazon-genomics-cli/commit/eed4189b2d988d88ef87be63b8f39610be569291))
* Add AGC version tag to AWS resources ([#243](https://github.com/aws/amazon-genomics-cli/issues/243)) ([fc6cad7](https://github.com/aws/amazon-genomics-cli/commit/fc6cad707eb4c90098d3c9435b23efccf3ac8a72))
* Add snakemake wes adapter + cdk + example workflow ([#304](https://github.com/aws/amazon-genomics-cli/issues/304)) ([6357b47](https://github.com/aws/amazon-genomics-cli/commit/6357b478522f718b6fcb44a56db142cb3edbb4db))
* Add WES API endpoint information to agc context describe output ([#253](https://github.com/aws/amazon-genomics-cli/issues/253)) ([6145db3](https://github.com/aws/amazon-genomics-cli/commit/6145db387344e8001e132320dbc404fecee38124))
* Bootstrap CDK during Account Activation ([#272](https://github.com/aws/amazon-genomics-cli/issues/272)) ([5228492](https://github.com/aws/amazon-genomics-cli/commit/5228492fba52aa9f6c075a9b27a293ca2298969a))
* Bootstrap CDK during Account Activation ([#272](https://github.com/aws/amazon-genomics-cli/issues/272)) ([219596e](https://github.com/aws/amazon-genomics-cli/commit/219596e8edf1666ca46a065bfa1ccc74848180f5))
* context destroy --force flag ([#118](https://github.com/aws/amazon-genomics-cli/issues/118)) ([4282093](https://github.com/aws/amazon-genomics-cli/commit/428209311aa247c999816348a972737739b1189f))
* don't use ebs auto scale for minidl ([#300](https://github.com/aws/amazon-genomics-cli/issues/300)) ([225d139](https://github.com/aws/amazon-genomics-cli/commit/225d139453fe233e93addfea257be13ab7f56952))
* Enable hosting assets in AWS commercial regions ([#240](https://github.com/aws/amazon-genomics-cli/issues/240)) ([2eb5007](https://github.com/aws/amazon-genomics-cli/commit/2eb5007ed304847f03c75ac16f7a58acc5e8ace1))
* Improve context stack deployment speed ([#150](https://github.com/aws/amazon-genomics-cli/issues/150)) ([f7d6c7f](https://github.com/aws/amazon-genomics-cli/commit/f7d6c7f4e5e0fb91b3851bf806bc411b7f4526e5))
* improve reliability of provision scripts ([#308](https://github.com/aws/amazon-genomics-cli/issues/308)) ([16eff36](https://github.com/aws/amazon-genomics-cli/commit/16eff364d56132712ec5fb67f5c16e9996acc7c2))
* Improved Workflow logs ([#156](https://github.com/aws/amazon-genomics-cli/issues/156)) ([f71231c](https://github.com/aws/amazon-genomics-cli/commit/f71231c556c7ecb92940f3685f55c7078b2d1028))
* Initial infrastructure for MiniWdl support ([#125](https://github.com/aws/amazon-genomics-cli/issues/125)) ([387393a](https://github.com/aws/amazon-genomics-cli/commit/387393a64593d08ec2016a69382f7d64de37914e))
* Introducing AWS Lambda based WES Adapter for running the workflows ([#155](https://github.com/aws/amazon-genomics-cli/issues/155)) ([bd90f47](https://github.com/aws/amazon-genomics-cli/commit/bd90f47b94ede34c31ea109221225ff3cd65d200))


### Bug Fixes

* Add AGC assets bucket to project read list ([#247](https://github.com/aws/amazon-genomics-cli/issues/247)) ([2019b6a](https://github.com/aws/amazon-genomics-cli/commit/2019b6ac0226276337ee4e93237c7cfb19c4d28e))
* add build wes adapter step to run dev script and readme ([#307](https://github.com/aws/amazon-genomics-cli/issues/307)) ([466ea69](https://github.com/aws/amazon-genomics-cli/commit/466ea69fa8076dd0db348490ccac81f6b5a7180f))
* Adds a message when new logs aren't shown to the user immediately ([#131](https://github.com/aws/amazon-genomics-cli/issues/131)) ([54349d2](https://github.com/aws/amazon-genomics-cli/commit/54349d2858a837da26e6479c409e4a8445055562))
* Allow home directory path in workflow project config ([#137](https://github.com/aws/amazon-genomics-cli/issues/137)) ([760ed3c](https://github.com/aws/amazon-genomics-cli/commit/760ed3cec26daacc1f7dd5f328a63b119aa98886))
* Asserts order deterministically ([#153](https://github.com/aws/amazon-genomics-cli/issues/153)) ([c4299e8](https://github.com/aws/amazon-genomics-cli/commit/c4299e86e499edeb4695152f558a88e82bcf2da3))
* container validation check with the correct region from user profile ([#310](https://github.com/aws/amazon-genomics-cli/issues/310)) ([c36a1f0](https://github.com/aws/amazon-genomics-cli/commit/c36a1f0ee052aea0cef721227c79a4b1370828fc))
* correctly handle missing return code from Cromwell executions ([#276](https://github.com/aws/amazon-genomics-cli/issues/276)) ([3c61735](https://github.com/aws/amazon-genomics-cli/commit/3c6173525b6e1b991c9ba52258822ad37fb024bc))
* correctly link to core app ([#133](https://github.com/aws/amazon-genomics-cli/issues/133)) ([ce103b2](https://github.com/aws/amazon-genomics-cli/commit/ce103b202d50c7a8a40e6d94daca0c4dd5141da7))
* Delete bootstrap bucket on account deactivate ([#306](https://github.com/aws/amazon-genomics-cli/issues/306)) ([a81710b](https://github.com/aws/amazon-genomics-cli/commit/a81710b1298648cae6a040065d1403d373cb0ba6))
* deploy wes adapter code as cdk asset ([#321](https://github.com/aws/amazon-genomics-cli/issues/321)) ([cbaa38c](https://github.com/aws/amazon-genomics-cli/commit/cbaa38c83115f8f4b10837a2bf175b23a324712a))
* Deregionalize min permissions ([#128](https://github.com/aws/amazon-genomics-cli/issues/128)) ([c4dc96b](https://github.com/aws/amazon-genomics-cli/commit/c4dc96b1641431ed7c20fad348e7d87d2156a4b8))
* Fix project data path prefix and docs ([#315](https://github.com/aws/amazon-genomics-cli/issues/315)) ([69a6bb4](https://github.com/aws/amazon-genomics-cli/commit/69a6bb4faca4b2ec3839905b02bbf79289dc81e4))
* fixed arg validation for format ([#198](https://github.com/aws/amazon-genomics-cli/issues/198)) ([0b68019](https://github.com/aws/amazon-genomics-cli/commit/0b6801942dd2fea22bcf8291833732f0c337d0e9))
* Fixes how inputs are processed in AGC via Manifest ([#277](https://github.com/aws/amazon-genomics-cli/issues/277)) ([be5c684](https://github.com/aws/amazon-genomics-cli/commit/be5c684ca0967baeda7cb8055d24d39b54cd3084))
* Fixes how users interact with the context commands ([#115](https://github.com/aws/amazon-genomics-cli/issues/115)) ([ffb3bb6](https://github.com/aws/amazon-genomics-cli/commit/ffb3bb6fdffeabd09a33288086c5442aa5e14c60))
* Fixes the mechanism for copying the directory contents ([#311](https://github.com/aws/amazon-genomics-cli/issues/311)) ([29f8e9e](https://github.com/aws/amazon-genomics-cli/commit/29f8e9ee012e09ca294547ea3dba225304b36bfa))
* fixing monocdk imports ([#190](https://github.com/aws/amazon-genomics-cli/issues/190)) ([9c94fd8](https://github.com/aws/amazon-genomics-cli/commit/9c94fd87bcd860fe3842536dbcedd2165fdcc80e))
* force use node v14 in cdk package ([d67a18c](https://github.com/aws/amazon-genomics-cli/commit/d67a18cdb31387642b015daeab44ac80a5972c14))
* function for default config ([#201](https://github.com/aws/amazon-genomics-cli/issues/201)) ([e6983ae](https://github.com/aws/amazon-genomics-cli/commit/e6983ae774f9d85cca634d4f9a0c04098ae6a234))
* Handle errors from CDK command ([#245](https://github.com/aws/amazon-genomics-cli/issues/245)) ([99d9378](https://github.com/aws/amazon-genomics-cli/commit/99d9378a3f5ad3f408ef24ac3dae42b46b79bbd5))
* improve contrast in docs ([#149](https://github.com/aws/amazon-genomics-cli/issues/149)) ([beb10f4](https://github.com/aws/amazon-genomics-cli/commit/beb10f4b02f9533da13ce0b3579ae2fd55a337aa))
* Increasing WES adapter lambda timeout ([#180](https://github.com/aws/amazon-genomics-cli/issues/180)) ([20cd77d](https://github.com/aws/amazon-genomics-cli/commit/20cd77dac1def4414f0159509caa8dfe853d62bb))
* invalid AWS Health url ([#130](https://github.com/aws/amazon-genomics-cli/issues/130)) ([1aef83b](https://github.com/aws/amazon-genomics-cli/commit/1aef83b682ba276ae5d8720ccaffc97a66bb34cb))
* miniwdl interpolation workaround ([27f4bf5](https://github.com/aws/amazon-genomics-cli/commit/27f4bf571712c6509e6352f4459e452fdd6a1cb1))
* nest batch artifacts and disable bucket clean up ([#262](https://github.com/aws/amazon-genomics-cli/issues/262)) ([73bb98d](https://github.com/aws/amazon-genomics-cli/commit/73bb98d23a18f210d4bfa231b0f8002bb294529c))
* Nextflow workflow logs unmarshal number ([#286](https://github.com/aws/amazon-genomics-cli/issues/286)) ([1ddf67c](https://github.com/aws/amazon-genomics-cli/commit/1ddf67c50e53e64c70acd003616a0f6e2f3af716))
* Nextflow workflow logs unmarshal number ([#286](https://github.com/aws/amazon-genomics-cli/issues/286)) ([f739f03](https://github.com/aws/amazon-genomics-cli/commit/f739f03b8825084a0f679ceb7b83a3db87b78c2e))
* Pass engine endpoint directly the wes adapter ([#122](https://github.com/aws/amazon-genomics-cli/issues/122)) ([81ed484](https://github.com/aws/amazon-genomics-cli/commit/81ed484a94ce195259315826377ece0443b582e1))
* Pinned running of rnaseq to version 3.4 ([5386a3c](https://github.com/aws/amazon-genomics-cli/commit/5386a3c4afd619989b61ad853d4c23aeecbe154e))
* progress bar doesn't show any progress ([#166](https://github.com/aws/amazon-genomics-cli/issues/166)) ([9dd17ae](https://github.com/aws/amazon-genomics-cli/commit/9dd17aef9e5edee719ed1c78f9d09aa3a8b4f9c8))
* Remove Chdir when looking for project files ([#228](https://github.com/aws/amazon-genomics-cli/issues/228)) ([3564d6f](https://github.com/aws/amazon-genomics-cli/commit/3564d6f5a2433191c3c31456f73d9b65f2cf1a87))
* rename unassigned variable ([#329](https://github.com/aws/amazon-genomics-cli/issues/329)) ([a3c26a8](https://github.com/aws/amazon-genomics-cli/commit/a3c26a805674e717901569123df3c44683b4cdab))
* Respect maxVCpus in miniwdl contexts ([#202](https://github.com/aws/amazon-genomics-cli/issues/202)) ([e4ad17e](https://github.com/aws/amazon-genomics-cli/commit/e4ad17efa9a3724c7dd9bc2f37430c21ebcf9d4e))
* Return timezone information in RunLog ([#174](https://github.com/aws/amazon-genomics-cli/issues/174)) ([118e1e2](https://github.com/aws/amazon-genomics-cli/commit/118e1e24ce8b6c454d641f27129bb04d6a44cf36))
* Revert add a global --silent flag ([#219](https://github.com/aws/amazon-genomics-cli/issues/219)) ([#279](https://github.com/aws/amazon-genomics-cli/issues/279)) ([a6cc990](https://github.com/aws/amazon-genomics-cli/commit/a6cc9904cd2d1899f7f64bf3c024a242c1a33749)), closes [#274](https://github.com/aws/amazon-genomics-cli/issues/274)
* Revert shrinkwrap upgrade ([#298](https://github.com/aws/amazon-genomics-cli/issues/298)) ([9dd6806](https://github.com/aws/amazon-genomics-cli/commit/9dd68060e044e4af487f0f67dc9d3da2a8439e8b))
* Scope down batch role permissions ([#230](https://github.com/aws/amazon-genomics-cli/issues/230)) ([14a5aca](https://github.com/aws/amazon-genomics-cli/commit/14a5acaf3f931cbaf53c0a495034cbff462797b0))
* show logs for workflows with more than 100 tasks ([#114](https://github.com/aws/amazon-genomics-cli/issues/114)) ([4e54c3b](https://github.com/aws/amazon-genomics-cli/commit/4e54c3bae5ad8242fb1af0ab171aeb4c5b818923))
* show warning when user choose ARM based instance type in verbose mode during context deployment ([3967b81](https://github.com/aws/amazon-genomics-cli/commit/3967b81558d7b115159ff24125bf20af6d045a29))
* Shows the relevant error if the workflow logs can't be retrieved ([#103](https://github.com/aws/amazon-genomics-cli/issues/103)) ([2051b95](https://github.com/aws/amazon-genomics-cli/commit/2051b9542d07c5f999bd149e2a9f65aefaccba00))
* switch from a progress bar to a spinner ([#203](https://github.com/aws/amazon-genomics-cli/issues/203)) ([1ae9c8a](https://github.com/aws/amazon-genomics-cli/commit/1ae9c8aae867ac56ca488d9e4cb263ea597c0726))
* temporary folder potential leak in some error scenarios. unit test for cdk command execution ([#140](https://github.com/aws/amazon-genomics-cli/issues/140)) ([c50608e](https://github.com/aws/amazon-genomics-cli/commit/c50608e594b528a7bddd33b678da984feabc50b4))
* Typo in documentation ([#255](https://github.com/aws/amazon-genomics-cli/issues/255)) ([db57703](https://github.com/aws/amazon-genomics-cli/commit/db577031de8af7a64d61d1699b720bf2509dd260))
* unique project name ([#171](https://github.com/aws/amazon-genomics-cli/issues/171)) ([2cb4303](https://github.com/aws/amazon-genomics-cli/commit/2cb4303ab445b5893689f831fcd24309923db18c))
* update go dependency versions ([#284](https://github.com/aws/amazon-genomics-cli/issues/284)) ([55f9bed](https://github.com/aws/amazon-genomics-cli/commit/55f9bedf9170ddebff2bce014ff6bbb16adce004))
* update go dependency versions ([#284](https://github.com/aws/amazon-genomics-cli/issues/284)) ([5925da5](https://github.com/aws/amazon-genomics-cli/commit/5925da50c44075208d9d19252bbbb5e406acc12c))
* updates context describe to be consistent with context destroy ([#143](https://github.com/aws/amazon-genomics-cli/issues/143)) ([bb7667b](https://github.com/aws/amazon-genomics-cli/commit/bb7667b44027b3374b8011da11418d6ee0054b79))
* updates how the logs are shown from cloudwatch ([#142](https://github.com/aws/amazon-genomics-cli/issues/142)) ([1502578](https://github.com/aws/amazon-genomics-cli/commit/1502578415c7db4c5a633982301a887bcd393514))
* Use correct context name ([#132](https://github.com/aws/amazon-genomics-cli/issues/132)) ([c1516b6](https://github.com/aws/amazon-genomics-cli/commit/c1516b60b5706d06b30d7516a3aa1d80efd216af))
* use fixed python 3.9 in wes adapter ([f88a69e](https://github.com/aws/amazon-genomics-cli/commit/f88a69efb530d96ccc8ac35a60bad2d3fdd6d30c))
* use proper assertion in order to avoid sporadic unit test failures ([#167](https://github.com/aws/amazon-genomics-cli/issues/167)) ([a836f20](https://github.com/aws/amazon-genomics-cli/commit/a836f20f6c25ed049524ca7ab97a938bdc9d3a37))
* use proper go tags for windows build ([#117](https://github.com/aws/amazon-genomics-cli/issues/117)) ([7dfa94a](https://github.com/aws/amazon-genomics-cli/commit/7dfa94a775fdba5193c99d0c697c8013a52a23ce))
* Using shrinkwrap and ci to install set NPM dependencies ([#193](https://github.com/aws/amazon-genomics-cli/issues/193)) ([27f456e](https://github.com/aws/amazon-genomics-cli/commit/27f456edbcad31182b281fc0fbb61de72029d400))
* Workflow status now ignores unqueryable stacks ([#138](https://github.com/aws/amazon-genomics-cli/issues/138)) ([cf817a8](https://github.com/aws/amazon-genomics-cli/commit/cf817a882de2160d8e333c17d4eb28508cd886e1))
* workflows from demo-wdl-project should run without errors out of the box ([#108](https://github.com/aws/amazon-genomics-cli/issues/108)) ([6960eac](https://github.com/aws/amazon-genomics-cli/commit/6960eacf236e744d3c5658c5557061ab9cd3d468))

## [1.1.2](https://github.com/aws/amazon-genomics-cli/compare/v1.1.1...v1.1.2) (2021-11-24)


### Bug Fixes

* fixing monocdk imports ([#190](https://github.com/aws/amazon-genomics-cli/issues/190)) ([6e4b295](https://github.com/aws/amazon-genomics-cli/commit/6e4b29551a0316cf8207871ec84681c34e91e96e))

## 1.1.1 (2021-11-15)


### Bug Fixes

* Correcting MiniWDL output location  ([#173](https://github.com/aws/amazon-genomics-cli/issues/173)) ([8e6f3fd](https://github.com/aws/amazon-genomics-cli/commit/8e6f3fda595531f0733b1d061cb84c05ed635923))
* fix for installation.md ([#161](https://github.com/aws/amazon-genomics-cli/issues/161)) ([72ee50c](https://github.com/aws/amazon-genomics-cli/commit/72ee50c9e9786f23e937a62e8e22a13f7aa909d5))
* progress bar doesn't show any progress ([#166](https://github.com/aws/amazon-genomics-cli/issues/166)) ([5b4a5b4](https://github.com/aws/amazon-genomics-cli/commit/5b4a5b4c81cf4ef99d82a146b378c66936ff7be4))
* Return timezone information in RunLog ([#174](https://github.com/aws/amazon-genomics-cli/issues/174)) ([37b33a7](https://github.com/aws/amazon-genomics-cli/commit/37b33a753d96010a581588c8734d482163ad1161))
* unique project name ([#171](https://github.com/aws/amazon-genomics-cli/issues/171)) ([a464f81](https://github.com/aws/amazon-genomics-cli/commit/a464f81d86575f1f3c0c95b1161f0474681537ac))
* Increasing WES adapter lambda timeout ([#180](https://github.com/aws/amazon-genomics-cli/issues/180)) ([20cd77d](https://github.com/aws/amazon-genomics-cli/commit/20cd77dac1def4414f0159509caa8dfe853d62bb))

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
