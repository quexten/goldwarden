#!/usr/bin/env python3
import time
import subprocess
from tendo import singleton
from .monitors import dbus_autofill_monitor
from .monitors import dbus_monitor
from .monitors import locked_monitor
import sys
from src.services import goldwarden
from src.services import pinentry
from threading import Thread
import os
import secrets
import time
import os
import src.linux.flatpak.api as flatpak_api

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir))
is_flatpak = os.path.exists("/.flatpak-info")

def main():
    token = secrets.token_hex(32)
    if not os.environ.get("GOLDWARDEN_DAEMON_AUTH_TOKEN") == None:
        token = os.environ["GOLDWARDEN_DAEMON_AUTH_TOKEN"]

    # check if already running
    try:
        me = singleton.SingleInstance()
    except:
        import dbus

        bus = dbus.SessionBus()
        the_object = bus.get_object("com.quexten.Goldwarden.gui", "/com/quexten/Goldwarden/gui")
        the_interface = dbus.Interface(the_object, "com.quexten.Goldwarden.gui.actions")
        reply = the_interface.settings()
        exit()

    goldwarden.run_daemon_background(token)
    # start daemons
    dbus_autofill_monitor.run_daemon(token) # todo: remove after migration
    dbus_monitor.run_daemon(token)
    locked_monitor.run_daemon(token)
    pinentry.daemonize()

    if not "--hidden" in sys.argv:
        p = subprocess.Popen(["python3", "-m", "src.gui.settings"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
        p.stdin.write(f"{token}\n".encode())
        p.stdin.flush()

    flatpak_api.register_autostart(True)

    while True:
        time.sleep(60)
