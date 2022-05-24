#!/usr/bin/env cwl-runner
cwlVersion: v1.2
class: CommandLineTool
hints:
  DockerRequirement:
    # Make sure to use a Docker image without an ENTRYPOINT, or it will behave
    # differently depending on if you pass --singularity to your CWL runner or
    # not.
    dockerPull: oittaa/imagemagick
  ResourceRequirement:
    coresMin: 2
    coresMax: 1
    outdirMin: 1024
    tmpdirMin: 1024
    ramMin: 4096
    ramMax: 4096
inputs:
  input_image:
    type: File
    inputBinding:
      position: 1
outputs:
  output_image:
    type: File
    outputBinding:
      glob: output.jpg
baseCommand: convert
arguments:
  - position: 2
    valueFrom: "-resize"
  - position: 3
    valueFrom: "1024x1024"
  - position: 4
    valueFrom: "output.jpg"
