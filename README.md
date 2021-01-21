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

### Getent

The getent command is required for the uid to user and group mapping used for the process names throughput metrics.

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

| Metric                      | Labels                           | Description                                                                                       |
| --------------------------- | -------------------------------- | ------------------------------------------------------------------------------------------------- |
| proc_read_throughput_bytes  | proc_name, user_name, group_name | Total IO read throughput of process names on the cluster per group and user in bytes per second.  |
| proc_write_throughput_bytes | proc_name, user_name, group_name | Total IO write throughput of process names on the cluster per group and user in bytes per second. |

## Multiple Srape Prevention

Since the forked processes do not have a timeout handling, they might block for a uncertain amount of time.  
It is very unlikely that reexecuting the processes will solve the problem of beeing blocked.  
Therefore multiple scrapes at a time will be prevented by the exporter.  

The following warning will be displayed on afterward scrape executions, were a scrape is still active:  
    *"Collect is still active... - Skipping now"*

Besides that, the cluster_exporter_scrape_ok metric will be set to 0 for skipped scrape attempts.  

## Building the Exporter

```go
go build -o cluster-exporter *.go
```

## Known Issues / Future Work

### Running Multiple Cluster Exporter for a Lustre Instance

#### Description

Since the cluster exporter exports the to be processed Lustre Jobstats out of Prometheus, this has significant disadvantages, when multiple cluster exporter should be used for the same Lustre instance:

1. Costs of Resources  
If the amount of exported jobstats gets high, this will lead to timeouts in the build\_read/write\_throughput\_metrics stage.

2. Conflict with procname\_uid Jobstats  
If more than one cluster exporter is used and the Lustre clients are set to the jobstats option procname_uid for non-cluster machines e.g. submit-nodes or service machines that acccess Lustre directly,  
the exporter will process the same jobstats for processes each time.  This might result in to high calculated throughput values for the same scrape point by multiple cluster exporter.  


#### Solution

By using the Lustre Complex JobID feature it could be distinguished between different Lustre client groups e.g. compute cluster and service machine groups.  
With such an identification of groups, multiple cluster exporter could be used to just process a specific group.  

