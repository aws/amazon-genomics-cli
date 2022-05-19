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
      baseCommand: ["bash", "script.sh"]
      inputs:
        words: File
        vowel: string
      requirements:
        InitialWorkDirRequirement:
          listing:
            - entryname: script.sh
              entry: |-
                set -e
                VOWEL=$(inputs.vowel)
                WORD_FILE=$(inputs.words.path)

                grep \${VOWEL} \${WORD_FILE} | wc -l
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