from pulumi_kubernetes.core.v1 import Pod, ConfigMap

# pod = Pod("foo")
cm = ConfigMap('test', data={"foo": "bar"})
