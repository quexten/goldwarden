import subprocess
import json
import os
from pathlib import Path
from threading import Thread
from shutil import which
import sys

is_flatpak = os.path.exists("/.flatpak-info")
log_directory = str(Path.home()) + "/.local/share/goldwarden"
if is_flatpak:
    log_directory = str(Path.home()) + "/.var/app/com.quexten.goldwarden/data/goldwarden"
os.makedirs(log_directory, exist_ok=True)

# detect goldwarden binary
BINARY_PATHS = [
    "/app/bin/goldwarden",
    "/usr/bin/goldwarden",
    str(Path.home()) + "/go/src/github.com/quexten/goldwarden/goldwarden"
]

BINARY_PATH = None
for path in BINARY_PATHS:
    if os.path.exists(path):
        BINARY_PATH = path
        break

if BINARY_PATH is None:
    BINARY_PATH = which('goldwarden')
    if isinstance(BINARY_PATH,str):
        BINARY_PATH = BINARY_PATH.strip()

if BINARY_PATH is None:
    print("goldwarden executable not found")
    sys.exit()

authenticated_connection = None

def create_authenticated_connection(token):
    print("create authenticated connection")
    global authenticated_connection
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

    authenticated_connection.stdin.write(cmd + "\n")
    authenticated_connection.stdin.flush()
    result = authenticated_connection.stdout.readline()
    return result

def set_api_url(url):
    send_authenticated_command(f"config set-api-url {url}")
    
def set_identity_url(url):
    send_authenticated_command(f"config set-identity-url {url}")

def set_notification_url(url):
    send_authenticated_command(f"config set-notifications-url {url}")

def set_vault_url(url):
    send_authenticated_command(f"config set-vault-url {url}")

def set_server(url):
    result = send_authenticated_command(f"config set-server {url}")
    if result.strip() != "Done":
        raise Exception("Failed to set server")

def get_environment():
    result = send_authenticated_command(f"config get-environment")
    try:
        return json.loads(result)
    except Exception as e:
        return None

def set_client_id(client_id):
    send_authenticated_command(f"config set-client-id \"{client_id}\"")

def set_client_secret(client_secret):
    send_authenticated_command(f"config set-client-secret \"{client_secret}\"")

def login_with_password(email, password):
    result = send_authenticated_command(f"vault login --email {email}")
    if "Login failed" in result and "username or password" in result.lower():
        raise Exception("errorbadpassword")
    if "Login failed" in result and ("error code 7" in result.lower() or "error code 6" in result.lower()):
        raise Exception("errorcaptcha")
    if "Login failed" in result and "two-factor" in result.lower():
        raise Exception("errortotp")

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
        return None
    
def get_vault_logins():
    result = send_authenticated_command(f"logins list")
    try:
        return json.loads(result)
    except Exception as e:
        return None

def get_runtime_config():
    result = send_authenticated_command(f"config get-runtime-config")
    try:
        return json.loads(result)
    except Exception as e:
        return None
    
def autotype(text):
    goldwarden_cmd = f"{BINARY_PATH} autotype"
    process = subprocess.Popen(goldwarden_cmd.split(), stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    text_hex = text.encode("utf-8").hex()
    process.stdin.write(text_hex + "\n")
    process.stdin.flush()
    process.wait()

def version():
    result = send_authenticated_command(f"version")
    return result.strip()

def is_daemon_running():
    result = send_authenticated_command(f"vault status")
    daemon_not_running = ("daemon running" in result)
    return not daemon_not_running

def listen_for_pinentry(on_pinentry, on_pin_approval):
    print("listening for pinentry", BINARY_PATH)
    pinentry_process = subprocess.Popen([f"{BINARY_PATH}", "pinentry"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    while True:
        line = pinentry_process.stdout.readline()
        # starts with pin-request
        if line.startswith("pin-request"):
            text = line.split(",")[1].strip()
            pin = on_pinentry(text)
            if pin == None:
                pin = ""
            pinentry_process.stdin.write(pin + "\n")
            pinentry_process.stdin.flush()
        if line.startswith("approval-request"):
            text = line.split(",")[1].strip()
            approval = on_pin_approval(text)
            if approval:
                pinentry_process.stdin.write("true\n")
                pinentry_process.stdin.flush()
            else:
                pinentry_process.stdin.write("false\n")
                pinentry_process.stdin.flush()

def run_daemon(token):
    #todo replace with stdin
    daemon_env = os.environ.copy()
    daemon_env["GOLDWARDEN_DAEMON_AUTH_TOKEN"] = token
    
    print("starting goldwarden daemon", BINARY_PATH)

    # print while running
    result = subprocess.Popen([f"{BINARY_PATH}", "daemonize"], stdout=subprocess.PIPE, text=True, env=daemon_env)

    # write log to file until process exits
    log_file = open(f"{log_directory}/daemon.log", "w")
    while result.poll() == None:
        # read stdout and stder
        stdout = result.stdout.readline()
        log_file.write(stdout)
        log_file.flush()
    log_file.close()
    print("quitting goldwarden daemon")
    return result.returncode

def run_daemon_background(token):
    thread = Thread(target=lambda: run_daemon(token))
    thread.start()