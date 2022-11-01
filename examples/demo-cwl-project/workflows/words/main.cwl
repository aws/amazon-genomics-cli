cwlVersion: v1.0
class: Workflow
requirements:
  ScatterFeatureRequirement: {} 
inputs:
  words: File
  vowels: string[]
outputs:
  summaryFile:
    type: File
    outputSource: sumWords/summaryFile

steps:
  countWordsWithLetter:
    scatter: vowel
    in:
      words: words
      vowel: vowels
    out: [countFile]
    run:
      class: CommandLineTool
      baseCommand: grep
      inputs:
        words: File
        vowel: string
      arguments:
        - $(inputs.vowel)
        - $(inputs.words.path)
        - --count
      outputs:
        countFile:
          type: stdout
      stdout: count.txt
  sumWords:
    in:
      countFiles: [countWordsWithLetter/countFile]
    out: [summaryFile]
    run:
      class: CommandLineTool
      baseCommand: ["awk", "{ sum += $1 } END { print sum }"]
      inputs:
        countFiles:
          type: File[]
          inputBinding:
            position: 1
      outputs:
        summaryFile:
          type: stdout
      stdout: summary.txt