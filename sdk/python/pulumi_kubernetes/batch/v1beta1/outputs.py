# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from ... import _utilities
from . import outputs
from ... import batch as _batch
from ... import core as _core
from ... import meta as _meta

__all__ = [
    'CronJob',
    'CronJobSpec',
    'CronJobSpecPatch',
    'CronJobStatus',
    'CronJobStatusPatch',
    'JobTemplateSpec',
    'JobTemplateSpecPatch',
]

@pulumi.output_type
class CronJob(dict):
    """
    CronJob represents the configuration of a single cron job.
    """
    @staticmethod
    def __key_warning(key: str):
        suggest = None
        if key == "apiVersion":
            suggest = "api_version"

        if suggest:
            pulumi.log.warn(f"Key '{key}' not found in CronJob. Access the value via the '{suggest}' property getter instead.")

    def __getitem__(self, key: str) -> Any:
        CronJob.__key_warning(key)
        return super().__getitem__(key)

    def get(self, key: str, default = None) -> Any:
        CronJob.__key_warning(key)
        return super().get(key, default)

    def __init__(__self__, *,
                 api_version: Optional[_builtins.str] = None,
                 kind: Optional[_builtins.str] = None,
                 metadata: Optional['_meta.v1.outputs.ObjectMeta'] = None,
                 spec: Optional['outputs.CronJobSpec'] = None,
                 status: Optional['outputs.CronJobStatus'] = None):
        """
        CronJob represents the configuration of a single cron job.
        :param _builtins.str api_version: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        :param _builtins.str kind: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        :param '_meta.v1.ObjectMetaArgs' metadata: Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        :param 'CronJobSpecArgs' spec: Specification of the desired behavior of a cron job, including the schedule. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        :param 'CronJobStatusArgs' status: Current status of a cron job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        if api_version is not None:
            pulumi.set(__self__, "api_version", 'batch/v1beta1')
        if kind is not None:
            pulumi.set(__self__, "kind", 'CronJob')
        if metadata is not None:
            pulumi.set(__self__, "metadata", metadata)
        if spec is not None:
            pulumi.set(__self__, "spec", spec)
        if status is not None:
            pulumi.set(__self__, "status", status)

    @_builtins.property
    @pulumi.getter(name="apiVersion")
    def api_version(self) -> Optional[_builtins.str]:
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        return pulumi.get(self, "api_version")

    @_builtins.property
    @pulumi.getter
    def kind(self) -> Optional[_builtins.str]:
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        return pulumi.get(self, "kind")

    @_builtins.property
    @pulumi.getter
    def metadata(self) -> Optional['_meta.v1.outputs.ObjectMeta']:
        """
        Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        return pulumi.get(self, "metadata")

    @_builtins.property
    @pulumi.getter
    def spec(self) -> Optional['outputs.CronJobSpec']:
        """
        Specification of the desired behavior of a cron job, including the schedule. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        return pulumi.get(self, "spec")

    @_builtins.property
    @pulumi.getter
    def status(self) -> Optional['outputs.CronJobStatus']:
        """
        Current status of a cron job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        return pulumi.get(self, "status")


@pulumi.output_type
class CronJobSpec(dict):
    """
    CronJobSpec describes how the job execution will look like and when it will actually run.
    """
    @staticmethod
    def __key_warning(key: str):
        suggest = None
        if key == "jobTemplate":
            suggest = "job_template"
        elif key == "concurrencyPolicy":
            suggest = "concurrency_policy"
        elif key == "failedJobsHistoryLimit":
            suggest = "failed_jobs_history_limit"
        elif key == "startingDeadlineSeconds":
            suggest = "starting_deadline_seconds"
        elif key == "successfulJobsHistoryLimit":
            suggest = "successful_jobs_history_limit"

        if suggest:
            pulumi.log.warn(f"Key '{key}' not found in CronJobSpec. Access the value via the '{suggest}' property getter instead.")

    def __getitem__(self, key: str) -> Any:
        CronJobSpec.__key_warning(key)
        return super().__getitem__(key)

    def get(self, key: str, default = None) -> Any:
        CronJobSpec.__key_warning(key)
        return super().get(key, default)

    def __init__(__self__, *,
                 job_template: 'outputs.JobTemplateSpec',
                 schedule: _builtins.str,
                 concurrency_policy: Optional[_builtins.str] = None,
                 failed_jobs_history_limit: Optional[_builtins.int] = None,
                 starting_deadline_seconds: Optional[_builtins.int] = None,
                 successful_jobs_history_limit: Optional[_builtins.int] = None,
                 suspend: Optional[_builtins.bool] = None):
        """
        CronJobSpec describes how the job execution will look like and when it will actually run.
        :param 'JobTemplateSpecArgs' job_template: Specifies the job that will be created when executing a CronJob.
        :param _builtins.str schedule: The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
        :param _builtins.str concurrency_policy: Specifies how to treat concurrent executions of a Job. Valid values are: - "Allow" (default): allows CronJobs to run concurrently; - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet; - "Replace": cancels currently running job and replaces it with a new one
        :param _builtins.int failed_jobs_history_limit: The number of failed finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1.
        :param _builtins.int starting_deadline_seconds: Optional deadline in seconds for starting the job if it misses scheduled time for any reason.  Missed jobs executions will be counted as failed ones.
        :param _builtins.int successful_jobs_history_limit: The number of successful finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 3.
        :param _builtins.bool suspend: This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.
        """
        pulumi.set(__self__, "job_template", job_template)
        pulumi.set(__self__, "schedule", schedule)
        if concurrency_policy is not None:
            pulumi.set(__self__, "concurrency_policy", concurrency_policy)
        if failed_jobs_history_limit is not None:
            pulumi.set(__self__, "failed_jobs_history_limit", failed_jobs_history_limit)
        if starting_deadline_seconds is not None:
            pulumi.set(__self__, "starting_deadline_seconds", starting_deadline_seconds)
        if successful_jobs_history_limit is not None:
            pulumi.set(__self__, "successful_jobs_history_limit", successful_jobs_history_limit)
        if suspend is not None:
            pulumi.set(__self__, "suspend", suspend)

    @_builtins.property
    @pulumi.getter(name="jobTemplate")
    def job_template(self) -> 'outputs.JobTemplateSpec':
        """
        Specifies the job that will be created when executing a CronJob.
        """
        return pulumi.get(self, "job_template")

    @_builtins.property
    @pulumi.getter
    def schedule(self) -> _builtins.str:
        """
        The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
        """
        return pulumi.get(self, "schedule")

    @_builtins.property
    @pulumi.getter(name="concurrencyPolicy")
    def concurrency_policy(self) -> Optional[_builtins.str]:
        """
        Specifies how to treat concurrent executions of a Job. Valid values are: - "Allow" (default): allows CronJobs to run concurrently; - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet; - "Replace": cancels currently running job and replaces it with a new one
        """
        return pulumi.get(self, "concurrency_policy")

    @_builtins.property
    @pulumi.getter(name="failedJobsHistoryLimit")
    def failed_jobs_history_limit(self) -> Optional[_builtins.int]:
        """
        The number of failed finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1.
        """
        return pulumi.get(self, "failed_jobs_history_limit")

    @_builtins.property
    @pulumi.getter(name="startingDeadlineSeconds")
    def starting_deadline_seconds(self) -> Optional[_builtins.int]:
        """
        Optional deadline in seconds for starting the job if it misses scheduled time for any reason.  Missed jobs executions will be counted as failed ones.
        """
        return pulumi.get(self, "starting_deadline_seconds")

    @_builtins.property
    @pulumi.getter(name="successfulJobsHistoryLimit")
    def successful_jobs_history_limit(self) -> Optional[_builtins.int]:
        """
        The number of successful finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 3.
        """
        return pulumi.get(self, "successful_jobs_history_limit")

    @_builtins.property
    @pulumi.getter
    def suspend(self) -> Optional[_builtins.bool]:
        """
        This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.
        """
        return pulumi.get(self, "suspend")


@pulumi.output_type
class CronJobSpecPatch(dict):
    """
    CronJobSpec describes how the job execution will look like and when it will actually run.
    """
    @staticmethod
    def __key_warning(key: str):
        suggest = None
        if key == "concurrencyPolicy":
            suggest = "concurrency_policy"
        elif key == "failedJobsHistoryLimit":
            suggest = "failed_jobs_history_limit"
        elif key == "jobTemplate":
            suggest = "job_template"
        elif key == "startingDeadlineSeconds":
            suggest = "starting_deadline_seconds"
        elif key == "successfulJobsHistoryLimit":
            suggest = "successful_jobs_history_limit"

        if suggest:
            pulumi.log.warn(f"Key '{key}' not found in CronJobSpecPatch. Access the value via the '{suggest}' property getter instead.")

    def __getitem__(self, key: str) -> Any:
        CronJobSpecPatch.__key_warning(key)
        return super().__getitem__(key)

    def get(self, key: str, default = None) -> Any:
        CronJobSpecPatch.__key_warning(key)
        return super().get(key, default)

    def __init__(__self__, *,
                 concurrency_policy: Optional[_builtins.str] = None,
                 failed_jobs_history_limit: Optional[_builtins.int] = None,
                 job_template: Optional['outputs.JobTemplateSpecPatch'] = None,
                 schedule: Optional[_builtins.str] = None,
                 starting_deadline_seconds: Optional[_builtins.int] = None,
                 successful_jobs_history_limit: Optional[_builtins.int] = None,
                 suspend: Optional[_builtins.bool] = None):
        """
        CronJobSpec describes how the job execution will look like and when it will actually run.
        :param _builtins.str concurrency_policy: Specifies how to treat concurrent executions of a Job. Valid values are: - "Allow" (default): allows CronJobs to run concurrently; - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet; - "Replace": cancels currently running job and replaces it with a new one
        :param _builtins.int failed_jobs_history_limit: The number of failed finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1.
        :param 'JobTemplateSpecPatchArgs' job_template: Specifies the job that will be created when executing a CronJob.
        :param _builtins.str schedule: The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
        :param _builtins.int starting_deadline_seconds: Optional deadline in seconds for starting the job if it misses scheduled time for any reason.  Missed jobs executions will be counted as failed ones.
        :param _builtins.int successful_jobs_history_limit: The number of successful finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 3.
        :param _builtins.bool suspend: This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.
        """
        if concurrency_policy is not None:
            pulumi.set(__self__, "concurrency_policy", concurrency_policy)
        if failed_jobs_history_limit is not None:
            pulumi.set(__self__, "failed_jobs_history_limit", failed_jobs_history_limit)
        if job_template is not None:
            pulumi.set(__self__, "job_template", job_template)
        if schedule is not None:
            pulumi.set(__self__, "schedule", schedule)
        if starting_deadline_seconds is not None:
            pulumi.set(__self__, "starting_deadline_seconds", starting_deadline_seconds)
        if successful_jobs_history_limit is not None:
            pulumi.set(__self__, "successful_jobs_history_limit", successful_jobs_history_limit)
        if suspend is not None:
            pulumi.set(__self__, "suspend", suspend)

    @_builtins.property
    @pulumi.getter(name="concurrencyPolicy")
    def concurrency_policy(self) -> Optional[_builtins.str]:
        """
        Specifies how to treat concurrent executions of a Job. Valid values are: - "Allow" (default): allows CronJobs to run concurrently; - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet; - "Replace": cancels currently running job and replaces it with a new one
        """
        return pulumi.get(self, "concurrency_policy")

    @_builtins.property
    @pulumi.getter(name="failedJobsHistoryLimit")
    def failed_jobs_history_limit(self) -> Optional[_builtins.int]:
        """
        The number of failed finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1.
        """
        return pulumi.get(self, "failed_jobs_history_limit")

    @_builtins.property
    @pulumi.getter(name="jobTemplate")
    def job_template(self) -> Optional['outputs.JobTemplateSpecPatch']:
        """
        Specifies the job that will be created when executing a CronJob.
        """
        return pulumi.get(self, "job_template")

    @_builtins.property
    @pulumi.getter
    def schedule(self) -> Optional[_builtins.str]:
        """
        The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
        """
        return pulumi.get(self, "schedule")

    @_builtins.property
    @pulumi.getter(name="startingDeadlineSeconds")
    def starting_deadline_seconds(self) -> Optional[_builtins.int]:
        """
        Optional deadline in seconds for starting the job if it misses scheduled time for any reason.  Missed jobs executions will be counted as failed ones.
        """
        return pulumi.get(self, "starting_deadline_seconds")

    @_builtins.property
    @pulumi.getter(name="successfulJobsHistoryLimit")
    def successful_jobs_history_limit(self) -> Optional[_builtins.int]:
        """
        The number of successful finished jobs to retain. This is a pointer to distinguish between explicit zero and not specified. Defaults to 3.
        """
        return pulumi.get(self, "successful_jobs_history_limit")

    @_builtins.property
    @pulumi.getter
    def suspend(self) -> Optional[_builtins.bool]:
        """
        This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.
        """
        return pulumi.get(self, "suspend")


@pulumi.output_type
class CronJobStatus(dict):
    """
    CronJobStatus represents the current state of a cron job.
    """
    @staticmethod
    def __key_warning(key: str):
        suggest = None
        if key == "lastScheduleTime":
            suggest = "last_schedule_time"

        if suggest:
            pulumi.log.warn(f"Key '{key}' not found in CronJobStatus. Access the value via the '{suggest}' property getter instead.")

    def __getitem__(self, key: str) -> Any:
        CronJobStatus.__key_warning(key)
        return super().__getitem__(key)

    def get(self, key: str, default = None) -> Any:
        CronJobStatus.__key_warning(key)
        return super().get(key, default)

    def __init__(__self__, *,
                 active: Optional[Sequence['_core.v1.outputs.ObjectReference']] = None,
                 last_schedule_time: Optional[_builtins.str] = None):
        """
        CronJobStatus represents the current state of a cron job.
        :param Sequence['_core.v1.ObjectReferenceArgs'] active: A list of pointers to currently running jobs.
        :param _builtins.str last_schedule_time: Information when was the last time the job was successfully scheduled.
        """
        if active is not None:
            pulumi.set(__self__, "active", active)
        if last_schedule_time is not None:
            pulumi.set(__self__, "last_schedule_time", last_schedule_time)

    @_builtins.property
    @pulumi.getter
    def active(self) -> Optional[Sequence['_core.v1.outputs.ObjectReference']]:
        """
        A list of pointers to currently running jobs.
        """
        return pulumi.get(self, "active")

    @_builtins.property
    @pulumi.getter(name="lastScheduleTime")
    def last_schedule_time(self) -> Optional[_builtins.str]:
        """
        Information when was the last time the job was successfully scheduled.
        """
        return pulumi.get(self, "last_schedule_time")


@pulumi.output_type
class CronJobStatusPatch(dict):
    """
    CronJobStatus represents the current state of a cron job.
    """
    @staticmethod
    def __key_warning(key: str):
        suggest = None
        if key == "lastScheduleTime":
            suggest = "last_schedule_time"

        if suggest:
            pulumi.log.warn(f"Key '{key}' not found in CronJobStatusPatch. Access the value via the '{suggest}' property getter instead.")

    def __getitem__(self, key: str) -> Any:
        CronJobStatusPatch.__key_warning(key)
        return super().__getitem__(key)

    def get(self, key: str, default = None) -> Any:
        CronJobStatusPatch.__key_warning(key)
        return super().get(key, default)

    def __init__(__self__, *,
                 active: Optional[Sequence['_core.v1.outputs.ObjectReferencePatch']] = None,
                 last_schedule_time: Optional[_builtins.str] = None):
        """
        CronJobStatus represents the current state of a cron job.
        :param Sequence['_core.v1.ObjectReferencePatchArgs'] active: A list of pointers to currently running jobs.
        :param _builtins.str last_schedule_time: Information when was the last time the job was successfully scheduled.
        """
        if active is not None:
            pulumi.set(__self__, "active", active)
        if last_schedule_time is not None:
            pulumi.set(__self__, "last_schedule_time", last_schedule_time)

    @_builtins.property
    @pulumi.getter
    def active(self) -> Optional[Sequence['_core.v1.outputs.ObjectReferencePatch']]:
        """
        A list of pointers to currently running jobs.
        """
        return pulumi.get(self, "active")

    @_builtins.property
    @pulumi.getter(name="lastScheduleTime")
    def last_schedule_time(self) -> Optional[_builtins.str]:
        """
        Information when was the last time the job was successfully scheduled.
        """
        return pulumi.get(self, "last_schedule_time")


@pulumi.output_type
class JobTemplateSpec(dict):
    """
    JobTemplateSpec describes the data a Job should have when created from a template
    """
    def __init__(__self__, *,
                 metadata: Optional['_meta.v1.outputs.ObjectMeta'] = None,
                 spec: Optional['_batch.v1.outputs.JobSpec'] = None):
        """
        JobTemplateSpec describes the data a Job should have when created from a template
        :param '_meta.v1.ObjectMetaArgs' metadata: Standard object's metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        :param '_batch.v1.JobSpecArgs' spec: Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        if metadata is not None:
            pulumi.set(__self__, "metadata", metadata)
        if spec is not None:
            pulumi.set(__self__, "spec", spec)

    @_builtins.property
    @pulumi.getter
    def metadata(self) -> Optional['_meta.v1.outputs.ObjectMeta']:
        """
        Standard object's metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        return pulumi.get(self, "metadata")

    @_builtins.property
    @pulumi.getter
    def spec(self) -> Optional['_batch.v1.outputs.JobSpec']:
        """
        Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        return pulumi.get(self, "spec")


@pulumi.output_type
class JobTemplateSpecPatch(dict):
    """
    JobTemplateSpec describes the data a Job should have when created from a template
    """
    def __init__(__self__, *,
                 metadata: Optional['_meta.v1.outputs.ObjectMetaPatch'] = None,
                 spec: Optional['_batch.v1.outputs.JobSpecPatch'] = None):
        """
        JobTemplateSpec describes the data a Job should have when created from a template
        :param '_meta.v1.ObjectMetaPatchArgs' metadata: Standard object's metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        :param '_batch.v1.JobSpecPatchArgs' spec: Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        if metadata is not None:
            pulumi.set(__self__, "metadata", metadata)
        if spec is not None:
            pulumi.set(__self__, "spec", spec)

    @_builtins.property
    @pulumi.getter
    def metadata(self) -> Optional['_meta.v1.outputs.ObjectMetaPatch']:
        """
        Standard object's metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        return pulumi.get(self, "metadata")

    @_builtins.property
    @pulumi.getter
    def spec(self) -> Optional['_batch.v1.outputs.JobSpecPatch']:
        """
        Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        return pulumi.get(self, "spec")


