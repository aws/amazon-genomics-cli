cwlVersion: v1.2
class: Workflow

requirements:
  SubworkflowFeatureRequirement: {}

inputs:
  - id: image_file
    type: File
  - id: image_directory
    type: Directory
  - id: image_filename
    type: string

steps:
  resize_file:
    run: shrink-image.cwl
    in:
      input_image: image_file
    out: [output_image]

  file_from_directory:
    run: file-from-directory.cwl
    in:
      dir: image_directory
      filename: image_filename
    out: [file]

  resize_from_directory:
    run: shrink-image.cwl
    in:
      input_image: file_from_directory/file
    out: [output_image]

outputs:
  - id: image_from_file
    type: File
    outputSource: resize_file/output_image
  - id: image_from_directory
    type: File
    outputSource: resize_from_directory/output_image
