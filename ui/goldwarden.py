import subprocess
import json
import os

# if flatpak
if os.path.exists("/app/bin/goldwarden"):
    BINARY_PATH = "/app/bin/goldwarden"
else:
    res = subprocess.run(["which", "goldwarden"])
    BINARY_PATH = res.stdout.decode("utf-8").strip()

def set_api_url(url):
    restic_cmd = f"{BINARY_PATH} config set-api-url {url}"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)

def set_identity_url(url):
    restic_cmd = f"{BINARY_PATH} config set-identity-url {url}"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    
def set_notification_url(url):
    restic_cmd = f"{BINARY_PATH} config set-notifications-url {url}"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)

def login_with_password(email, password):
    restic_cmd = f"{BINARY_PATH} vault login --email {email}"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
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
    restic_cmd = f"{BINARY_PATH} vault login --email {email} --passwordless"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    
def is_pin_enabled():
    restic_cmd = f"{BINARY_PATH} vault pin status"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    # check if contains enabled
    return "enabled" in result.stderr

def enable_pin():
    restic_cmd = f"{BINARY_PATH} vault pin set"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    
def unlock():
    restic_cmd = f"{BINARY_PATH} vault unlock"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    
def lock():
    restic_cmd = f"{BINARY_PATH} vault lock"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    
def purge():
    restic_cmd = f"{BINARY_PATH} vault purge"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)

def get_vault_status():
    restic_cmd = f"{BINARY_PATH} vault status"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    try:
        return json.loads(result.stdout)
    except Exception as e:
        print(e)
        return None
    
def get_vault_logins():
    restic_cmd = f"{BINARY_PATH} logins list"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)
    try:
        return json.loads(result.stdout)
    except Exception as e:
        print(e)
        return None
    
def autotype(username, password):
    # environment
    env = os.environ.copy()
    env["PASSWORD"] = password
    restic_cmd = f"{BINARY_PATH} autotype --username {username}"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True, env=env)
    if result.returncode != 0:
        raise Exception("Failed to initialize repository, err", result.stderr)

def is_daemon_running():
    restic_cmd = f"{BINARY_PATH} vault status"
    result = subprocess.run(restic_cmd.split(), capture_output=True, text=True)
    if result.returncode != 0:
        return False
    daemon_not_running = ("daemon running?" in result.stderr or "daemon running" in result.stderr)
    return not daemon_not_running

def run_daemon():
    restic_cmd = f"{BINARY_PATH} daemonize"
    # print while running
    result = subprocess.Popen(restic_cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    if result.returncode != 0:
        print("Failed err", result.stderr)
    for line in result.stdout:
        print(line.decode("utf-8"))
    result.wait()
    print("quitting goldwarden daemon")
    return result.returncode