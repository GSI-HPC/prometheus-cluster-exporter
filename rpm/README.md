# How-To Create the RPM Package

Create required rpmbuild directory structure in the users home directory:  

`mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}`

Create the following file structure within the `~/rpmbuild/SOURCES` directory to build the rpm package:  

* `prometheus-cluster-exporter-1.1/usr/bin/cluster_exporter`
* `prometheus-cluster-exporter-1.1/usr/lib/systemd/system/prometheus-cluster-exporter.service`

Create the tar ball:  

`tar -czvf prometheus-cluster-exporter-1.1.tar.gz prometheus-cluster-exporter-1.1`

    Use relative paths here, otherwise rpmbuild will not find the extracted files!

Create RPM package:  

`rpmbuild -ba ~/rpmbuild/SPECS/prometheus-cluster-exporter.spec`

# Resources

* https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment
* https://rpm-packaging-guide.github.io/
* https://docs.fedoraproject.org/en-US/packaging-guidelines/RPMMacros/
