#!/bin/sh
rm -rf vendor/github.com/pulumi/lumi
cp -R ../lumi vendor/github.com/pulumi/lumi
rm -rf vendor/github.com/pulumi/aws
cp -R packs/aws vendor/github.com/pulumi/aws

