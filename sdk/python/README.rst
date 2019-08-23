|Build Status|

Pulumi Kubernetes Resource Provider
===================================

The Kubernetes resource provider for Pulumi lets you create, deploy, and
manage Kubernetes API resources and workloads in a running cluster.

-  `Introduction <#introduction>`__
-  `Kubernetes API Version Support <#kubernetes-api-version-support>`__
-  `How does API support for Kubernetes
   work? <#how-does-api-support-for-kubernetes-work>`__
-  `References <#references>`__
-  `Prerequisites <#prerequisites>`__
-  `Installing <#installing>`__
-  `Quick Examples <#quick-examples>`__
-  `Deploying a YAML Manifest <#deploying-a-yaml-manifest>`__
-  `Deploying a Helm Chart <#deploying-a-helm-chart>`__
-  `Deploying a Workload using the Resource
   API <#deploying-a-workload-using-the-resource-api>`__
-  `Contributing <#contributing>`__
-  `Code of Conduct <#code-of-conduct>`__

Introduction
------------

``pulumi-kubernetes`` provides an SDK to create any of the API resources
available in Kubernetes.

This includes the resources you know and love, such as: - Deployments -
ReplicaSets - ConfigMaps - Secrets - Jobs etc.

Kubernetes API Version Support
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

The ``pulumi-kubernetes`` SDK closely tracks the latest upstream
release, and provides access to the full API surface, including
deprecated endpoints. The SDK API is 100% compatible with the Kubernetes
API, and is schematically identical to what Kubernetes users expect.

At a minimum, we support the `current Kubernetes
version <https://kubernetes.io/docs/setup/release/>`__ and the previous
two versions. Older versions will likely work as well, but are not
officially supported.

See the
`CHANGELOG <https://github.com/pulumi/pulumi-kubernetes/blob/master/CHANGELOG.md>`__
for details on supported versions of Kubernetes for each version of the
``pulumi-kubernetes`` package.

How does API support for Kubernetes work?
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Pulumi’s Kubernetes SDK is manufactured by automatically wrapping our
library functionality around the Kubernetes resource `OpenAPI
spec <https://github.com/kubernetes/kubernetes/tree/master/api/openapi-spec>`__
as soon as a new version is released! Ultimately, this means that Pulumi
users do not have to learn a new Kubernetes API model, nor wait long to
work with the latest available versions.

    Note: Pulumi also supports alpha and beta APIs.

Visit the
`FAQ <https://www.pulumi.com/docs/reference/clouds/kubernetes/faq/>`__
for more details.

References
----------

-  `Reference
   Documentation <https://www.pulumi.com/docs/reference/clouds/kubernetes/>`__
-  API Documentation

   -  `Node.js
      API <https://pulumi.io/reference/pkg/nodejs/@pulumi/kubernetes>`__
   -  `Python
      API <https://www.pulumi.com/docs/reference/pkg/python/pulumi_kubernetes/>`__

-  `All Examples <./examples>`__
-  `Tutorials <https://www.pulumi.com/docs/reference/tutorials/kubernetes/>`__

Prerequisites
-------------

See the
`quickstart <https://www.pulumi.com/docs/get-started/kubernetes/>`__ for
complete documentation.

1. `Install
   Pulumi <https://www.pulumi.com/docs/quickstart/kubernetes/install-pulumi/>`__
2. Install a language runtime such as
   `Node.js <https://nodejs.org/en/download>`__ or
   `Python <https://www.python.org/downloads/>`__.
3. Install a package manager

   -  For Node.js, use `NPM <https://www.npmjs.com/get-npm>`__ or
      `Yarn <https://yarnpkg.com/lang/en/docs/install>`__.
   -  For Python, use
      `pip <https://pip.pypa.io/en/stable/installing/>`__.

4. Have access to a running Kubernetes cluster

   -  If ``kubectl`` already works for your running cluster, Pulumi
      respects and uses this configuration.
   -  If you do not have a cluster already running and available, we
      encourage you to explore Pulumi's SDKs for AWS EKS, Azure AKS, and
      GCP GKE. Visit the `reference
      docs <https://www.pulumi.com/docs/reference/clouds/kubernetes/>`__
      for more details.

5. `Install
   ``kubectl`` <https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl>`__.

Installing
----------

This package is available in JavaScript/TypeScript for use with Node.js,
as well as in Python.

For Node.js use either ``npm`` or ``yarn``:

``npm``:

.. code:: bash

    npm install @pulumi/kubernetes

``yarn``:

.. code:: bash

    yarn add @pulumi/kubernetes

For Python use ``pip``:

.. code:: bash

    pip install pulumi-kubernetes

Quick Examples
--------------

The following examples demonstrate how to work with
``pulumi-kubernetes`` in a couple of ways.

Examples may include the creation of an AWS EKS cluster, although an EKS
cluster is **not** required to use ``pulumi/kubernetes``. It is simply
used to ensure we have access to a running Kubernetes cluster to deploy
resources and workloads into.

Deploying a YAML Manifest
~~~~~~~~~~~~~~~~~~~~~~~~~

This example deploys resources from a YAML manifest file path, using the
transient, default ``kubeconfig`` credentials on the local machine, just
as ``kubectl`` does.

    Note: This capabality is primarily targeted for experimentation and
    transitioning to Pulumi. Pulumi's `desired state
    model <https://www.pulumi.com/docs/intro/concepts/how-pulumi-works>`__
    greatly benefits from having resources be directly defined in your
    Pulumi program as demonstrated in the `workload
    example <#deploying-a-workload-on-aws-eks>`__.

.. code:: typescript

    import * as k8s as "@pulumi/kubernetes";

    const myApp = new k8s.yaml.ConfigFile("app", {
        file: "app.yaml"
    });

Deploying a Helm Chart
~~~~~~~~~~~~~~~~~~~~~~

This example creates an EKS cluster with
```pulumi/eks`` <https://github.com/pulumi/pulumi-eks>`__, and then
deploys a Helm chart from the stable repo using the ``kubeconfig``
credentials from the cluster's `Pulumi
provider <https://www.pulumi.com/docs/reference/programming-model/#providers>`__.

    Note: This capabality is primarily targeted for experimentation and
    transitioning to Pulumi. Pulumi's `desired state
    model <https://www.pulumi.com/docs/intro/concepts/how-pulumi-works>`__
    greatly benefits from having resources be directly defined in your
    Pulumi program as demonstrated in the `workload
    example <#deploying-a-workload-on-aws-eks>`__.

.. code:: typescript

    import * as eks from "@pulumi/eks";
    import * as k8s from "@pulumi/kubernetes";

    // Create an EKS cluster.
    const cluster = new eks.Cluster("my-cluster");

    // Deploy Wordpress into our cluster.
    const wordpress = new k8s.helm.v2.Chart("wordpress", {
        repo: "stable",
        chart: "wordpress",
        values: {
            wordpressBlogName: "My Cool Kubernetes Blog!",
        },
    }, { providers: { "kubernetes": cluster.provider } });

    // Export the cluster's kubeconfig.
    export const kubeconfig = cluster.kubeconfig;

Deploying a Workload using the Resource API
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This example creates a EKS cluster with
```pulumi/eks`` <https://github.com/pulumi/pulumi-eks>`__, and then
deploys an NGINX Deployment and Service using the SDK resource API, and
the ``kubeconfig`` credentials from the cluster's `Pulumi
provider <https://www.pulumi.com/docs/reference/programming-model/#providers>`__.

.. code:: typescript

    import * as eks from "@pulumi/eks";
    import * as k8s from "@pulumi/kubernetes";

    // Create an EKS cluster with the default configuration.
    const cluster = new eks.Cluster("my-cluster");

    // Create a NGINX Deployment and Service.
    const appName = "my-app";
    const appLabels = { appClass: appName };
    const deployment = new k8s.apps.v1.Deployment(`${appName}-dep`, {
        metadata: { labels: appLabels },
        spec: {
            replicas: 2,
            selector: { matchLabels: appLabels },
            template: {
                metadata: { labels: appLabels },
                spec: {
                    containers: [{
                        name: appName,
                        image: "nginx",
                        ports: [{ name: "http", containerPort: 80 }]
                    }],
                }
            }
        },
    }, { provider: cluster.provider });

    const service = new k8s.core.v1.Service(`${appName}-svc`, {
        metadata: { labels: appLabels },
        spec: {
            type: "LoadBalancer",
            ports: [{ port: 80, targetPort: "http" }],
            selector: appLabels,
        },
    }, { provider: cluster.provider });

    // Export the URL for the load balanced service.
    export const url = service.status.loadBalancer.ingress[0].hostname;

    // Export the cluster's kubeconfig.
    export const kubeconfig = cluster.kubeconfig;

Contributing
------------

If you are interested in contributing, please see the `contributing
docs <CONTRIBUTING.md>`__.

Code of Conduct
---------------

You can read the code of conduct `here <CODE-OF-CONDUCT.md>`__.

.. |Build Status| image:: https://travis-ci.com/pulumi/pulumi-kubernetes.svg?token=eHg7Zp5zdDDJfTjY8ejq&branch=master
   :target: https://travis-ci.com/pulumi/pulumi-kubernetes
