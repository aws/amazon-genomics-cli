#!/usr/bin/env cwl-runner
# Modified from the CWL docs
cwlVersion: v1.2
class: CommandLineTool
baseCommand: node
hints:
  DockerRequirement:
    dockerPull: node:slim
inputs:
  src:
    type: File
    inputBinding:
      position: 1
  arg:
    type: string
    inputBinding:
      position: 2
outputs:
  script_output:
    type: stdout
stdout: output.txt
