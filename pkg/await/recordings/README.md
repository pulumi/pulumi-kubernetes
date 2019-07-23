This folder contains JSON-encoded recordings of k8s watch Events that can be
used to validate await logic for the provider. These recordings are organized
into the following structure:

* workflows - Each file contains a JSON array of watch Events corresponding to
a recorded workflow of interest.
* states - Each file contains a JSON-encoded watch Event corresponding to a
state of interest.
