import subprocess
import json
import os
from pathlib import Path

BINARY_PATHS = [
    "/app/bin/goldwarden",
    "/usr/bin/goldwarden",
    str(Path.home()) + "/go/src/github.com/quexten/goldwarden/goldwarden"
]

authenticated_connection = None

def create_authenticated_connection(token):
    global authenticated_connection
    BINARY_PATH = None
    for path in BINARY_PATHS:
        if os.path.exists(path):
            BINARY_PATH = path
            break
    if BINARY_PATH == None:
        raise Exception("Could not find goldwarden binary")
    authenticated_connection = subprocess.Popen([f"{BINARY_PATH}", "session"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if not token == None:
        authenticated_connection.stdin.write("authenticate-session " + token + "\n")
        authenticated_connection.stdin.flush()
        # read entire message
        result = authenticated_connection.stdout.readline()
        if "true" not in result:
            raise Exception("Failed to authenticate")

def send_authenticated_command(cmd):
    if authenticated_connection == None:
        print("No daemon connection running, please restart the application completely.")
        return ""

    print("sending command", cmd)
    authenticated_connection.stdin.write(cmd + "\n")
    authenticated_connection.stdin.flush()
    result = authenticated_connection.stdout.readline()
    print("result", result)
    return result

def set_api_url(url):
    send_authenticated_command(f"config set-api-url {url}")
    
def set_identity_url(url):
    send_authenticated_command(f"config set-identity-url {url}")

def set_notification_url(url):
    send_authenticated_command(f"config set-notifications-url {url}")

def set_vault_url(url):
    send_authenticated_command(f"config set-vault-url {url}")

def set_url(url):
    send_authenticated_command(f"config set-url {url}")

def get_environment():
    result = send_authenticated_command(f"config get-environment")
    try:
        return json.loads(result)
    except Exception as e:
        print(e)
        return None

def set_client_id(client_id):
    send_authenticated_command(f"config set-client-id \"{client_id}\"")

def set_client_secret(client_secret):
    send_authenticated_command(f"config set-client-secret \"{client_secret}\"")

def login_with_password(email, password):
    result = send_authenticated_command(f"vault login --email {email}")
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    if len(result.stderr.strip()) > 0:
        print(result.stderr)
        if "password" in result.stderr:
            return "badpass"
        else:
            if "Logged in" in result.stderr:
                print("ok")
                return "ok"
            return "error"
    print("ok")
    return "ok"

def login_passwordless(email):
    send_authenticated_command(f"vault login --email {email} --passwordless")
  
def is_pin_enabled():
    result = send_authenticated_command("vault pin status")
    return "enabled" in result

def enable_pin():
    send_authenticated_command(f"vault pin set")

def unlock():
    send_authenticated_command(f"vault unlock")
    
def lock():
    send_authenticated_command(f"vault lock")
    
def purge():
    send_authenticated_command(f"vault purge")

def get_vault_status():
    result = send_authenticated_command(f"vault status")
    try:
        return json.loads(result)
    except Exception as e:
        print(e)
        return None
    
def get_vault_logins():
    result = send_authenticated_command(f"logins list")
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    try:
        return json.loads(result)
    except Exception as e:
        return None

def get_runtime_config():
    result = send_authenticated_command(f"config get-runtime-config")
    print(result)
    try:
        return json.loads(result)
    except Exception as e:
        return None
    
# def autotype(username, password):
#     # environment
#     env = os.environ.copy()
#     env["PASSWORD"] = password
#     restic_cmd = f"{BINARY_PATH} autotype --username {username}"
#     result = subprocess.run(restic_cmd.split(), capture_output=True, text=True, env=env)
#     print(result.stderr)
#     print(result.stdout)
#     if result.returncode != 0:
#         raise Exception("Failed to initialize repository, err", result.stderr)

def is_daemon_running():
    result = send_authenticated_command(f"vault status")
    daemon_not_running = ("daemon running" in result)
    return not daemon_not_running

# def run_daemon():
#     restic_cmd = f"daemonize"
#     # print while running
#     result = subprocess.Popen(restic_cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
#     if result.returncode != 0:
#         print("Failed err", result.stderr)
#     for line in result.stdout:
#         print(line.decode("utf-8"))
#     result.wait()
#     print("quitting goldwarden daemon")
#     return result.returncode