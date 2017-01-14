# The name of your package
Name: backup_cydex

# A short summary of your package
Summary: None

# The version of your package
Version: 3.0

# The release number of your package
Release: 1

# Any license you wish to list
License: GNU GPL

# What group this RPM would typically reside in
Group: Applications/System

# Who packaged this RPM
Packager: YourName <your_email@address.com>

# The build architecture of this RPM (noarch/x86_64/i386/etc)
Buildarch: noarch
#Warn this is specially  add,qhen you meet the situation that noarch problem
%define _binaries_in_noarch_packages_terminate_build   0
# You generally should not need to mess with this setting
Buildroot: %{_tmppath}/%{name}

# Change this extension to change the compression level in your RPM
#  tar / tar.gz / tar.bz2
Source0: %{name}.tar

# If you are having trouble building a package and need to disable
#  automatic dependency/provides checking, uncomment this:
# AutoReqProv: no

# If this package has prerequisites, uncomment this line and
#  list them here - examples are already listed
#Requires: bash, python >= 2.7

# A more verbose description of your package
%description
None

# You probably do not need to change this
%define debug_package %{nil}


%prep
%setup -q -c

%build

%install
rsync -a . %{buildroot}/

%clean
rm -rf %{buildroot}

%pre

%post

%preun

%postun

#%trigger

#%triggerin

#%triggerun

%changelog
* Sat Dec 31 2016 YourName <your_email@address.com>
- Initial version.

%files
%attr(0755, root, root) "/opt/cydex/ftpupload.py"
%attr(0644, root, root) "/opt/cydex/ftpupload.pyc"
%attr(0644, root, root) "/opt/cydex/ftpupload.pyo"
%attr(0755, root, root) "/opt/cydex/ftpuploadincre.py"
%attr(0644, root, root) "/opt/cydex/ftpuploadincre.pyc"
%attr(0644, root, root) "/opt/cydex/ftpuploadincre.pyo"
%attr(0755, root, root) "/opt/cydex/fullbackup.sh"
%attr(0755, root, root) "/opt/cydex/increbackup.sh"
%attr(0644, root, root) "/opt/cydex/libauth.so"
%attr(0755, root, root) "/opt/cydex/start.sh"
%attr(0644, root, root) "/opt/cydex/upload.py"
%attr(0644, root, root) "/opt/cydex/upload.pyc"
%attr(0644, root, root) "/opt/cydex/upload.pyo"