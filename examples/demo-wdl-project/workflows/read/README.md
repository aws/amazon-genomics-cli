## How to provide input arguments

This demo workflow requires an "input_file" input argument to be defined. You can use the definition in read.inputs.json file by providing it as a workflow argument:

agc workflow run read --context myContext --args workflows/read/read.inputs.json

See [Workflows](https://aws.github.io/amazon-genomics-cli/docs/concepts/workflows/) documentation for further details.
