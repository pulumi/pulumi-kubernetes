// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Pulumi;
using CoreV1 = Pulumi.Kubernetes.Core.V1;
using Yaml = Pulumi.Kubernetes.Yaml;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;

class Program
{
    public static string Base64Encode(string plainText)
    {
        var plainTextBytes = System.Text.Encoding.UTF8.GetBytes(plainText);
        return System.Convert.ToBase64String(plainTextBytes);
    }
    private static Random random = new Random();
    public static string RandomString()
    {
        const string chars = "abcdefghijklmnopqrstuvwxyz";
        return new string(Enumerable.Repeat(chars, 5)
          .Select(s => s[random.Next(s.Length)]).ToArray());
    }

    static Task<int> Main(string[] args)
    {
        return Pulumi.Deployment.RunAsync(() =>
        {
            var config = new Pulumi.Config();
            var pw = config.RequireSecret("message");
            var rawPw = config.Require("message");

            var cmData = new CoreV1.ConfigMap("cmdata", new ConfigMapArgs
            {
                Data = new InputMap<string>
                {
                    {"password", pw},
                }
            });

            var cmBinaryData = new CoreV1.ConfigMap("cmbinarydata", new ConfigMapArgs
            {
                BinaryData = new InputMap<string>
                {
                    {"password", pw.Apply(v => Base64Encode(v))},
                }
            });

            var sStringData = new CoreV1.Secret("sstringdata", new SecretArgs
            {
                StringData = new InputMap<string>
                {
                    {"password", rawPw}
                }
            });

            var sData = new CoreV1.Secret("sdata", new SecretArgs
            {
                Data = new InputMap<string>
                {
                    {"password", Base64Encode(rawPw)}
                }
            });

            var name = $"test-{RandomString()}";
            var secretYAML = $@"
apiVersion: v1
kind: Secret
metadata:
  name: {name}
stringData:
  password: {rawPw}
";
            var cg = new Yaml.ConfigGroup("example", new Yaml.ConfigGroupArgs
            {
                Yaml = secretYAML
            });
            var cgSecret = cg.GetResource<CoreV1.Secret>(name);

            return new Dictionary<string, object>
            {
                {"cmData", cmData.Data},
                {"cmBinaryData", cmData.BinaryData},
                {"sStringData", sStringData.StringData},
                {"sData", sStringData.Data},
                {"cgData", cgSecret.Apply(v => v.Data)},
            };

        });
    }
}
