#!/usr/bin/env cwl-runner
cwlVersion: v1.2
class: CommandLineTool
baseCommand: ["sort"]
hints:
  ResourceRequirement:
    coresMax: 1
    outdirMin: $(inputs.input_file.size)
inputs:
  input_file:
    type: File
    inputBinding:
      position: 1
outputs:
  sorted_file:
    type: stdout
stdout: sorted.txt
