# Prometheus Cluster Exporter

A Prometheus exporter for Lustre IO throughput metrics associated to SLURM accounts and process names on a cluster.

## Requirements

### Lustre Exporter

A Lustre exporter that exposes the two metrics to Prometheus with a label jobid is required:  

* lustre_job_read_bytes_total
* lustre_job_write_bytes_total

The Lustre exporter from HP provides such metrics:  
https://github.com/HewlettPackard/lustre_exporter

### Squeue Command

The squeue command from SLURM must be accessable locally to the exporter to retrieve the running jobs.

## Metrics

Cluster exporter metrics are prefixed with "cluster_".

### Global

These metrics are always exported.  

| Metric                           | Labels        | Description                                                       |
| -------------------------------- | ------------- | ----------------------------------------------------------------- |
| exporter_scrape_ok               | -             | Indicates if the scrape of the exporter was successful or not.    |
| exporter_stage_execution_seconds | name          | Execution duration in seconds spend in a specific exporter stage. |

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