import paramiko
import ipaddress
import multiprocessing
import os
import subprocess
from configparser import ConfigParser

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
    iDRACVersion = config['Version']['iDRACVersion']
    BiosVersion = config['Version']['BiosVersion']
    username = config['SSH']['username']
    password = config['SSH']['password']
    ip_list_file = config['Paths']['ip_list_file']
    log_dir = config['Paths']['log_dir']

    if not is_pingable(ip):
        print(f"Node {ip} is not pingable, skipping...")
        return

    print(f"Ping successful for Node: {ip}")

    try:
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(ip, username=username, password=password, timeout=10)

        log_entry = f"\nNode: {ip}\n"
        print(f"Node: {ip}")

        commands = ["racadm getversion"]  # Move commands here if they're specific to each node
        for command in commands:
            result = execute_ssh_command(ssh, command)
            log_entry += f"Node: {ip}\n"
            log_entry += f"Command: {command}\n"
            log_entry += f"Result:\n{result}\n"
            log_entry += "---------------------------\n"
            print(f"Node: {ip}")
            print(f"Command: {command}")
            print(f"Result:\n{result}")
            print("---------------------------")

            if result.strip() != iDRACVersion and result.strip() != BiosVersion:
                version_mismatch_log_path = os.path.join(log_dir, "version_mismatch.log")
                with open(version_mismatch_log_path, 'a') as log:
                    log.write(f"IP: {ip}\n")

        ssh.close()

        if result.strip() == iDRACVersion or result.strip() == BiosVersion:
            print(f"Versions match for Node: {ip}")
        else:
            log_path = os.path.join(log_dir, f"{ip}.log")
            with open(log_path, 'a') as log:
                log.write(log_entry)

    except Exception as e:
        error_msg = f"Error connecting to {ip}: {str(e)}"
        print(error_msg)
        log_path = os.path.join(log_dir, f"{ip}_error.log")
        with open(log_path, 'a') as log:
            log.write(error_msg + "\n\n")

def main():
    config = ConfigParser()
    config.read('config.ini')

    ip_list = []
    with open(config['Paths']['ip_list_file'], 'r') as file:
        for line in file:
            ip = line.strip()
            if re.match(r'^\d+\.\d+\.\d+\.\d+$', ip):
                ipaddress.ip_address(ip)
                ip_list.append(ip)

    os.makedirs(config['Paths']['log_dir'], exist_ok=True)

    num_processes = min(len(ip_list), multiprocessing.cpu_count())

    with multiprocessing.Pool(processes=num_processes) as pool:
        pool.map(process_node, ip_list, [config]*len(ip_list))

if __name__ == "__main__":
    main()
