# How-To Create the RPM Package

Create required rpmbuild directory structure in the users home directory:  

* `mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}`
* `mkdir -p ~/rpmbuild/SOURCES/prometheus-cluster-exporter-1.1/usr/bin/`
* `mkdir -p ~/rpmbuild/SOURCES/prometheus-cluster-exporter-1.1/usr/lib/systemd/system/`

Provide the following files to build the rpm package:  

* `~/rpmbuild/SOURCES/prometheus-cluster-exporter-1.1/usr/bin/cluster_exporter`
* `~/rpmbuild/SOURCES/prometheus-cluster-exporter-1.1/usr/lib/systemd/system/prometheus-cluster-exporter.service`
* `~/rpmbuild/SPECS/prometheus-cluster-exporter.spec`

Create the tar ball:  

* `cd ~/rpmbuild/SOURCES`
* `tar -czvf prometheus-cluster-exporter-1.1.tar.gz prometheus-cluster-exporter-1.1`

    Use relative paths here, otherwise rpmbuild will not find the extracted files!

Create RPM package:  

`rpmbuild -ba ~/rpmbuild/SPECS/prometheus-cluster-exporter.spec`

# Resources

* https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment
* https://rpm-packaging-guide.github.io/
* https://docs.fedoraproject.org/en-US/packaging-guidelines/RPMMacros/
