import sys
import subprocess
import os
import secrets
from src.services import goldwarden
from src.services import pinentry
import time

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir))

def main():
    token = secrets.token_hex(32)
    if not os.environ.get("GOLDWARDEN_DAEMON_AUTH_TOKEN") == None:
        token = os.environ["GOLDWARDEN_DAEMON_AUTH_TOKEN"]
    print("Starting Goldwarden GUI")
    goldwarden.run_daemon_background(token)
    time.sleep(1)
    #pinentry.daemonize()
    if not "--hidden" in sys.argv:
        p = subprocess.Popen(["python3", "-m", "src.gui.settings"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=root_path, start_new_session=True)
        p.stdin.write(f"{token}\n".encode())
        p.stdin.flush()
        # print stdout
        while True:
            line = p.stderr.readline()
            if not line:
                break
            print(line.decode().strip())
    while True:
        time.sleep(60)