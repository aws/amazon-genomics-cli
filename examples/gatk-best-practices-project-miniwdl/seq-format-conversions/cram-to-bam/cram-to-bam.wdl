version 1.0
## Copyright Broad Institute, 2017
## This script should convert a CRAM to SAM to BAM and output a BAM, BAM Index, and validation report to a Google bucket. If you'd like to do ## this on multiple CRAMS, create a sample set in the Data tab.
## The reason this approach was chosen instead of converting CRAM to BAM directly using Samtools is because Samtools 1.3 produces incorrect
## bins due to an old version of htslib included in the package. Samtools versions 1.4 & 1.5 have an NM issue that causes them to not validate ## with Picard.
##
## TESTED: It was tested using the Genomes in the Cloud Docker image version 2.3.1-1500064817.
## Versions of other tools on this image at the time of testing:
## PICARD_VER=1.1150
## GATK34_VER=3.4-g3c929b0
## GATK35_VER=3.5-0-g36282e4
## GATK36_VER=3.6-44-ge7d1cd2
## GATK4_VER=4.beta.1
## SAMTOOLS_VER=1.3.1
## BWA_VER=0.7.15.r1140
## TABIX_VER=0.2.5_r1005
## BGZIP_VER=1.3
## SVTOOLKIT_VER=2.00-1650
## It was tested pulling the HG38 reference Fasta and Fai.
## Successfully tested on Cromwell version 47. Does not work on versions < v23 due to output syntax
## Runtime parameters are optimized for Broad's Google Cloud Platform implementation.
##
## LICENSING : This script is released under the WDL source code license (BSD-3) (see LICENSE in https://github.com/broadinstitute/wdl).
## Note however that the programs it calls may be subject to different licenses. Users are responsible for checking that they are authorized to run all programs before running this script.
## Please see the docker for detailed licensing information pertaining to the included programs.
##
#WORKFLOW DEFINITION
workflow CramToBamFlow {
	input {
		File ref_fasta
		File ref_fasta_index
		File ref_dict
		File input_cram
		String sample_name
		String gotc_docker = "public.ecr.aws/aws-genomics/broadinstitute/genomes-in-the-cloud:2.4.7-1603303710"
		Int preemptible_tries = 3
	}

	#converts CRAM to SAM to BAM and makes BAI
	call CramToBamTask{
		input:
			ref_fasta = ref_fasta,
			ref_fasta_index = ref_fasta_index,
			ref_dict = ref_dict,
			input_cram = input_cram,
			sample_name = sample_name,
			docker_image = gotc_docker,
			preemptible_tries = preemptible_tries
	}

	#validates Bam
	call ValidateSamFile{
		input:
			input_bam = CramToBamTask.outputBam,
			docker_image = gotc_docker,
			preemptible_tries = preemptible_tries
	}

	#Outputs Bam, Bai, and validation report to the FireCloud data model
	output {
		File outputBam = CramToBamTask.outputBam
		File outputBai = CramToBamTask.outputBai
		File validation_report = ValidateSamFile.report
	}

}

#Task Definitions
task CramToBamTask {
	input {
		# Command parameters
		File ref_fasta
		File ref_fasta_index
		File ref_dict
		File input_cram
		String sample_name

		# Runtime parameters
		Int machine_mem_size = 15
		String docker_image
		Int preemptible_tries
	}

	#Calls samtools view to do the conversion
	command {
		set -eo pipefail

		samtools view -h -T ~{ref_fasta} ~{input_cram} |
		samtools view -b -o ~{sample_name}.bam -
		samtools index -b ~{sample_name}.bam
		mv ~{sample_name}.bam.bai ~{sample_name}.bai
	}

	#Run time attributes:
	#Use a docker with samtools. Set this up as a workspace attribute.
	#cpu of one because no multi-threading is required. This is also default, so don't need to specify.
	#disk_size should equal input size + output size + buffer
	runtime {
		docker: docker_image
		memory: machine_mem_size + " GB"
		preemptible: preemptible_tries
		awsBatchRetryAttempts: preemptible_tries
	}

	#Outputs a BAM and BAI with the same sample name
	output {
		File outputBam = "~{sample_name}.bam"
		File outputBai = "~{sample_name}.bai"
	}
}

#Validates BAM output to ensure it wasn't corrupted during the file conversion
task ValidateSamFile {
	input {
		File input_bam
		Int machine_mem_size = 4
		String docker_image
		Int preemptible_tries
	}
	String output_name = basename(input_bam, ".bam") + ".validation_report"
	Int command_mem_size = machine_mem_size - 1
	command {
		java -Xmx~{command_mem_size}G -jar /usr/gitc/picard.jar \
		ValidateSamFile \
		INPUT=~{input_bam} \
		OUTPUT=~{output_name} \
		MODE=SUMMARY \
		IS_BISULFITE_SEQUENCED=false
	}
	#Run time attributes:
	#Use a docker with the picard.jar. Set this up as a workspace attribute.
	#Read more about return codes here: https://github.com/broadinstitute/cromwell#continueonreturncode
	runtime {
		docker: docker_image
		memory: machine_mem_size + " GB"
		preemptible: preemptible_tries
		awsBatchRetryAttempts: preemptible_tries
		continueOnReturnCode: [0,1]
	}
	#A text file is generated that will list errors or warnings that apply.
	output {
		File report = "~{output_name}"
	}
}
