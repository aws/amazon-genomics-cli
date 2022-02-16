version 1.0
workflow Words_With_Vowels {
	input {
		File word_file
		Array[String] vowels = ["a", "e", "i", "o", "u"]
	}

	scatter (vowel in vowels) {
		call FindWordsWithLetter {input: vowel = vowel, word_file = word_file}
	}

	output{
		Array[File] words_w_vowel = FindWordsWithLetter.filtered_file
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
	}
	runtime {
		docker: "public.ecr.aws/amazonlinux/amazonlinux:2"
		#cromwell retry
		awsBatchRetryAttempts: 3
		#miniwdl retry
		maxRetries: 3
	}
	output { File filtered_file = "words_with_~{vowel}.txt" }
}