%if 0%{?rhel} == 7
  %define dist .el7
%endif
%define _unpackaged_files_terminate_build 0
Name: rweb
Version: 0.1
Release:	1%{?dist}
Summary: A golang gin tonic front end for Relevy.

License: GPLv2
URL: https://github.com/Jmainguy/relevy-web
Source0: relevy-web.tar.gz

%description
A golang gin tonic front end for Relevy.

%prep
%setup -q -n relevy-web
%install
mkdir -p $RPM_BUILD_ROOT/usr/sbin
mkdir -p $RPM_BUILD_ROOT/opt/rweb
mkdir -p $RPM_BUILD_ROOT/usr/lib/systemd/system
mkdir -p $RPM_BUILD_ROOT/etc/relevy
install -m 0755 $RPM_BUILD_DIR/relevy-web/rweb %{buildroot}/usr/sbin
install -m 0644 $RPM_BUILD_DIR/relevy-web/service/rweb.service %{buildroot}/usr/lib/systemd/system
install -m 0644 $RPM_BUILD_DIR/relevy-web/template.html %{buildroot}/opt/rweb
install -m 0644 $RPM_BUILD_DIR/relevy-web/sorttable.js %{buildroot}/opt/rweb
install -m 0644 $RPM_BUILD_DIR/relevy-web/config.yaml %{buildroot}/etc/relevy/

%files
/usr/sbin/rweb
/usr/lib/systemd/system/rweb.service
%dir /opt/rweb
%dir /etc/relevy
/opt/rweb/template.html
/opt/rweb/sorttable.js
%config(noreplace) /etc/relevy/config.yaml

%pre
getent group rweb >/dev/null || groupadd -r rweb
getent passwd rweb >/dev/null || \
    useradd -r -g rweb -d /opt/rweb -s /sbin/nologin \
    -c "User to run rweb service" rweb
exit 0

%post
chown -R rweb:rweb /opt/rweb
if [ -f /usr/bin/systemctl ]; then
  systemctl daemon-reload
fi

%changelog

