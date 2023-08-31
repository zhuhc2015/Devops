#!/usr/bin/env python3

import os
import subprocess
from configparser import ConfigParser

# Read the package.cfg configuration file
config = ConfigParser()
config.read("package.cfg")
vsftpd_package_path = config.get("vsftpd", "offline_package_path")

# Function to install vsftpd from an offline package on SUSE
def install_vsftpd_offline_suse():
    subprocess.run(["sudo", "zypper", "install", vsftpd_package_path, "-y"])

# Function to install vsftpd from an offline package on CentOS
def install_vsftpd_offline_centos():
    subprocess.run(["sudo", "yum", "install", vsftpd_package_path, "-y"])

# Function to install vsftpd from an offline package on Ubuntu/Debian
def install_vsftpd_offline_debian():
    subprocess.run(["sudo", "dpkg", "-i", f"{vsftpd_package_path}.deb"])
    subprocess.run(["sudo", "apt-get", "install", "-f", "-y"])

# Function to install vsftpd online on Ubuntu/Debian
def install_vsftpd_online_debian():
    subprocess.run(["sudo", "apt-get", "update"])
    subprocess.run(["sudo", "apt-get", "install", "vsftpd", "-y"])

# Function to install vsftpd online on SUSE
def install_vsftpd_online_suse():
    subprocess.run(["sudo", "zypper", "install", "vsftpd", "-y"])

# Function to install vsftpd online on CentOS
def install_vsftpd_online_centos():
    subprocess.run(["sudo", "yum", "install", "vsftpd", "-y"])

# Function to configure vsftpd
def configure_vsftpd():
    with open("/etc/vsftpd.conf", "r") as f:
        config = f.read()
    config = config.replace("anonymous_enable=YES", "anonymous_enable=NO")
    config += """
local_enable=YES
write_enable=YES
local_umask=022
chroot_local_user=YES
local_root=/home/admin/
"""
    with open("/etc/vsftpd.conf", "w") as f:
        f.write(config)

# Function to set permissions for the root directory
def set_root_directory_permissions():
    os.system("sudo chown root:root /home/admin/")
    os.system("sudo chmod 755 /home/admin")

# Function to create the admin user
def create_admin_user():
    os.system("sudo useradd -m admin")
    os.system("echo 'admin:admin' | sudo chpasswd")

# Function to set permissions for the admin user
def set_admin_user_permissions():
    os.system("sudo chown -R admin:admin /home/admin")
    os.system("sudo chmod -R 755 /home/admin")

# Function to restart vsftpd service
def restart_vsftpd():
    os.system("sudo service vsftpd restart")
    os.system("sudo service vsftpd status")

# Main script
if os.path.exists("/etc/os-release"):
    with open("/etc/os-release") as f:
        os_info = dict(line.strip().split('=') for line in f if '=' in line)
    os_id = os_info.get("ID")
    if os_id == "ubuntu" or os_id == "debian":
        subprocess.run(["sudo", "apt-get", "update"])
        install_vsftpd_online_debian()
    elif os_id == "suse":
        install_vsftpd_online_suse()
    elif os_id == "centos":
        install_vsftpd_online_centos()
    else:
        print("Unsupported operating system.")
        exit(1)

    configure_vsftpd()
    set_root_directory_permissions()
    create_admin_user()
    set_admin_user_permissions()
    restart_vsftpd()
    print("FTP server installed, and admin user created.")
else:
    print("Unable to determine the operating system.")
    exit(1)
