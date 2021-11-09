version 1.0

## Copyright Broad Institute, 2019
##
## The haplotypecaller-gvcf-gatk4 workflow runs the HaplotypeCaller tool
## from GATK4 in GVCF mode on a single sample according to GATK Best Practices.
## When executed the workflow scatters the HaplotypeCaller tool over a sample
## using an intervals list file. The output file produced will be a
## single gvcf file which can be used by the joint-discovery workflow.
##
## Requirements/expectations :
## - One analysis-ready BAM file for a single sample (as identified in RG:SM)
## - Set of variant calling intervals lists for the scatter, provided in a file
##
## Outputs :
## - One GVCF file and its index
##
## Cromwell version support
## - Successfully tested on v53
##
## Runtime parameters are optimized for Broad's Google Cloud Platform implementation.
##
## LICENSING :
## This script is released under the WDL source code license (BSD-3) (see LICENSE in
## https://github.com/broadinstitute/wdl). Note however that the programs it calls may
## be subject to different licenses. Users are responsible for checking that they are
## authorized to run all programs before running this script. Please see the dockers
## for detailed licensing information pertaining to the included programs.

# WORKFLOW DEFINITION
workflow HaplotypeCallerGvcf_GATK4 {
    input {
        File input_bam
        File input_bam_index
        File ref_dict
        File ref_fasta
        File ref_fasta_index
        File scattered_calling_intervals_list

        Boolean make_gvcf = true
        Boolean make_bamout = false
        String gatk_docker = "us.gcr.io/broad-gatk/gatk:4.2.0.0"
        String gatk_path = "/gatk/gatk"
        String gitc_docker = "us.gcr.io/broad-gotc-prod/genomes-in-the-cloud:2.4.7-1603303710"
        String samtools_path = "samtools"
    }

    Array[File] scattered_calling_intervals = read_lines(scattered_calling_intervals_list)

    #is the input a cram file?
    Boolean is_cram = sub(basename(input_bam), ".*\\.", "") == "cram"

    String sample_basename = if is_cram then  basename(input_bam, ".cram") else basename(input_bam, ".bam")
    String vcf_basename = sample_basename
    String output_suffix = if make_gvcf then ".g.vcf.gz" else ".vcf.gz"
    String output_filename = vcf_basename + output_suffix

    # We need disk to localize the sharded input and output due to the scatter for HaplotypeCaller.
    # If we take the number we are scattering by and reduce by 20 we will have enough disk space
    # to account for the fact that the data is quite uneven across the shards.
    Int potential_hc_divisor = length(scattered_calling_intervals) - 20
    Int hc_divisor = if potential_hc_divisor > 1 then potential_hc_divisor else 1

    if ( is_cram ) {
        call CramToBamTask {
            input:
                input_cram = input_bam,
                sample_name = sample_basename,
                ref_dict = ref_dict,
                ref_fasta = ref_fasta,
                ref_fasta_index = ref_fasta_index,
                docker = gitc_docker,
                samtools_path = samtools_path
        }
    }

    # Call variants in parallel over grouped calling intervals
    scatter (interval_file in scattered_calling_intervals) {

        # Generate GVCF by interval
        call HaplotypeCaller {
            input:
                input_bam = select_first([CramToBamTask.output_bam, input_bam]),
                input_bam_index = select_first([CramToBamTask.output_bai, input_bam_index]),
                interval_list = interval_file,
                output_filename = output_filename,
                ref_dict = ref_dict,
                ref_fasta = ref_fasta,
                ref_fasta_index = ref_fasta_index,
                hc_scatter = hc_divisor,
                make_gvcf = make_gvcf,
                make_bamout = make_bamout,
                docker = gatk_docker,
                gatk_path = gatk_path
        }
    }

    # Merge per-interval GVCFs
    call MergeGVCFs {
        input:
            input_vcfs = HaplotypeCaller.output_vcf,
            input_vcfs_indexes = HaplotypeCaller.output_vcf_index,
            output_filename = output_filename,
            docker = gatk_docker,
            gatk_path = gatk_path
    }

    # Outputs that will be retained when execution is complete
    output {
        File output_vcf = MergeGVCFs.output_vcf
        File output_vcf_index = MergeGVCFs.output_vcf_index
    }
}

# TASK DEFINITIONS

task CramToBamTask {
    input {
        # Command parameters
        File ref_fasta
        File ref_fasta_index
        File ref_dict
        File input_cram
        String sample_name

        # Runtime parameters
        String docker
        Int? machine_mem_gb
        String samtools_path
        Int? retry_attempts
    }

    command {
        set -e
        set -o pipefail

        ~{samtools_path} view -h -T ~{ref_fasta} ~{input_cram} |
        ~{samtools_path} view -b -o ~{sample_name}.bam -
        ~{samtools_path} index -b ~{sample_name}.bam
        mv ~{sample_name}.bam.bai ~{sample_name}.bai
    }
    runtime {
        docker: docker
        memory: select_first([machine_mem_gb, 15]) + " GB"
        awsBatchRetryAttempts: select_first([retry_attempts, 3])
    }
    output {
        File output_bam = "~{sample_name}.bam"
        File output_bai = "~{sample_name}.bai"
    }
}

# HaplotypeCaller per-sample in GVCF mode
task HaplotypeCaller {
    input {
        # Command parameters
        File input_bam
        File input_bam_index
        File interval_list
        String output_filename
        File ref_dict
        File ref_fasta
        File ref_fasta_index
        Float? contamination
        Boolean make_gvcf
        Boolean make_bamout
        Int hc_scatter

        String? gcs_project_for_requester_pays

        String gatk_path
        String? java_options

        # Runtime parameters
        String docker
        Int? mem_gb
        Int? retry_attempts
    }

    String java_opt = select_first([java_options, "-XX:GCTimeLimit=50 -XX:GCHeapFreeLimit=10"])

    Int machine_mem_gb = select_first([mem_gb, 7])
    Int command_mem_gb = machine_mem_gb - 1

    String vcf_basename = if make_gvcf then  basename(output_filename, ".gvcf") else basename(output_filename, ".vcf")
    String bamout_arg = if make_bamout then "-bamout ~{vcf_basename}.bamout.bam" else ""

    parameter_meta {
        input_bam: {
                       description: "a bam file",
                       localization_optional: true
                   }
        input_bam_index: {
                             description: "an index file for the bam input",
                             localization_optional: true
                         }
    }
    command {
        set -e

        ~{gatk_path} --java-options "-Xmx~{command_mem_gb}G ~{java_opt}" \
        HaplotypeCaller \
        -R ~{ref_fasta} \
        -I ~{input_bam} \
        -L ~{interval_list} \
        -O ~{output_filename} \
        -contamination ~{default="0" contamination} \
        -G StandardAnnotation -G StandardHCAnnotation ~{true="-G AS_StandardAnnotation" false="" make_gvcf} \
        -GQB 10 -GQB 20 -GQB 30 -GQB 40 -GQB 50 -GQB 60 -GQB 70 -GQB 80 -GQB 90 \
        ~{true="-ERC GVCF" false="" make_gvcf} \
        ~{if defined(gcs_project_for_requester_pays) then "--gcs-project-for-requester-pays ~{gcs_project_for_requester_pays}" else ""} \
        ~{bamout_arg}

        # Cromwell doesn't like optional task outputs, so we have to touch this file.
        touch ~{vcf_basename}.bamout.bam
    }
    runtime {
        docker: docker
        memory: machine_mem_gb + " GB"
        awsBatchRetryAttempts: select_first([retry_attempts, 3])
    }
    output {
        File output_vcf = "~{output_filename}"
        File output_vcf_index = "~{output_filename}.tbi"
        File bamout = "~{vcf_basename}.bamout.bam"
    }
}
# Merge GVCFs generated per-interval for the same sample
task MergeGVCFs {
    input {
        # Command parameters
        Array[File] input_vcfs
        Array[File] input_vcfs_indexes
        String output_filename

        String gatk_path

        # Runtime parameters
        String docker
        Int? mem_gb
        Int? retry_attempts
    }
    Int machine_mem_gb = select_first([mem_gb, 3])
    Int command_mem_gb = machine_mem_gb - 1

    command {
        set -e

        ~{gatk_path} --java-options "-Xmx~{command_mem_gb}G"  \
        MergeVcfs \
        --INPUT ~{sep=' --INPUT ' input_vcfs} \
        --OUTPUT ~{output_filename}
    }
    runtime {
        docker: docker
        memory: machine_mem_gb + " GB"
        awsBatchRetryAttempts: select_first([retry_attempts, 3])
    }
    output {
        File output_vcf = "~{output_filename}"
        File output_vcf_index = "~{output_filename}.tbi"
    }
}

