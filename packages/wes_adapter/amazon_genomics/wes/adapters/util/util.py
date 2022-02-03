def describe_batch_jobs_with_tag(tag_key, tag_value, aws_batch, aws_tags):
    """
    Retrieve descriptions of all Batch jobs with the given tag
    """
    pagination_token = None
    all_descriptions = []
    get_resources_kwargs = {
        "TagFilters": [{"Key": tag_key, "Values": [tag_value]}],
        "ResourceTypeFilters": ["batch:job"],
    }
    while True:
        if pagination_token:
            get_resources_kwargs["PaginationToken"] = pagination_token
        resources = aws_tags.get_resources(**get_resources_kwargs)
        resource_tag_mappings = resources.get("ResourceTagMappingList", [])
        job_arns = map(
            lambda tag_mapping: tag_mapping["ResourceARN"], resource_tag_mappings
        )
        job_ids = list(map(job_id_from_arn, job_arns))
        if job_ids:
            descriptions = aws_batch.describe_jobs(jobs=job_ids)["jobs"]
            all_descriptions += descriptions
        pagination_token = resources.get("PaginationToken", None)
        if not pagination_token:
            return all_descriptions


def job_id_from_arn(job_arn: str) -> str:
    return job_arn[job_arn.rindex("/") + 1:]
