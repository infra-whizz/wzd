# spec file for whizz-client package

Name:           whizz-client
Version:        0.9
Release:        0
Summary:        Whizz client daemon
License:        MIT
Group:          System/Tools
Url:            https://gitlab.com/infra-whizz/wzd
Source:         %{name}-%{version}.tar.gz
Source1:        vendor.tar.gz

BuildRequires:  golang-packaging
BuildRequires:  golang(API) >= 1.13

%description
The client component of Whizz configuration management system

%prep
%setup -q
%setup -q -T -D -a 1

%build
go build -x -mod=vendor -buildmode=pie -o wzd ./cmd/*.go

%install
install -D -m 0755 wzd %{buildroot}%{_bindir}/wzd
mkdir -p %{buildroot}%{_sysconfdir}
install -m 0644 ./etc/wzd.conf.example %{buildroot}%{_sysconfdir}/wzd.conf

%files
%defattr(-,root,root)
%{_bindir}/wzd
%dir %{_sysconfdir}
%config /etc/wzd.conf

%changelog
