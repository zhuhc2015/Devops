#!/bin/bash

# Path to the vsftpd offline installation package
vsftpd_package_path="/path/to/vsftpd_package"

# Function to install vsftpd from an offline package on SUSE
install_vsftpd_offline_suse() {
    sudo zypper install "$vsftpd_package_path" -y
}

# Function to install vsftpd from an offline package on CentOS
install_vsftpd_offline_centos() {
    sudo yum install "$vsftpd_package_path" -y
}

# Function to install vsftpd from an offline package on Ubuntu/Debian
install_vsftpd_offline_debian() {
    sudo dpkg -i "$vsftpd_package_path.deb"
    sudo apt-get install -f
}

# Function to configure vsftpd
configure_vsftpd() {
    sudo sed -i 's/anonymous_enable=YES/anonymous_enable=NO/' /etc/vsftpd.conf
    echo "local_enable=YES" | sudo tee -a /etc/vsftpd.conf
    echo "write_enable=YES" | sudo tee -a /etc/vsftpd.conf
    echo "local_umask=022" | sudo tee -a /etc/vsftpd.conf
    echo "chroot_local_user=YES" | sudo tee -a /etc/vsftpd.conf
    echo "local_root=/home/admin/" | sudo tee -a /etc/vsftpd.conf
}

# Function to set permissions for the root directory
set_root_directory_permissions() {
    sudo chown root:root /home/admin/
    sudo chmod 755 /home/admin
}

# Function to create the admin user
create_admin_user() {
    sudo useradd -m admin
    echo "admin:admin" | sudo chpasswd
}

# Function to set permissions for the admin user
set_admin_user_permissions() {
    sudo chown -R admin:admin /home/admin
    sudo chmod -R 755 /home/admin
}

# Function to restart vsftpd service
restart_vsftpd() {
    sudo service vsftpd restart
    sudo service vsftpd status
}

# Main script
if [ -f /etc/os-release ]; then
    source /etc/os-release
    if [ "$ID" == "ubuntu" ] || [ "$ID" == "debian" ]; then
        install_vsftpd_offline_debian
    elif [ "$ID" == "suse" ]; then
        install_vsftpd_offline_suse
    elif [ "$ID" == "centos" ]; then
        install_vsftpd_offline_centos
    else
        echo "Unsupported operating system."
        exit 1
    fi

    configure_vsftpd
    set_root_directory_permissions
    create_admin_user
    set_admin_user_permissions
    restart_vsftpd
    echo "FTP server installed, and admin user created."
else
    echo "Unable to determine the operating system."
    exit 1
fi
