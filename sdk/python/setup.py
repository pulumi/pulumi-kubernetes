import errno

from subprocess import check_call

from setuptools import setup, find_packages
from setuptools.command.install import install

class InstallPluginCommand(install):
    def run(self):
        install.run(self)
        try:
            check_call(['pulumi', 'plugin', 'install', 'resource', 'kubernetes', '${PLUGIN_VERSION}'])
        except OSError as error:
            if error.errno == errno.ENOENT:
                print("""
                There was an error installing the kubernetes resource provider plugin.
                It looks like `pulumi` is not installed on your system.
                Please visit https://pulumi.com/ to install the Pulumi CLI.
                You may try manually installing the plugin by running
                `pulumi plugin install resource kubernetes ${PLUGIN_VERSION}`
                """)
            else:
                raise


def readme():
    with open('README.md', 'r') as f:
        return f.read()


setup(name='pulumi_kubernetes',
      version='${VERSION}',
      description='A Pulumi package for creating and managing Kubernetes resources.',
      long_description=readme(),
      long_description_content_type='text/markdown',
      cmdclass={
          'install': InstallPluginCommand,
      },
      keywords='pulumi kubernetes',
      url='https://pulumi.io',
      project_urls={
          'Repository': 'https://github.com/pulumi/pulumi-kubernetes'
      },
      license='Apache-2.0',
      packages=find_packages(),
      install_requires=[
          'pulumi>=1.11.0,<2.0.0',
          'requests>=2.21.0,<2.22.0',
          'semver>=2.8.1',
          'parver>=0.2.1',
      ],
      zip_safe=False)
