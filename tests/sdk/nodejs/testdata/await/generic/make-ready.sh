#!/bin/bash
set -ex

kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-ready-condition -p '{"status": {"conditions": [{"type": "Ready", "status": "True"}]}}'
kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-generation-increment -p '{"status": {"observedGeneration": 6}}'
kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-foo-condition -p '{"status": {"conditions": [{"type": "Foo", "status": "True"}]}}'
kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-field -p '{"spec": {"someField": "foo"}}'
kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-field-and-foo-condition -p '{"spec": {"someField": "expected"}}'
kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-field-and-foo-condition -p '{"status": {"conditions": [{"type": "Foo", "status": "True"}]}}'
