#!/usr/bin/env cwl-runner
cwlVersion: v1.2
class: CommandLineTool
baseCommand: ["sort"]
requirements:
  InlineJavascriptRequirement: {}
hints:
  ResourceRequirement:
    coresMax: 1
    outdirMin: $(parseInt(Math.ceil(inputs.input_file.size / Math.pow(2, 20))))
inputs:
  input_file:
    type: File
    inputBinding:
      position: 1
outputs:
  sorted_file:
    type: stdout
stdout: sorted.txt
