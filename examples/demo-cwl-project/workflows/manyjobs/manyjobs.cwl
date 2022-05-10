cwlVersion: v1.2
class: Workflow

requirements:
  SubworkflowFeatureRequirement: {}
  ScatterFeatureRequirement: {}

inputs:
  - id: script_file
    type: File
  - id: script_argument_count
    type: int
    
steps:
  make_array:
    run: make-array.cwl
    in:
      size: script_argument_count
    out: [array]

  scripts:
    run: run-script.cwl 
    scatter: arg
    in:
      src: script_file
      arg: make_array/array
    out: [script_output]
    
  sort:
    run: sort-file.cwl 
    scatter: input_file
    in:
      input_file: scripts/script_output
    out: [sorted_file]

  concat:
    run: concatenate-files.cwl
    in:
      files: sort/sorted_file
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
