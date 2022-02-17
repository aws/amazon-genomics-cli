cwlVersion: v1.2
class: Workflow

requirements:
  SubworkflowFeatureRequirement: {}
  ScatterFeatureRequirement: {}

inputs:
  - id: script_file
    type: File
  - id: script_arguments
    type: string[]
    
steps:
  scripts:
    run: run-script.cwl 
    scatter: arg
    in:
      src: script_file
      arg: script_arguments
    out: [script_output]

  concat:
    run: concatenate-files.cwl
    in:
      files: scripts/script_output
    out:
      [concatenated_file]
      
  compress:
    run: compress-file.cwl 
    in:
      input_file: concat/concatenated_file
    out: [compressed_file]
    
outputs:
  - id: output
    type: File
    outputSource: compress/compressed_file
