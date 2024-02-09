import subprocess
import os
from gui.src.services import goldwarden

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir))

def get_pin(message):
    p = subprocess.Popen(["python3", "-m", "src.ui.pinentry"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=root_path, start_new_session=True)
    p.stdin.write(f"{message}\n".encode())
    p.stdin.flush()
    pin = p.stdout.readline().decode().strip()
    if pin == "":
        return None
    return pin

def get_approval(message):
    p = subprocess.Popen(["python3", "-m", "src.ui.pinentry_approval"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
    p.stdin.write(f"{message}\n".encode())
    p.stdin.flush()
    result = p.stdout.readline().decode().strip()
    if result == "true":
        return True
    return False

def daemon():
    goldwarden.listen_for_pinentry(get_pin, get_approval)

daemon()