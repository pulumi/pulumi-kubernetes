#!/bin/bash
set -ex

# Modify our resources but leave them un-ready.
kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-ready-condition -p '{"spec": {"someField": "touched"}}'
