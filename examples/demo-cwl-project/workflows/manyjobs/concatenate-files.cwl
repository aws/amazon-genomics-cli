cwlVersion: v1.2
class: CommandLineTool
baseCommand: cat
inputs:
  files:
    type: File[]
    inputBinding:
      position: 1
  
outputs:
  concatenated_file:
    type: stdout

stdout: concatenated.txt
