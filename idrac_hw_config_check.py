import paramiko
import re
import ipaddress
import multiprocessing
import os
from configparser import ConfigParser

def execute_ssh_command(ssh, command):
    stdin, stdout, stderr = ssh.exec_command(command)
    result = stdout.read().decode('utf-8')
    return result

def process_node(ip, config):
    try:
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(ip, username=config["username"], password=config["password"], timeout=10)

        log_entry = f"\nNode: {ip}\n"
        print(f"Node: {ip}")

        for command in config["commands"]:
            result = execute_ssh_command(ssh, command)
            log_entry += f"Node: {ip}\n"
            log_entry += f"Command: {command}\n"
            log_entry += f"Result:\n{result}\n"
            log_entry += "---------------------------\n"
            print(f"Node: {ip}")
            print(f"Command: {command}")
            print(f"Result:\n{result}")
            print("---------------------------")

        ssh.close()

        log_path = os.path.join(config["log_dir"], f"{ip}.log")
        with open(log_path, 'a') as log:
            log.write(log_entry)
    except Exception as e:
        error_msg = f"Error connecting to {ip}: {str(e)}"
        print(error_msg)
        log_path = os.path.join(config["log_dir"], f"{ip}_error.log")
        with open(log_path, 'a') as log:
            log.write(error_msg + "\n\n")

def main():
    config = ConfigParser()
    config.read('config.ini')

    ip_list = []
    with open(config["SSH"]["ip_list_file"], 'r') as file:
        for line in file:
            ip = line.strip()
            if re.match(r'^\d+\.\d+\.\d+\.\d+$', ip):
                ipaddress.ip_address(ip)
                ip_list.append(ip)

    os.makedirs(config["SSH"]["log_dir"], exist_ok=True)

    num_processes = min(len(ip_list), multiprocessing.cpu_count())

    with multiprocessing.Pool(processes=num_processes) as pool:
        pool.map(process_node, ip_list, [config["SSH"]] * len(ip_list))

if __name__ == "__main__":
    main()
