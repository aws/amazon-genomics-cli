#!/usr/bin/env cwl-runner
cwlVersion: v1.2
class: ExpressionTool
requirements:
  InlineJavascriptRequirement: {}
inputs:
  size:
    type: int
outputs:
  array:
    type: string[]
expression: "$({array: function(){var arr = []; for (var i = 0; i < inputs.size; i++) {arr.push('' + i)}; return arr;}()})"

