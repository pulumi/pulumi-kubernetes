# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
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
from ... import meta as _meta

__all__ = [
    'LeaseCandidateSpecPatchArgs',
    'LeaseCandidateSpecPatchArgsDict',
    'LeaseCandidateSpecArgs',
    'LeaseCandidateSpecArgsDict',
    'LeaseCandidateArgs',
    'LeaseCandidateArgsDict',
]

MYPY = False

if not MYPY:
    class LeaseCandidateSpecPatchArgsDict(TypedDict):
        """
        LeaseCandidateSpec is a specification of a Lease.
        """
        binary_version: NotRequired[pulumi.Input[str]]
        """
        BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
        """
        emulation_version: NotRequired[pulumi.Input[str]]
        """
        EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is "OldestEmulationVersion"
        """
        lease_name: NotRequired[pulumi.Input[str]]
        """
        LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
        """
        ping_time: NotRequired[pulumi.Input[str]]
        """
        PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
        """
        renew_time: NotRequired[pulumi.Input[str]]
        """
        RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
        """
        strategy: NotRequired[pulumi.Input[str]]
        """
        Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
        """
elif False:
    LeaseCandidateSpecPatchArgsDict: TypeAlias = Mapping[str, Any]

@pulumi.input_type
class LeaseCandidateSpecPatchArgs:
    def __init__(__self__, *,
                 binary_version: Optional[pulumi.Input[str]] = None,
                 emulation_version: Optional[pulumi.Input[str]] = None,
                 lease_name: Optional[pulumi.Input[str]] = None,
                 ping_time: Optional[pulumi.Input[str]] = None,
                 renew_time: Optional[pulumi.Input[str]] = None,
                 strategy: Optional[pulumi.Input[str]] = None):
        """
        LeaseCandidateSpec is a specification of a Lease.
        :param pulumi.Input[str] binary_version: BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
        :param pulumi.Input[str] emulation_version: EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is "OldestEmulationVersion"
        :param pulumi.Input[str] lease_name: LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
        :param pulumi.Input[str] ping_time: PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
        :param pulumi.Input[str] renew_time: RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
        :param pulumi.Input[str] strategy: Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
        """
        if binary_version is not None:
            pulumi.set(__self__, "binary_version", binary_version)
        if emulation_version is not None:
            pulumi.set(__self__, "emulation_version", emulation_version)
        if lease_name is not None:
            pulumi.set(__self__, "lease_name", lease_name)
        if ping_time is not None:
            pulumi.set(__self__, "ping_time", ping_time)
        if renew_time is not None:
            pulumi.set(__self__, "renew_time", renew_time)
        if strategy is not None:
            pulumi.set(__self__, "strategy", strategy)

    @property
    @pulumi.getter(name="binaryVersion")
    def binary_version(self) -> Optional[pulumi.Input[str]]:
        """
        BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
        """
        return pulumi.get(self, "binary_version")

    @binary_version.setter
    def binary_version(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "binary_version", value)

    @property
    @pulumi.getter(name="emulationVersion")
    def emulation_version(self) -> Optional[pulumi.Input[str]]:
        """
        EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is "OldestEmulationVersion"
        """
        return pulumi.get(self, "emulation_version")

    @emulation_version.setter
    def emulation_version(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "emulation_version", value)

    @property
    @pulumi.getter(name="leaseName")
    def lease_name(self) -> Optional[pulumi.Input[str]]:
        """
        LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
        """
        return pulumi.get(self, "lease_name")

    @lease_name.setter
    def lease_name(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "lease_name", value)

    @property
    @pulumi.getter(name="pingTime")
    def ping_time(self) -> Optional[pulumi.Input[str]]:
        """
        PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
        """
        return pulumi.get(self, "ping_time")

    @ping_time.setter
    def ping_time(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "ping_time", value)

    @property
    @pulumi.getter(name="renewTime")
    def renew_time(self) -> Optional[pulumi.Input[str]]:
        """
        RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
        """
        return pulumi.get(self, "renew_time")

    @renew_time.setter
    def renew_time(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "renew_time", value)

    @property
    @pulumi.getter
    def strategy(self) -> Optional[pulumi.Input[str]]:
        """
        Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
        """
        return pulumi.get(self, "strategy")

    @strategy.setter
    def strategy(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "strategy", value)


if not MYPY:
    class LeaseCandidateSpecArgsDict(TypedDict):
        """
        LeaseCandidateSpec is a specification of a Lease.
        """
        binary_version: pulumi.Input[str]
        """
        BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
        """
        lease_name: pulumi.Input[str]
        """
        LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
        """
        strategy: pulumi.Input[str]
        """
        Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
        """
        emulation_version: NotRequired[pulumi.Input[str]]
        """
        EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is "OldestEmulationVersion"
        """
        ping_time: NotRequired[pulumi.Input[str]]
        """
        PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
        """
        renew_time: NotRequired[pulumi.Input[str]]
        """
        RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
        """
elif False:
    LeaseCandidateSpecArgsDict: TypeAlias = Mapping[str, Any]

@pulumi.input_type
class LeaseCandidateSpecArgs:
    def __init__(__self__, *,
                 binary_version: pulumi.Input[str],
                 lease_name: pulumi.Input[str],
                 strategy: pulumi.Input[str],
                 emulation_version: Optional[pulumi.Input[str]] = None,
                 ping_time: Optional[pulumi.Input[str]] = None,
                 renew_time: Optional[pulumi.Input[str]] = None):
        """
        LeaseCandidateSpec is a specification of a Lease.
        :param pulumi.Input[str] binary_version: BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
        :param pulumi.Input[str] lease_name: LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
        :param pulumi.Input[str] strategy: Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
        :param pulumi.Input[str] emulation_version: EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is "OldestEmulationVersion"
        :param pulumi.Input[str] ping_time: PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
        :param pulumi.Input[str] renew_time: RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
        """
        pulumi.set(__self__, "binary_version", binary_version)
        pulumi.set(__self__, "lease_name", lease_name)
        pulumi.set(__self__, "strategy", strategy)
        if emulation_version is not None:
            pulumi.set(__self__, "emulation_version", emulation_version)
        if ping_time is not None:
            pulumi.set(__self__, "ping_time", ping_time)
        if renew_time is not None:
            pulumi.set(__self__, "renew_time", renew_time)

    @property
    @pulumi.getter(name="binaryVersion")
    def binary_version(self) -> pulumi.Input[str]:
        """
        BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
        """
        return pulumi.get(self, "binary_version")

    @binary_version.setter
    def binary_version(self, value: pulumi.Input[str]):
        pulumi.set(self, "binary_version", value)

    @property
    @pulumi.getter(name="leaseName")
    def lease_name(self) -> pulumi.Input[str]:
        """
        LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
        """
        return pulumi.get(self, "lease_name")

    @lease_name.setter
    def lease_name(self, value: pulumi.Input[str]):
        pulumi.set(self, "lease_name", value)

    @property
    @pulumi.getter
    def strategy(self) -> pulumi.Input[str]:
        """
        Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved. (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
        """
        return pulumi.get(self, "strategy")

    @strategy.setter
    def strategy(self, value: pulumi.Input[str]):
        pulumi.set(self, "strategy", value)

    @property
    @pulumi.getter(name="emulationVersion")
    def emulation_version(self) -> Optional[pulumi.Input[str]]:
        """
        EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is "OldestEmulationVersion"
        """
        return pulumi.get(self, "emulation_version")

    @emulation_version.setter
    def emulation_version(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "emulation_version", value)

    @property
    @pulumi.getter(name="pingTime")
    def ping_time(self) -> Optional[pulumi.Input[str]]:
        """
        PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
        """
        return pulumi.get(self, "ping_time")

    @ping_time.setter
    def ping_time(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "ping_time", value)

    @property
    @pulumi.getter(name="renewTime")
    def renew_time(self) -> Optional[pulumi.Input[str]]:
        """
        RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
        """
        return pulumi.get(self, "renew_time")

    @renew_time.setter
    def renew_time(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "renew_time", value)


if not MYPY:
    class LeaseCandidateArgsDict(TypedDict):
        """
        LeaseCandidate defines a candidate for a Lease object. Candidates are created such that coordinated leader election will pick the best leader from the list of candidates.
        """
        api_version: NotRequired[pulumi.Input[str]]
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        kind: NotRequired[pulumi.Input[str]]
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        metadata: NotRequired[pulumi.Input['_meta.v1.ObjectMetaArgsDict']]
        """
        More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        spec: NotRequired[pulumi.Input['LeaseCandidateSpecArgsDict']]
        """
        spec contains the specification of the Lease. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
elif False:
    LeaseCandidateArgsDict: TypeAlias = Mapping[str, Any]

@pulumi.input_type
class LeaseCandidateArgs:
    def __init__(__self__, *,
                 api_version: Optional[pulumi.Input[str]] = None,
                 kind: Optional[pulumi.Input[str]] = None,
                 metadata: Optional[pulumi.Input['_meta.v1.ObjectMetaArgs']] = None,
                 spec: Optional[pulumi.Input['LeaseCandidateSpecArgs']] = None):
        """
        LeaseCandidate defines a candidate for a Lease object. Candidates are created such that coordinated leader election will pick the best leader from the list of candidates.
        :param pulumi.Input[str] api_version: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        :param pulumi.Input[str] kind: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        :param pulumi.Input['_meta.v1.ObjectMetaArgs'] metadata: More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        :param pulumi.Input['LeaseCandidateSpecArgs'] spec: spec contains the specification of the Lease. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        if api_version is not None:
            pulumi.set(__self__, "api_version", 'coordination.k8s.io/v1alpha2')
        if kind is not None:
            pulumi.set(__self__, "kind", 'LeaseCandidate')
        if metadata is not None:
            pulumi.set(__self__, "metadata", metadata)
        if spec is not None:
            pulumi.set(__self__, "spec", spec)

    @property
    @pulumi.getter(name="apiVersion")
    def api_version(self) -> Optional[pulumi.Input[str]]:
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        return pulumi.get(self, "api_version")

    @api_version.setter
    def api_version(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "api_version", value)

    @property
    @pulumi.getter
    def kind(self) -> Optional[pulumi.Input[str]]:
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        return pulumi.get(self, "kind")

    @kind.setter
    def kind(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "kind", value)

    @property
    @pulumi.getter
    def metadata(self) -> Optional[pulumi.Input['_meta.v1.ObjectMetaArgs']]:
        """
        More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        return pulumi.get(self, "metadata")

    @metadata.setter
    def metadata(self, value: Optional[pulumi.Input['_meta.v1.ObjectMetaArgs']]):
        pulumi.set(self, "metadata", value)

    @property
    @pulumi.getter
    def spec(self) -> Optional[pulumi.Input['LeaseCandidateSpecArgs']]:
        """
        spec contains the specification of the Lease. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
        """
        return pulumi.get(self, "spec")

    @spec.setter
    def spec(self, value: Optional[pulumi.Input['LeaseCandidateSpecArgs']]):
        pulumi.set(self, "spec", value)

