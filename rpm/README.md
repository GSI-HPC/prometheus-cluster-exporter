# How-To Create the RPM Package

An automated Vagrant-based RPM builder is found in the subfolder `exporter-rpm-builder`. Using it is recommended (see the included README.md file for instructions). The following steps are for doing a manual build instead. 

Prerequisites: 

* Target version tag (1.1.7 in this example) needs to match contents of .spec file.
* Entire code needs to be in folder named prometheus-cluster-exporter-1.1.7, relative to current folder.

Create required rpmbuild directory structure in the users home directory:

* `mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}`

Create tarball for full source code (top level folder name is transformed to add version tag):

* `tar czf prometheus-cluster-exporter-1.1.7.tar.gz prometheus-cluster-exporter-1.1.7/`

Provide required files to build the rpm package:

* `cp prometheus-cluster-exporter-1.1.7.tar.gz ~/rpmbuild/SOURCES/`
* `cp prometheus-cluster-exporter-1.1.7/rpm/prometheus-cluster-exporter.spec ~/rpmbuild/SPECS/`
* `cp prometheus-cluster-exporter-1.1.7/systemd/prometheus-cluster-exporter.service ~/rpmbuild/SOURCES/`
* `cp prometheus-cluster-exporter-1.1.7/systemd/prometheus-cluster-exporter.options ~/rpmbuild/SOURCES/`

Create RPM package:

`rpmbuild -ba ~/rpmbuild/SPECS/prometheus-cluster-exporter.spec`

# Resources

* https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment
* https://rpm-packaging-guide.github.io/
* https://docs.fedoraproject.org/en-US/packaging-guidelines/RPMMacros/
