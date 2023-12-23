import subprocess

def write(text):
    process = subprocess.Popen(["/bin/sh", "-c", "wl-copy"], stdin=subprocess.PIPE)
    process.communicate(text.encode('utf-8'))