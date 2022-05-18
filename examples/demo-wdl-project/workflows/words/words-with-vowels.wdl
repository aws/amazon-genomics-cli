version 1.0
workflow Words_With_Vowels {
	input {
		File word_file
		Array[String] vowels = ["a", "e", "i", "o", "u"]
	}

	scatter (vowel in vowels) {
		call FindWordsWithLetter {input: vowel = vowel, word_file = word_file}
	}

	call SumWords {
		input:
			count_files = FindWordsWithLetter.count_file

	}

	output{
		Array[File] words_w_vowel = FindWordsWithLetter.filtered_file
		File sum_file = SumWords.sum_file
	}
}

task FindWordsWithLetter {
	input {
		String vowel
		File word_file
	}
	command {
		touch "words_with_~{vowel}.txt"
		grep ~{vowel} ~{word_file} > "words_with_~{vowel}.txt"
		wc -l "words_with_~{vowel}.txt" | cut -f1 -d " " > "~{vowel}-count.txt"
	}
	runtime {
		docker: "public.ecr.aws/amazonlinux/amazonlinux:2"
		#cromwell retry
		awsBatchRetryAttempts: 3
		#miniwdl retry
		maxRetries: 3
	}
	output {
		File filtered_file = "words_with_~{vowel}.txt"
		File count_file = "~{vowel}-count.txt"
	}
}

task SumWords {
	input {
		Array[File] count_files
	}
	command <<<
		awk '{ sum += $1 } END { print sum }' ~{sep=' ' count_files} > sum.txt
	>>>
	runtime {
		docker: "public.ecr.aws/amazonlinux/amazonlinux:2"
		#cromwell retry
		awsBatchRetryAttempts: 3
		#miniwdl retry
		maxRetries: 3
	}
   output {
	   File sum_file = "sum.txt"
       String out = read_string( "sum.txt" )
   }
}