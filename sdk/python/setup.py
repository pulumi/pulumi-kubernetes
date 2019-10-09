from subprocess import check_call

from setuptools import setup, find_packages
from setuptools.command.install import install


class InstallPluginCommand(install):
    def run(self):
        install.run(self)
        check_call(['pulumi', 'plugin', 'install', 'resource', 'kubernetes', '${PLUGIN_VERSION}'])


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
          'pulumi>=1.3.0,<2.0.0',
          'requests>=2.21.0,<2.22.0',
          'pyyaml>=5.1,<5.2',
          'semver>=2.8.1',
          'parver>=0.2.1',
      ],
      zip_safe=False)
