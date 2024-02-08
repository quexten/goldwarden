#!/usr/bin/python
import time
import subprocess
from tendo import singleton
from .monitors import dbus_autofill_monitor
from .monitors import dbus_monitor
import sys
from services import goldwarden
from threading import Thread
import os
import secrets
import time
import os

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir))

def main():
    token = secrets.token_hex(32)
    print("token", token)
    if not os.environ.get("GOLDWARDEN_DAEMON_AUTH_TOKEN") == None:
        token = os.environ["GOLDWARDEN_DAEMON_AUTH_TOKEN"]

    # check if already running
    try:
        me = singleton.SingleInstance()
    except:
        import dbus

        bus = dbus.SessionBus()
        the_object = bus.get_objeect("com.quexten.Goldwarden.dbus", "/com/quexten/Goldwarden")
        the_interface = dbus.Interface(the_object, "com.quexten.Goldwarden.Settings")
        reply = the_interface.settings()
        exit()

    # start daemons
    dbus_autofill_monitor.run_daemon() # todo: remove after migration
    dbus_monitor.run_daemon(token)

    if not "--hidden" in sys.argv:
        subprocess.Popen(["python3", "-m", "src.ui.settings"], cwd=root_path, start_new_session=True)

    # try:
    #     subprocess.Popen(["python3", f'{source_path}/background.py'], start_new_session=True)
    # except Exception as e:
    #     pass

    while True:
        time.sleep(60)


# def run_daemon():
#     # todo: do a proper check
#     if is_hidden:
#         time.sleep(20)
#     print("IS daemon running", goldwarden.is_daemon_running())
#     if not goldwarden.is_daemon_running():
#         print("running daemon")
#         goldwarden.run_daemon()
#         print("daemon running")
