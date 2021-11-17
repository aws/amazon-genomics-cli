version 1.0
##Copyright Broad Institute, 2018
##
## This WDL converts paired FASTQ to uBAM and adds read group information
##
## Requirements/expectations :
## - Pair-end sequencing data in FASTQ format (one file per orientation)
## - The following metada descriptors per sample:
##  - readgroup
##  - sample_name
##  - library_name
##  - platform_unit
##  - run_date
##  - platform_name
##  - sequecing_center
##
## Outputs :
## - Set of unmapped BAMs, one per read group
## - File of a list of the generated unmapped BAMs
##
## Cromwell version support
## - Successfully tested on v47
## - Does not work on versions < v23 due to output syntax
##
## Runtime parameters are optimized for Broad's Google Cloud Platform implementation.
## For program versions, see docker containers.
##
## LICENSING :
## This script is released under the WDL source code license (BSD-3) (see LICENSE in
## https://github.com/broadinstitute/wdl). Note however that the programs it calls may
## be subject to different licenses. Users are responsible for checking that they are
## authorized to run all programs before running this script. Please see the docker
## page at https://hub.docker.com/r/broadinstitute/genomes-in-the-cloud/ for detailed
## licensing information pertaining to the included programs.

# WORKFLOW DEFINITION
workflow ConvertPairedFastQsToUnmappedBamWf {
	input {
		String sample_name
		File fastq_1
		File fastq_2
		String readgroup_name
		String library_name
		String platform_unit
		String run_date
		String platform_name
		String sequencing_center

		Boolean make_fofn = false

		String gatk_docker = "public.ecr.aws/aws-genomics/broadinstitute/gatk:4.1.9.0"
		String gatk_path = "/gatk/gatk"

	}

	String ubam_list_name = sample_name

	# Convert pair of FASTQs to uBAM
	call PairedFastQsToUnmappedBAM {
		input:
			sample_name = sample_name,
			fastq_1 = fastq_1,
			fastq_2 = fastq_2,
			readgroup_name = readgroup_name,
			library_name = library_name,
			platform_unit = platform_unit,
			run_date = run_date,
			platform_name = platform_name,
			sequencing_center = sequencing_center,
			gatk_path = gatk_path,
			docker = gatk_docker,
	}

	#Create a file with the generated ubam
	if (make_fofn) {
		call CreateFoFN {
			input:
				ubam = PairedFastQsToUnmappedBAM.output_unmapped_bam,
				fofn_name = ubam_list_name + ".ubam"
		}
	}

	# Outputs that will be retained when execution is complete
	output {
		File output_unmapped_bam = PairedFastQsToUnmappedBAM.output_unmapped_bam
		File? unmapped_bam_list = CreateFoFN.fofn_list
	}
}

# TASK DEFINITIONS

# Convert a pair of FASTQs to uBAM
task PairedFastQsToUnmappedBAM {
	input {
		# Command parameters
		String sample_name
		File fastq_1
		File fastq_2
		String readgroup_name
		String library_name
		String platform_unit
		String run_date
		String platform_name
		String sequencing_center
		String gatk_path

		# Runtime parameters
		Int machine_mem_gb = 7
		Int preemptible_attempts = 3
		String docker
	}
	Int command_mem_gb = machine_mem_gb - 1
	command {
		~{gatk_path} --java-options "-Xmx~{command_mem_gb}g" \
		FastqToSam \
		--FASTQ ~{fastq_1} \
		--FASTQ2 ~{fastq_2} \
		--OUTPUT ~{readgroup_name}.unmapped.bam \
		--READ_GROUP_NAME ~{readgroup_name} \
		--SAMPLE_NAME ~{sample_name} \
		--LIBRARY_NAME ~{library_name} \
		--PLATFORM_UNIT ~{platform_unit} \
		--RUN_DATE ~{run_date} \
		--PLATFORM ~{platform_name} \
		--SEQUENCING_CENTER ~{sequencing_center}
	}
	runtime {
		docker: docker
		memory: machine_mem_gb + " GB"
		preemptible: preemptible_attempts
		awsBatchRetryAttempts: preemptible_attempts
	}
	output {
		File output_unmapped_bam = "~{readgroup_name}.unmapped.bam"
	}
}

# Creats a file of file names of the uBAM, which is a text file with each row having the path to the file.
# In this case there will only be one file path in the txt file but this format is used by
# the pre-processing for variant discvoery workflow.
task CreateFoFN {
	input {
		# Command parameters
		String ubam
		String fofn_name
	}
	command {
		echo ~{ubam} > ~{fofn_name}.list
	}
	output {
		File fofn_list = "~{fofn_name}.list"
	}
	runtime {
		docker: "public.ecr.aws/lts/ubuntu:20.04_stable"
		preemptible: 3
		awsBatchRetryAttempts: 3
	}
}
