#!/usr/bin/python
import time
import subprocess
from tendo import singleton
import monitors.dbus_autofill_monitor
import sys

try:
    subprocess.Popen(["python3", "/app/bin/background.py"], start_new_session=True)
except:
    pass

if "--hidden" not in sys.argv:
    try:
        subprocess.Popen(["python3", "/app/bin/settings.py"], start_new_session=True)
    except:
        pass

try:
    me = singleton.SingleInstance()
except:
    exit()

def on_autofill():
    subprocess.Popen(["python3", "/app/bin/autofill.py"], start_new_session=True)
monitors.dbus_autofill_monitor.on_autofill = on_autofill

while True:
    time.sleep(60)