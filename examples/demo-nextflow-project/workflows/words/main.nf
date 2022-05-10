nextflow.enable.dsl = 2

params.word_file = null
params.vowels = ["a", "e", "i", "o", "u"]

if (!params.word_file) exit 1, "required parameter 'word_file' missing"

word_file = file(params.word_file)

process FindWordsWithLetter {
    container "public.ecr.aws/amazonlinux/amazonlinux:2"

    input:
        val vowel
        path word_file
    
    output:
        path "words_with_${vowel}.txt", emit: filtered_file
        path "${vowel}-count.txt", emit: count_file
    
    script:
        """
        touch "words_with_${vowel}.txt"
        grep ${vowel} ${word_file} > "words_with_${vowel}.txt"
        wc -l "words_with_${vowel}.txt" | cut -f1 -d " " > "${vowel}-count.txt"
        """
    
    stub:
        """
        touch "words_with_${vowel}.txt"
        touch "${vowel}-count.txt"
        """
}

process SumWords {
    container "public.ecr.aws/amazonlinux/amazonlinux:2"

    input:
        path count_files
    
    output:
        path "sum.txt", emit: sum_file
        stdout emit: out
    
    shell:
        count_files_params = count_files
            .sort(false)
            .join(" ")
        
        '''
        awk '{ sum += $1 } END { print sum }' !{count_files_params} > sum.txt
        cat sum.txt
        '''
    
    stub:
        """
        touch "sum.txt"
        echo "999"
        """
}

workflow {

    // nextflow implicitly scatters processes over channels
    // create a channle of intervals to run FindWordsWithLetter
    vowels = Channel.fromList(params.vowels)
    FindWordsWithLetter(
        vowels,
        word_file
    )

    SumWords(
        FindWordsWithLetter.out.filtered_file.collect()
    )

    FindWordsWithLetter.out.filtered_file.view()
    SumWords.out.sum_file.view()
}