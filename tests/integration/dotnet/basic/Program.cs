// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;

class Program
{
    static Task<int> Main(string[] args)
    {
        return Deployment.RunAsync(() =>
        {
            var pod = new Pulumi.Kubernetes.Core.V1.Pod("pod", new Pod
            {
                Spec = new PodSpec
                {
                    Containers = 
                    {
                        new Container
                        {
                            Name = "nginx",
                            Image = "nginx",
                        },
                    },
                },
            });
        });
    }
}
