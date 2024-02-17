import subprocess
import os
from src.services import goldwarden
from threading import Thread
import time

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir))

def get_pin(message):
    p = subprocess.Popen(["python3", "-m", "src.gui.pinentry"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=root_path, start_new_session=True)
    p.stdin.write(f"{message}\n".encode())
    p.stdin.flush()
    pin = p.stdout.readline().decode().strip()
    if pin == "":
        return None
    return pin

def get_approval(message):
    p = subprocess.Popen(["python3", "-m", "src.gui.pinentry_approval"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
    p.stdin.write(f"{message}\n".encode())
    p.stdin.flush()
    result = p.stdout.readline().decode().strip()
    if result == "true":
        return True
    return False

def daemon():
    goldwarden.listen_for_pinentry(get_pin, get_approval)

def daemonize():
    #todo fix this
    time.sleep(3)
    thread = Thread(target=daemon)
    thread.start()