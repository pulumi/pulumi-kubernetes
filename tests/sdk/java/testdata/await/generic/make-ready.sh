#!/bin/bash
set -ex

kubectl patch --type=merge -n=generic-await "genericawaiters.test.pulumi.com" wants-ready-condition -p '{"status": {"conditions": [{"type": "Ready", "status": "True"}]}}'
