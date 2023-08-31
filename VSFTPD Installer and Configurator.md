```
# VSFTPD Installer and Configurator

This Python script helps you install and configure the VSFTPD FTP server on different Linux distributions. It supports both offline installation from a package file and online installation from repositories. The script reads configuration details from the `package.cfg` file.

## Prerequisites

- Python 3.x
- Administrative privileges (`sudo`)

## Getting Started

1. Clone or download this repository.

2. Place your VSFTPD offline installation package in the same directory as the script. The package can be in `.rpm` or `.deb` format.

3. Create a `package.cfg` file in the same directory with the following content:

   ```ini
   [vsftpd]
   offline_package_path = /path/to/vsftpd_package
```

Replace `/path/to/vsftpd_package` with the actual path to your offline package.

1. Open a terminal and navigate to the directory containing the script and the `package.cfg` file.

2. Make the script executable:

   ```
   chmod +x vsftpd_installer.py
   ```

3. Run the script:

   ```
   ./vsftpd_installer.py
   ```

   The script will install and configure VSFTPD on your system.

## Configuration

- Edit the `package.cfg` file to change the offline package path if needed.

## Supported Distributions

- Ubuntu
- Debian
- CentOS
- openSUSE

## Notes

- Ensure that you have the necessary administrative privileges to install software.
- Use this script at your own risk. It modifies system configurations and requires careful use.
- For online installation, ensure your system has access to the internet and appropriate package repositories.
- For offline installation, make sure to provide the correct package file path.
- If you encounter issues, check the terminal output for error messages.