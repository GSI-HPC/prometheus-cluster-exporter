# Prometheus Cluster Exporter

A Prometheus exporter for Lustre IO throughput metrics associated to SLURM accounts and process names on a cluster.

## Metrics

Cluster exporter metrics are prefixed with "cluster_".

### Global

These metrics are always exported.  

| Metric                           | Labels        | Description                                                       |
| -------------------------------- | ------------- | ----------------------------------------------------------------- |
| exporter_scrape_ok               | -             | Indicates if the scrape of the exporter was successful or not.    |
| exporter_stage_execution_seconds | name          | Execution duration in seconds spend in a specific exporter stage. |

TODO: List/Describe exporter stages.  

    HELP cluster_exporter_stage_execution_seconds Execution duration in seconds spend in a specific exporter stage.  
    cluster_exporter_stage_execution_seconds{name="build_read_throughput_metrics"} 3.524443994  
    cluster_exporter_stage_execution_seconds{name="build_write_throughput_metrics"} 3.603413776  
    cluster_exporter_stage_execution_seconds{name="retrieve_running_jobs"} 0.161250696  

### Throughput

#### **Jobs**

| Metric                     | Labels        | Description                                                                           |
| -------------------------- | ------------- | ------------------------------------------------------------------------------------- |
| job_read_throughput_bytes  | account, user | Total IO read throughput of all jobs on the cluster per account in bytes per second.  |
| job_write_throughput_bytes | account, user | Total IO write throughput of all jobs on the cluster per account in bytes per second. |

#### **Process Names**

| Metric                      | Labels         | Description                                                                            |
| --------------------------- | -------------- | -------------------------------------------------------------------------------------- |
| proc_read_throughput_bytes  | proc_name, uid | Total IO read throughput of process names on the cluster per uid in bytes per second.  |
| proc_write_throughput_bytes | proc_name, uid | Total IO write throughput of process names on the cluster per uid in bytes per second. |
