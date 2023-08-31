import paramiko
import re
import ipaddress
import multiprocessing
import os
from configparser import ConfigParser

def read_configuration():
    config = ConfigParser()
    config.read('config.ini')

    return {
        'SSH': {
            'username': config.get('SSH', 'username'),
            'password': config.get('SSH', 'password'),
            'ip_list_file': config.get('SSH', 'ip_list_file'),
            'log_dir': config.get('SSH', 'log_dir'),
        },
        'Firmware': {
            'bios_firmware_name': config.get('Firmware', 'bios_firmware_name'),
            'idrac_firmware_name': config.get('Firmware', 'idrac_firmware_name'),
        },
        'FTP': {
            'ftp_server': config.get('FTP', 'ftp_server'),
            'ftp_path': config.get('FTP', 'ftp_path'),
            'ftp_user': config.get('FTP', 'ftp_user'),
            'ftp_password': config.get('FTP', 'ftp_password'),
        },
        'Version': {
            'iDRACVersion': config.get('Version', 'iDRACVersion'),
            'BiosVersion': config.get('Version', 'BiosVersion'),
        }
    }



commands = [
    "racadm set BIOS.MemSettings.CECriticalSEL Enabled",
    "racadm set bios.SysProfileSettings.SysProfile Custom",
    "racadm set BIOS.SysProfileSettings.CpuInterconnectBusLinkPower Disabled",
    "racadm set BIOS.SysProfileSettings.EnergyPerformanceBias MaxPower",
    "racadm set BIOS.SysProfileSettings.PcieAspmL1 Disabled",
    "racadm set BIOS.SysProfileSettings.ProcC1E Disabled",
    "racadm set BIOS.SysProfileSettings.ProcCStates Disabled",
    "racadm set BIOS.SysProfileSettings.ProcPwrPerf MaxPerf",
    "racadm set BIOS.SysProfileSettings.UncoreFrequency MaxUFS",
    "racadm set bios.SysProfileSettings.SysProfile PerfOptimized",
    "racadm set BIOS.SysSecurity.Tpm2Hierarchy Disabled",
    "racadm set BIOS.SysSecurity.TpmSecurity Off",
    "racadm set biOS.integratedDevices.sriovGlobalEnable Enabled",
    "racadm jobqueue create BIOS.Setup.1-1 -r Graceful"
    "racadm set NIC.DeviceLevelConfig.5.VirtualizationMode SRIOV"
    "racadm set NIC.DeviceLevelConfig.7.VirtualizationMode SRIOV"
    "racadm jobqueue create NIC.Slot.1-1-1"
    "racadm jobqueue create NIC.Slot.2-1-1"
    "racadm serveraction powercycle"  
]


def execute_ssh_command(ssh, command):
    stdin, stdout, stderr = ssh.exec_command(command)
    result = stdout.read().decode('utf-8')
    return result

def process_node(ip):
    ssh_config = config['SSH']
    username = ssh_config['username']
    password = ssh_config['password']
    log_dir = ssh_config['log_dir']
    try:
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(ip, username=config["username"], password=["password"], timeout=10)

        log_entry = f"\nNode: {ip}\n"
        print(f"Node: {ip}")

        for command in commands:
            result = execute_ssh_command(ssh, command)
            log_entry += f"Command: {command}\n"
            log_entry += f"Result:\n{result}\n"
            log_entry += "---------------------------\n"
            print(f"Command: {command}")
            print(f"Result:\n{result}")
            print("---------------------------")

        ssh.close()

        with open(os.path.join(log_dir, f"{ip}.log"), 'a') as log:
            log.write(log_entry)
    except Exception as e:
        error_msg = f"Error connecting to {ip}: {str(e)}"
        print(error_msg)
        with open(os.path.join(log_dir, f"{ip}_error.log"), 'a') as log:
            log.write(error_msg + "\n\n")

def main():
    ip_list = []
    with open(config['SSH']['ip_list_file'], 'r') as file:
        for line in file:
            ip = line.strip()
            if re.match(r'^\d+\.\d+\.\d+\.\d+$', ip):
                ipaddress.ip_address(ip)  
                ip_list.append(ip)

    os.makedirs(log_dir, exist_ok=True)

    num_processes = min(len(ip_list), multiprocessing.cpu_count())  

    with multiprocessing.Pool(processes=num_processes) as pool:
        pool.map(process_node, ip_list)

if __name__ == "__main__":
    main()
