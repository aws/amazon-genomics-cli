## Copyright Broad Institute, 2019
##
## Workflows for processing RNA data for germline short variant discovery with GATK (v4) and related tools
##
## Requirements/expectations :
## - BAM
##
## Output :
## - A BAM file and its index.
## - A VCF file and its index.
## - A Filtered VCF file and its index.
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
version 1.0


workflow RNAseq {
	input {
		Int? haplotypeScatterCount
		File? zippedStarReferences
		File annotationsGTF
		File dbSnpVcfIndex
		String? star_docker_override
		File inputBam
		String? gatk_path_override
		Array[File] knownVcfs
		File dbSnpVcf
		String? gatk4_docker_override
		File refFasta
		Int? preemptible_tries
		File refDict
		File refFastaIndex
		Int? minConfidenceForVariantCalling
		Int? readLength
		Array[File] knownVcfsIndices
	}
	String star_docker = select_first([star_docker_override, "quay.io/humancellatlas/secondary-analysis-star:v0.2.2-2.5.3a-40ead6e"])
	String gatk4_docker = select_first([gatk4_docker_override, "broadinstitute/gatk:latest"])
	File starReferences = select_first([zippedStarReferences, StarGenerateReferences.star_genome_refs_zipped, ""])
	String sampleName = basename(inputBam, ".bam")
	String gatk_path = select_first([gatk_path_override, "/gatk/gatk"])
	Int preemptible_count = select_first([preemptible_tries, 3])
	Int scatterCount = select_first([haplotypeScatterCount, 6])
	if (!defined(zippedStarReferences)) {
		call StarGenerateReferences {
			input:
				ref_fasta = refFasta,
				ref_fasta_index = refFastaIndex,
				annotations_gtf = annotationsGTF,
				read_length = readLength,
				docker = star_docker,
				preemptible_count = preemptible_count
		}
	}
	scatter (interval in ScatterIntervalList.out) {
		File HaplotypeCallerOutputVcfIndex = HaplotypeCaller.output_vcf_index
		File HaplotypeCallerOutputVcf = HaplotypeCaller.output_vcf
		call HaplotypeCaller {
			input:
				input_bam = ApplyBQSR.output_bam,
				input_bam_index = ApplyBQSR.output_bam_index,
				base_name = sampleName + ".hc",
				interval_list = interval,
				ref_dict = refDict,
				ref_fasta = refFasta,
				ref_fasta_index = refFastaIndex,
				dbSNP_vcf = dbSnpVcf,
				dbSNP_vcf_index = dbSnpVcfIndex,
				gatk_path = gatk_path,
				docker = gatk4_docker,
				preemptible_count = preemptible_count,
				stand_call_conf = minConfidenceForVariantCalling
		}
	}
	call ScatterIntervalList {
		input:
			interval_list = gtfToCallingIntervals.interval_list,
			scatter_count = scatterCount,
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call ApplyBQSR {
		input:
			input_bam = SplitNCigarReads.output_bam,
			input_bam_index = SplitNCigarReads.output_bam_index,
			base_name = sampleName + ".aligned.duplicates_marked.recalibrated",
			recalibration_report = BaseRecalibrator.recalibration_report,
			ref_dict = refDict,
			ref_fasta = refFasta,
			ref_fasta_index = refFastaIndex,
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call gtfToCallingIntervals {
		input:
			gtf = annotationsGTF,
			ref_dict = refDict,
			docker = gatk4_docker,
			gatk_path = gatk_path,
			preemptible_count = preemptible_count
	}
	call MergeBamAlignment {
		input:
			ref_fasta = refFasta,
			ref_dict = refDict,
			unaligned_bam = RevertSam.output_bam,
			star_bam = StarAlign.output_bam,
			base_name = ".merged",
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call StarAlign {
		input:
			star_genome_refs_zipped = starReferences,
			fastq1 = SamToFastq.fastq1,
			fastq2 = SamToFastq.fastq2,
			base_name = sampleName + ".star",
			read_length = readLength,
			docker = star_docker,
			preemptible_count = preemptible_count
	}
	call SamToFastq {
		input:
			unmapped_bam = RevertSam.output_bam,
			base_name = sampleName,
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call MergeVCFs {
		input:
			input_vcfs = HaplotypeCallerOutputVcf,
			input_vcfs_indexes = HaplotypeCallerOutputVcfIndex,
			output_vcf_name = sampleName + ".g.vcf.gz",
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call SplitNCigarReads {
		input:
			input_bam = MarkDuplicates.output_bam,
			input_bam_index = MarkDuplicates.output_bam_index,
			base_name = sampleName + ".split",
			ref_fasta = refFasta,
			ref_fasta_index = refFastaIndex,
			ref_dict = refDict,
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call MarkDuplicates {
		input:
			input_bam = MergeBamAlignment.output_bam,
			base_name = sampleName + ".dedupped",
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call VariantFiltration {
		input:
			input_vcf = MergeVCFs.output_vcf,
			input_vcf_index = MergeVCFs.output_vcf_index,
			base_name = sampleName + ".variant_filtered.vcf.gz",
			ref_dict = refDict,
			ref_fasta = refFasta,
			ref_fasta_index = refFastaIndex,
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call BaseRecalibrator {
		input:
			input_bam = SplitNCigarReads.output_bam,
			input_bam_index = SplitNCigarReads.output_bam_index,
			recal_output_file = sampleName + ".recal_data.csv",
			dbSNP_vcf = dbSnpVcf,
			dbSNP_vcf_index = dbSnpVcfIndex,
			known_indels_sites_VCFs = knownVcfs,
			known_indels_sites_indices = knownVcfsIndices,
			ref_dict = refDict,
			ref_fasta = refFasta,
			ref_fasta_index = refFastaIndex,
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}
	call RevertSam {
		input:
			input_bam = inputBam,
			base_name = sampleName + ".reverted",
			sort_order = "queryname",
			gatk_path = gatk_path,
			docker = gatk4_docker,
			preemptible_count = preemptible_count
	}

	output {
		File variant_filtered_vcf_index = VariantFiltration.output_vcf_index
		File merged_vcf = MergeVCFs.output_vcf
		File merged_vcf_index = MergeVCFs.output_vcf_index
		File variant_filtered_vcf = VariantFiltration.output_vcf
		File recalibrated_bam_index = ApplyBQSR.output_bam_index
		File recalibrated_bam = ApplyBQSR.output_bam
	}
}
task MarkDuplicates {
	input {
		File input_bam
		String base_name
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File output_bam = "${base_name}.bam"
		File output_bam_index = "${base_name}.bai"
		File metrics_file = "${base_name}.metrics"
	}
	command <<<

		~{gatk_path} \
		MarkDuplicates \
		--INPUT ~{input_bam} \
		--OUTPUT ~{base_name}.bam  \
		--CREATE_INDEX true \
		--VALIDATION_STRINGENCY SILENT \
		--METRICS_FILE ~{base_name}.metrics

	>>>
	runtime {
		# disks: "local-disk " + sub(((size(input_bam, "GB") + 1) * 3), "\..*", "") + " HDD"
		docker: docker
		memory: "4 GB"
		preemptible: preemptible_count
	}

}
task BaseRecalibrator {
	input {
		File input_bam
		File input_bam_index
		String recal_output_file
		File dbSNP_vcf
		File dbSNP_vcf_index
		Array[File] known_indels_sites_VCFs
		Array[File] known_indels_sites_indices
		File ref_dict
		File ref_fasta
		File ref_fasta_index
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File recalibration_report = recal_output_file
	}
	command <<<

		~{gatk_path} --java-options "-XX:GCTimeLimit=50 -XX:GCHeapFreeLimit=10 -XX:+PrintFlagsFinal \
		-XX:+PrintGCTimeStamps -XX:+PrintGCDateStamps -XX:+PrintGCDetails \
		-Xloggc:gc_log.log -Xms4000m" \
		BaseRecalibrator \
		-R ~{ref_fasta} \
		-I ~{input_bam} \
		--use-original-qualities \
		-O ~{recal_output_file} \
		-known-sites ~{dbSNP_vcf} \
		-known-sites ~{sep=" --known-sites "  known_indels_sites_VCFs}

	>>>
	runtime {
		memory: "6 GB"
		# disks: "local-disk " + sub((size(input_bam, "GB") * 3) + 30, "\..*", "") + " HDD"
		docker: docker
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task RevertSam {
	input {
		File input_bam
		String base_name
		String sort_order
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File output_bam = "${base_name}.bam"
	}
	command <<<

		~{gatk_path} \
		RevertSam \
		--INPUT ~{input_bam} \
		--OUTPUT ~{base_name}.bam \
		--VALIDATION_STRINGENCY SILENT \
		--ATTRIBUTE_TO_CLEAR FT \
		--ATTRIBUTE_TO_CLEAR CO \
		--SORT_ORDER ~{sort_order}

	>>>
	runtime {
		docker: docker
		# disks: "local-disk " + sub(((size(input_bam, "GB") + 1) * 5), "\..*", "") + " HDD"
		memory: "4 GB"
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task ApplyBQSR {
	input {
		File input_bam
		File input_bam_index
		String base_name
		File recalibration_report
		File ref_dict
		File ref_fasta
		File ref_fasta_index
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File output_bam = "${base_name}.bam"
		File output_bam_index = "${base_name}.bai"
	}
	command <<<

		~{gatk_path} \
		--java-options "-XX:+PrintFlagsFinal -XX:+PrintGCTimeStamps -XX:+PrintGCDateStamps \
		-XX:+PrintGCDetails -Xloggc:gc_log.log \
		-XX:GCTimeLimit=50 -XX:GCHeapFreeLimit=10 -Xms3000m" \
		ApplyBQSR \
		--add-output-sam-program-record \
		-R ~{ref_fasta} \
		-I ~{input_bam} \
		--use-original-qualities \
		-O ~{base_name}.bam \
		--bqsr-recal-file ~{recalibration_report}

	>>>
	runtime {
		memory: "3500 MB"
		# disks: "local-disk " + sub((size(input_bam, "GB") * 4) + 30, "\..*", "") + " HDD"
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
		docker: docker
	}

}
task StarAlign {
	input {
		File star_genome_refs_zipped
		File fastq1
		File fastq2
		String base_name
		Int? read_length
		Int length = select_first([read_length, 101])
		Int? num_threads
		Int threads = select_first([num_threads, 8])
		Int? star_mem_max_gb
		Int star_mem = select_first([star_mem_max_gb, 45])
		Int? star_limitOutSJcollapsed
		#Int? additional_disk
		#Int add_to_disk = select_first([additional_disk, 0])
		String docker
		Int preemptible_count
	}


	output {
		File output_bam = "${base_name}.Aligned.sortedByCoord.out.bam"
		File output_log_final = "${base_name}.Log.final.out"
		File output_log = "${base_name}.Log.out"
		File output_log_progress = "${base_name}.Log.progress.out"
		File output_SJ = "${base_name}.SJ.out.tab"
	}
	command <<<

		set -e

		tar -xvzf ~{star_genome_refs_zipped}

		STAR \
		--genomeDir STAR2_5 \
		--runThreadN ~{threads} \
		--readFilesIn ~{fastq1} ~{fastq2} \
		--readFilesCommand "gunzip -c" \
		~{"--sjdbOverhang " + (length - 1)} \
		--outSAMtype BAM SortedByCoordinate \
		--twopassMode Basic \
		--limitBAMsortRAM ~{star_mem + "000000000"} \
		--limitOutSJcollapsed ~{default="1000000"  star_limitOutSJcollapsed} \
		--outFileNamePrefix ~{base_name}.

	>>>
	runtime {
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
		# disks: "local-disk " + sub(((size(fastq1, "GB") + size(fastq2, "GB") * 10) + 30 + add_to_disk), "\..*", "") + " HDD"
		docker: docker
		cpu: threads
		memory: (star_mem + 1) + " GB"
	}

}
task MergeVCFs {
	input {
		Array[File] input_vcfs
		Array[File] input_vcfs_indexes
		String output_vcf_name
		#Int? disk_size = 5
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File output_vcf = output_vcf_name
		File output_vcf_index = "${output_vcf_name}.tbi"
	}
	command <<<

		~{gatk_path} --java-options "-Xms2000m"  \
		MergeVcfs \
		--INPUT ~{sep=" --INPUT "  input_vcfs} \
		--OUTPUT ~{output_vcf_name}

	>>>
	runtime {
		memory: "3 GB"
		# disks: "local-disk " + disk_size + " HDD"
		docker: docker
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task VariantFiltration {
	input {
		File input_vcf
		File input_vcf_index
		String base_name
		File ref_dict
		File ref_fasta
		File ref_fasta_index
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File output_vcf = "${base_name}"
		File output_vcf_index = "${base_name}.tbi"
	}
	command <<<

		~{gatk_path} \
		VariantFiltration \
		--R ~{ref_fasta} \
		--V ~{input_vcf} \
		--window 35 \
		--cluster 3 \
		--filter-name "FS" \
		--filter "FS > 30.0" \
		--filter-name "QD" \
		--filter "QD < 2.0" \
		-O ~{base_name}

	>>>
	runtime {
		docker: docker
		memory: "3 GB"
		# disks: "local-disk " + sub((size(input_vcf, "GB") * 2) + 30, "\..*", "") + " HDD"
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task gtfToCallingIntervals {
	input {
		File gtf
		File ref_dict
		String output_name = basename(gtf, ".gtf") + ".exons.interval_list"
		String docker
		String gatk_path
		Int preemptible_count
	}


	output {
		File interval_list = "${output_name}"
	}
	command <<<


		set -e

		Rscript --no-save -<<'RCODE'
		gtf = read.table("~{gtf}", sep="\t")
		gtf = subset(gtf, V3 == "exon")
		write.table(data.frame(chrom=gtf[,'V1'], start=gtf[,'V4'], end=gtf[,'V5']), "exome.bed", quote = F, sep="\t", col.names = F, row.names = F)
		RCODE

		awk '{print $1 "\t" ($2 - 1) "\t" $3}' exome.bed > exome.fixed.bed

		~{gatk_path} \
		BedToIntervalList \
		-I exome.fixed.bed \
		-O ~{output_name} \
		-SD ~{ref_dict}

	>>>
	runtime {
		docker: docker
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task ScatterIntervalList {
	input {
		File interval_list
		Int scatter_count
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		Array[File] out = glob("out/*/*.interval_list")
		Int interval_count = read_int("interval_count.txt")
	}
	command <<<

		set -e
		mkdir out
		~{gatk_path} --java-options "-Xms1g" \
		IntervalListTools \
		--SCATTER_COUNT ~{scatter_count} \
		--SUBDIVISION_MODE BALANCING_WITHOUT_INTERVAL_SUBDIVISION_WITH_OVERFLOW \
		--UNIQUE true \
		--SORT true \
		--INPUT ~{interval_list} \
		--OUTPUT out

		python3 <<CODE
import glob, os
# Works around a JES limitation where multiples files with the same name overwrite each other when globbed
intervals = sorted(glob.glob("out/*/*.interval_list"))
for i, interval in enumerate(intervals):
	(directory, filename) = os.path.split(interval)
newName = os.path.join(directory, str(i + 1) + filename)
os.rename(interval, newName)
print(len(intervals))
if len(intervals) == 0:
	raise ValueError("Interval list produced 0 scattered interval lists. Is the gtf or input interval list empty?")
f = open("interval_count.txt", "w+")
f.write(str(len(intervals)))
f.close()

CODE

	>>>
	runtime {
		# disks: "local-disk 1 HDD"
		memory: "2 GB"
		docker: docker
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task MergeBamAlignment {
	input {
		File ref_fasta
		File ref_dict
		File unaligned_bam
		File star_bam
		String base_name
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File output_bam = "${base_name}.bam"
	}
	command <<<

		~{gatk_path} \
		MergeBamAlignment \
		--REFERENCE_SEQUENCE ~{ref_fasta} \
		--UNMAPPED_BAM ~{unaligned_bam} \
		--ALIGNED_BAM ~{star_bam} \
		--OUTPUT ~{base_name}.bam \
		--INCLUDE_SECONDARY_ALIGNMENTS false \
		--VALIDATION_STRINGENCY SILENT

	>>>
	runtime {
		docker: docker
		# disks: "local-disk " + sub(((size(unaligned_bam, "GB") + size(star_bam, "GB") + 1) * 5), "\..*", "") + " HDD"
		memory: "4 GB"
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task SplitNCigarReads {
	input {
		File input_bam
		File input_bam_index
		String base_name
		File ref_fasta
		File ref_fasta_index
		File ref_dict
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File output_bam = "${base_name}.bam"
		File output_bam_index = "${base_name}.bai"
	}
	command <<<

		~{gatk_path} \
		SplitNCigarReads \
		-R ~{ref_fasta} \
		-I ~{input_bam} \
		-O ~{base_name}.bam

	>>>
	runtime {
		# disks: "local-disk " + sub(((size(input_bam, "GB") + 1) * 5 + size(ref_fasta, "GB")), "\..*", "") + " HDD"
		docker: docker
		memory: "4 GB"
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task HaplotypeCaller {
	input {
		File input_bam
		File input_bam_index
		String base_name
		File interval_list
		File ref_dict
		File ref_fasta
		File ref_fasta_index
		File dbSNP_vcf
		File dbSNP_vcf_index
		String gatk_path
		String docker
		Int preemptible_count
		Int? stand_call_conf
	}


	output {
		File output_vcf = "${base_name}.vcf.gz"
		File output_vcf_index = "${base_name}.vcf.gz.tbi"
	}
	command <<<

		~{gatk_path} --java-options "-Xms6000m -XX:GCTimeLimit=50 -XX:GCHeapFreeLimit=10" \
		HaplotypeCaller \
		-R ~{ref_fasta} \
		-I ~{input_bam} \
		-L ~{interval_list} \
		-O ~{base_name}.vcf.gz \
		-dont-use-soft-clipped-bases \
		--standard-min-confidence-threshold-for-calling ~{default="20"  stand_call_conf} \
		--dbsnp ~{dbSNP_vcf}

	>>>
	runtime {
		docker: docker
		memory: "6.5 GB"
		# disks: "local-disk " + sub((size(input_bam, "GB") * 2) + 30, "\..*", "") + " HDD"
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task SamToFastq {
	input {
		File unmapped_bam
		String base_name
		String gatk_path
		String docker
		Int preemptible_count
	}


	output {
		File fastq1 = "${base_name}.1.fastq.gz"
		File fastq2 = "${base_name}.2.fastq.gz"
	}
	command <<<

		~{gatk_path} \
		SamToFastq \
		--INPUT ~{unmapped_bam} \
		--VALIDATION_STRINGENCY SILENT \
		--FASTQ ~{base_name}.1.fastq.gz \
		--SECOND_END_FASTQ ~{base_name}.2.fastq.gz

	>>>
	runtime {
		docker: docker
		memory: "4 GB"
		# disks: "local-disk " + sub(((size(unmapped_bam, "GB") + 1) * 5), "\..*", "") + " HDD"
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
	}

}
task StarGenerateReferences {
	input {
		File ref_fasta
		File ref_fasta_index
		File annotations_gtf
		Int? read_length
		Int length = select_first([read_length, 101])
		Int? num_threads
		Int threads = select_first([num_threads, 8])
		#Int? additional_disk
		#Int add_to_disk = select_first([additional_disk, 0])
		#Int disk_size = select_first([100 + add_to_disk, 100])
		Int? mem_gb
		Int mem = select_first([100, mem_gb])
		String docker
		Int preemptible_count
	}


	output {
		Array[File] star_logs = glob("*.out")
		File star_genome_refs_zipped = "star-HUMAN-refs.tar.gz"
	}
	command <<<

		set -e
		mkdir STAR2_5

		STAR \
		--runMode genomeGenerate \
		--genomeDir STAR2_5 \
		--genomeFastaFiles ~{ref_fasta} \
		--sjdbGTFfile ~{annotations_gtf} \
		~{"--sjdbOverhang " + (length - 1)} \
		--runThreadN ~{threads}

		ls STAR2_5

		tar -zcvf star-HUMAN-refs.tar.gz STAR2_5

	>>>
	runtime {
		preemptible: preemptible_count
		awsBatchRetryAttempts: preemptible_count
		# disks: "local-disk " + disk_size + " HDD"
		docker: docker
		cpu: threads
		memory: mem + " GB"
	}

}