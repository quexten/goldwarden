#!/usr/bin/python
import time
import subprocess
from tendo import singleton
import monitors.dbus_autofill_monitor
import sys
import goldwarden
from threading import Thread

try:
    subprocess.Popen(["python3", "/app/bin/background.py"], start_new_session=True)
except:
    pass

is_hidden = "--hidden" in sys.argv

if not is_hidden:
    try:
        subprocess.Popen(["python3", "/app/bin/settings.py"], start_new_session=True)
    except:
        subprocess.Popen(["python3", "./settings.py"], start_new_session=True)
        pass

try:
    me = singleton.SingleInstance()
except:
    exit()

def run_daemon():
    # todo: do a proper check
    if is_hidden:
        time.sleep(20)
    if not goldwarden.is_daemon_running():
        goldwarden.run_daemon()

def on_autofill():
    subprocess.Popen(["python3", "/app/bin/autofill.py"], start_new_session=True)

monitors.dbus_autofill_monitor.on_autofill = lambda: on_autofill()
monitors.dbus_autofill_monitor.run_daemon()

while True:
    time.sleep(60)