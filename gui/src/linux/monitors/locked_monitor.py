from gi.repository import Gtk
import dbus
import dbus.service
from dbus.mainloop.glib import DBusGMainLoop
from threading import Thread
import subprocess
import os
from src.services import goldwarden
import time
import src.linux.flatpak.api as flatpak_api

daemon_token = None

def daemon():
    time.sleep(5)
    goldwarden.create_authenticated_connection(daemon_token)
    while True:
        status = goldwarden.get_vault_status()
        if status["locked"]:
            flatpak_api.set_status("Locked")
        else:
            flatpak_api.set_status("Unlocked")
        time.sleep(1)

def run_daemon(token):
    print("running locked status daemon")
    global daemon_token
    daemon_token = token
    thread = Thread(target=daemon)
    thread.start()