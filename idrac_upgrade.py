import paramiko
import re
import ipaddress
import multiprocessing
import os
import time
import subprocess
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

def execute_ssh_command(ssh, command):
    stdin, stdout, stderr = ssh.exec_command(command)
    result = stdout.read().decode('utf-8')
    return result

def is_pingable(ip):
    try:
        subprocess.run(["ping", "-c", "1", ip], stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=True)
        return True
    except subprocess.CalledProcessError:
        return False

def process_node(ip, config):
    ssh_config = config['SSH']
    username = ssh_config['username']
    password = ssh_config['password']
    log_dir = ssh_config['log_dir']

    firmware_config = config['Firmware']
    bios_firmware_name = firmware_config['bios_firmware_name']

    ftp_config = config['FTP']
    ftp_server = ftp_config['ftp_server']
    ftp_path = ftp_config['ftp_path']
    ftp_user = ftp_config['ftp_user']
    ftp_password = ftp_config['ftp_password']

    version_config = config['Version']
    iDRACVersion = version_config['iDRACVersion']
    BiosVersion = version_config['BiosVersion']

    update_command_template = (
        f'racadm update -f {bios_firmware_name} -u {ftp_user} -p {ftp_password} -l ftp://{ftp_server}{ftp_path}'
    )
    
    job_queue_command = "racadm jobqueue view"
    
    reboot_command = "racadm serveraction hardreset"

    if not is_pingable(ip):
        print(f"Node {ip} is not reachable, skipping...")
        return

    try:
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(ip, username=username, password=password, timeout=10)

        log_entry = f"\nNode: {ip}\n"
        print(f"Node: {ip}")

        version_command = "racadm getversion"
        version_result = execute_ssh_command(ssh, version_command)
        log_entry += f"Node: {ip}\n"
        log_entry += f"Command: {version_command}\n"
        log_entry += f"Result:\n{version_result}\n"
        log_entry += "---------------------------\n"
        print(f"Node: {ip}")
        print(f"Command: {version_command}")
        print(f"Result:\n{version_result}")
        print("---------------------------")

        bios_version = ""
        version_lines = version_result.split("\n")
        for line in version_lines:
            if "BIOS" in line:
                bios_version = line.split(":")[1].strip()

        if bios_version != BiosVersion:
            print(f"BIOS version does not match for Node: {ip}")
            
            # Execute update command
            update_command = update_command_template
            result = execute_ssh_command(ssh, update_command)
            log_entry += f"Node: {ip}\n"
            log_entry += f"Command: {update_command}\n"
            log_entry += f"Result:\n{result}\n"
            log_entry += "---------------------------\n"
            print(f"Node: {ip}")
            print(f"Command: {update_command}")
            print(f"Result:\n{result}")
            print("---------------------------")
            
            # Rest of the code for rebooting, job queue view, writing logs, etc.

        else:
            print(f"BIOS version matches for Node: {ip}, skipping update.")

    except Exception as e:
        error_msg = f"Error connecting to {ip}: {str(e)}"
        print(error_msg)
        log_path = os.path.join(log_dir, f"{ip}_error.log")
        with open(log_path, 'a') as log:
            log.write(error_msg + "\n\n")

def main():
    config = read_configuration()

    ip_list = []
    with open(config['SSH']['ip_list_file'], 'r') as file:
        for line in file:
            ip = line.strip()
            if re.match(r'^\d+\.\d+\.\d+\.\d+$', ip):
                ipaddress.ip_address(ip)
                ip_list.append(ip)

    os.makedirs(config['SSH']['log_dir'], exist_ok=True)

    num_processes = min(len(ip_list), multiprocessing.cpu_count())

    with multiprocessing.Pool(processes=num_processes) as pool:
        pool.starmap(process_node, [(ip, config) for ip in ip_list])

if __name__ == "__main__":
    main()