# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities

__all__ = [
    'KubeClientSettings',
]

@pulumi.output_type
class KubeClientSettings(dict):
    """
    Options for tuning the Kubernetes client used by a Provider.
    """
    def __init__(__self__, *,
                 burst: Optional[int] = None,
                 qps: Optional[float] = None):
        """
        Options for tuning the Kubernetes client used by a Provider.
        :param int burst: Maximum burst for throttle. Default value is 10.
        :param float qps: QPS indicates the maximum queries per second (QPS) to the API server from this client. Default value is 5.
        """
        if burst is None:
            burst = _utilities.get_env_int('PULUMI_K8S_CLIENT_BURST')
        if burst is not None:
            pulumi.set(__self__, "burst", burst)
        if qps is None:
            qps = _utilities.get_env_float('PULUMI_K8S_CLIENT_QPS')
        if qps is not None:
            pulumi.set(__self__, "qps", qps)

    @property
    @pulumi.getter
    def burst(self) -> Optional[int]:
        """
        Maximum burst for throttle. Default value is 10.
        """
        return pulumi.get(self, "burst")

    @property
    @pulumi.getter
    def qps(self) -> Optional[float]:
        """
        QPS indicates the maximum queries per second (QPS) to the API server from this client. Default value is 5.
        """
        return pulumi.get(self, "qps")

