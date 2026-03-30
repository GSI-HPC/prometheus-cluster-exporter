Name:           prometheus-cluster-exporter
Version:        1.1.7
Release:        1%{?dist}
Summary:        Prometheus exporter for Lustre IO throughput metrics associated to SLURM accounts and processes on a cluster.
Group:          Monitoring

License:        GPL-3.0-only
URL:            https://github.com/GSI-HPC/prometheus-cluster-exporter
Source0:        %{name}-%{version}.tar.gz
Source1:        %{name}.service
Source2:        %{name}.options

%{?systemd_requires}
BuildRequires:  systemd-rpm-macros
BuildRequires:  golang
Requires(pre): shadow-utils

%description
A Prometheus exporter for Lustre metadata operations and IO throughput metrics associated to SLURM accounts
and process names with user and group information on a cluster.

%global debug_package %{nil}

%prep
%setup -q

%build
GO111MODULE=on
go build -o prometheus-cluster-exporter ./

%install
install -Dm0755 prometheus-cluster-exporter %{buildroot}%{_bindir}/prometheus-cluster-exporter
install -Dm0644 %{SOURCE1} %{buildroot}%{_unitdir}/prometheus-cluster-exporter.service
install -Dm0644 %{SOURCE2} %{buildroot}/etc/default/prometheus-cluster-exporter.options

%check
go test ./...

%pre
getent group prometheus >/dev/null || groupadd -r prometheus
getent passwd prometheus >/dev/null || \
    useradd -r -g prometheus -d /var/lib/prometheus-cluster-exporter -s /sbin/nologin \
    -c "Prometheus exporter user" prometheus
exit 0

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%files
%defattr(-,root,root,-)
%config(noreplace) /etc/default/prometheus-cluster-exporter.options
%{_bindir}/prometheus-cluster-exporter
%{_unitdir}/prometheus-cluster-exporter.service
