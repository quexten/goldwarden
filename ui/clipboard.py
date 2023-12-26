import subprocess
import os

def write(text):
    # set path
    env = os.environ.copy()
    env["PATH"] = env["PATH"] + ":/app/bin"
    process = subprocess.Popen(["/bin/sh", "-c", "wl-copy"], stdin=subprocess.PIPE, env=env)
    process.communicate(text.encode('utf-8'))