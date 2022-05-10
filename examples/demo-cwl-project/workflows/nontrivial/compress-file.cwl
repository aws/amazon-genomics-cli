#!/usr/bin/env cwl-runner
cwlVersion: v1.2
class: CommandLineTool
baseCommand: ["-9", "-p", "8", "-c"]
hints:
  DockerRequirement:
    dockerPull: bytesco/pigz
  ResourceRequirement:
      coresMin: 8
inputs:
  input_file:
    type: File
    inputBinding:
      position: 1
outputs:
  compressed_file:
    type: stdout
stdout: compressed.gz
