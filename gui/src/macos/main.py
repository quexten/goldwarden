import sys
import subprocess
import os

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir))

def main():
    token = "abc"
    print("Starting Goldwarden GUI")
    if not "--hidden" in sys.argv:
        p = subprocess.Popen(["python3", "-m", "src.gui.settings"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
        p.stdin.write(f"{token}\n".encode())
        p.stdin.flush()