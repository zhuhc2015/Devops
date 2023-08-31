```
# iDRAC and BIOS Upgrade Scripts

This repository contains two scripts for upgrading iDRAC (Integrated Dell Remote Access Controller) and BIOS firmware versions of remote servers using SSH and FTP protocols.

## Installation and Setup

1. Install the required Python dependencies:
   pip install paramiko
```

This will install the necessary `paramiko` module for SSH connections.

1. Run the `install_ftp.sh` script to set up the FTP server on your system.

   ```
   ./install_ftp.sh
   ```

   This script will install the necessary FTP server packages and set up the required configuration.

2. Modify the `config.ini` file to match your environment.

   Open the `config.ini` file and make the following modifications:

   ### SSH Configuration

   In the `[SSH]` section, configure the SSH connection settings:

   ```
   iniCopy code[SSH]
   username = root
   password = calvin
   ip_list_file = ip_list.txt
   log_dir = logs
   ```

   ### Firmware Configuration

   In the `[Firmware]` section, configure firmware file names and FTP settings:

   ```
   iniCopy code[Firmware]
   bios_firmware_name = BIOS_NVD2K_WN64_2.17.1.EXE
   idrac_firmware_name = iDRAC-with-Lifecycle-Controller_Firmware_D92HF_WN64_6.00.30.00_A00.EXE
   ```

   ### FTP Configuration

   In the `[FTP]` section, configure FTP settings for uploading firmware files:

   ```
   iniCopy code[FTP]
   ftp_server = 10.137.2.45
   ftp_path = /
   ftp_user = admin
   ftp_password = admin
   ```

   ### Version Configuration

   In the `[Version]` section, set the current iDRAC and BIOS versions:

   ```
   iniCopy code[Version]
   iDRACVersion = 6.00.30.00
   BiosVersion = 2.17.1
   ```

3. Place the BIOS and iDRAC firmware packages in the appropriate directory.

   The firmware files should be placed in the directory specified in the `[Firmware]` section of the `config.ini` file.

## Running the Upgrade Scripts

### idrac_upgrade_bios.py

This script is used to upgrade the BIOS firmware of remote servers.

1. Ensure you have completed the installation and setup steps mentioned above.

2. Create a list of server IP addresses in the `ip_list.txt` file.

3. Run the script using the following command:

   ```
   python3 idrac_upgrade_bios.py
   ```

### idrac_upgrade.py

This script is used to upgrade the iDRAC firmware of remote servers.

1. Ensure you have completed the installation and setup steps mentioned above.

2. Create a list of server IP addresses in the `ip_list.txt` file.

3. Run the script using the following command:

   ```
   python3 idrac_upgrade.py
   ```

## Checking Firmware Versions

The `idrac_bios_check_version.py` script can be used to check the current versions of iDRAC and BIOS firmware on remote servers.

1. Ensure you have completed the installation and setup steps mentioned above.

2. Create a list of server IP addresses in the `ip_list.txt` file.

3. Run the script using the following command:

   ```
   python3 idrac_bios_check_version.py
   ```

