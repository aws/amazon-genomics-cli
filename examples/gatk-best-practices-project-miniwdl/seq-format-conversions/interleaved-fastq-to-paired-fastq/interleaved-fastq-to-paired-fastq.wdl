version 1.0
# This WDL takes in a single interleaved(R1+R2) FASTQ file and separates it into
# separate R1 and R2 FASTQ (i.e. paired FASTQ) files. Paired FASTQ files are the
# input format for the tool that generates unmapped BAMs (the format used in most
# GATK processing and analysis tools).
#
# Requirements/expectations
# - Interleaved Fastq file
#
# Outputs
# - Separate R1 and R2 FASTQ files (i.e. paired FASTQ)
#
# LICENSING : This script is released under the WDL source code license (BSD-3) (see LICENSE in https://github.com/broadinstitute/wdl).
# Note however that the programs it calls may be subject to different licenses. Users are responsible for checking that they are authorized to run all programs before running this script.
# Please see the docker for detailed licensing information pertaining to the included programs.
##################

workflow UninterleaveFastqs {
	input {
		File input_fastq
	}

	call uninterleave_fqs {
		input:
			input_fastq = input_fastq
	}

}

task uninterleave_fqs {
	input {
		File input_fastq

		Int machine_mem_gb = 8
	}
	String r1_name = basename(input_fastq, ".fastq") + "_reads_1.fastq"
	String r2_name = basename(input_fastq, ".fastq") + "_reads_2.fastq"

	command {
		cat ~{input_fastq} | paste - - - - - - - -  | \
		tee >(cut -f 1-4 | tr "\t" "\n" > ~{r1_name}) | \
		cut -f 5-8 | tr "\t" "\n" > ~{r2_name}
	}

	runtime {
		docker: "public.ecr.aws/lts/ubuntu:20.04_stable"
		memory: machine_mem_gb + " GB"
	}

	output {
		File r1_fastq = "~{r1_name}"
		File r2_fastq = "~{r2_name}"
	}
}