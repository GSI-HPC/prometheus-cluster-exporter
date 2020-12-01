# How-To Create RPM Package

Create required rpmbuild directory structure in the users home directory:  

`mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}`

Create a directory e.g. RPM_SOURCE with following structure:  

* prometheus-cluster-exporter-1.0/usr/bin/cluster_exporter
* prometheus-cluster-exporter-1.0/usr/lib/systemd/system/prometheus-cluster-exporter.service

Create tar ball:  

`tar -czvf prometheus-cluster-exporter-1.0.tar.gz prometheus-cluster-exporter-1.0`

Copy tar ball to:  

`~/rpmbuild/SOURCES/`

Create RPM package:  

`rpmbuild -ba ~/rpmbuild/SPECS/prometheus-cluster-exporter.spec`

# Resources

* https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment
* https://rpm-packaging-guide.github.io/
* https://docs.fedoraproject.org/en-US/packaging-guidelines/RPMMacros/