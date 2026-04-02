Name:           cinaseek
Version:        1.0.0
Release:        1%{?dist}
Summary:        CinaSeek - Lightweight VM Web Management Platform

License:        MIT
URL:            https://cinaseek.ai
Source0:        cinaseek-server-linux-amd64

ExclusiveArch:  x86_64

%description
CinaSeek is a lightweight VM remote management platform based on Multipass,
supporting OpenClaw AI development environment one-click deployment.

%prep
# No source prep needed for pre-built binaries

%install
mkdir -p %{buildroot}/usr/bin
install -m 0755 %{SOURCE0} %{buildroot}/usr/bin/cinaseek-server

mkdir -p %{buildroot}/etc/systemd/system
install -m 0644 cinaseek.service %{buildroot}/etc/systemd/system/cinaseek.service

%files
%attr(0755, root, root) /usr/bin/cinaseek-server
%config(noreplace) /etc/systemd/system/cinaseek.service

%post
systemctl daemon-reload
systemctl enable cinaseek

%preun
if [ $1 -eq 0 ]; then
    systemctl stop cinaseek || true
    systemctl disable cinaseek || true
fi

%postun
systemctl daemon-reload

%changelog
* Thu Apr 02 2026 CinaGroup <dev@cinagroup.com> - 1.0.0-1
- Initial package
