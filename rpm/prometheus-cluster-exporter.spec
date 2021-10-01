%define        __spec_install_post %{nil}
%define          debug_package %{nil}
%define        __os_install_post %{_dbpath}/brp-compress

Name:           prometheus-cluster-exporter
Version:        1.1.3
Release:        1%{?dist}
Summary:        Prometheus exporter for Lustre IO throughput metrics associated to SLURM accounts and processes on a cluster.
Group:          Monitoring

License:        GPL 3.0
URL:            https://github.com/GSI-HPC/prometheus-cluster-exporter
Source0:        %{name}-%{version}.tar.gz

Requires(pre): shadow-utils

Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd
%{?systemd_requires}
BuildRequires:  systemd

BuildRoot:      %{_tmppath}/%{name}-%{version}-1-root

%description
A Prometheus exporter for Lustre metadata operations and IO throughput metrics associated to SLURM accounts
and process names with user and group information on a cluster.

%prep
%setup -q

%build
# Empty section.

%install
rm -rf %{buildroot}
mkdir -vp  %{buildroot}
mkdir -vp %{buildroot}%{_unitdir}/
cp usr/lib/systemd/system/%{name}.service %{buildroot}%{_unitdir}/

# in builddir
cp -a * %{buildroot}

%clean
rm -rf %{buildroot}
 
%pre
getent group prometheus >/dev/null || groupadd -r prometheus
getent passwd prometheus >/dev/null || \
    useradd -r -g prometheus -d /var/lib/prometheus-cluster-exporter -s /sbin/nologin \
    -c "Prometheus exporter user" prometheus
exit 0

%post
systemctl enable %{name}.service
systemctl start %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%files
%defattr(-,root,root,-)
%config /etc/sysconfig/prometheus-cluster-exporter.options
%{_bindir}/prometheus-cluster-exporter
%{_unitdir}/%{name}.service
