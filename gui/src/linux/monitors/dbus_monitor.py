from gi.repository import Gtk
import dbus
import dbus.service
from dbus.mainloop.glib import DBusGMainLoop
from threading import Thread
import subprocess
import os

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir, os.pardir))
daemon_token = None

class GoldwardenDBUSService(dbus.service.Object):
    def __init__(self):
        bus_name = dbus.service.BusName('com.quexten.Goldwarden.ui', bus=dbus.SessionBus())
        dbus.service.Object.__init__(self, bus_name, '/com/quexten/Goldwarden/ui')

    @dbus.service.method('com.quexten.Goldwarden.ui.QuickAccess')
    def quickaccess(self):
        p = subprocess.Popen(["python3", "-m", "src.gui.quickaccess"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
        p.stdin.write(f"{daemon_token}\n".encode())
        p.stdin.flush()
        return ""

    @dbus.service.method('com.quexten.Goldwarden.ui.Settings')
    def settings(self):
        subprocess.Popen(["python3", "-m", "src.gui.settings"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
        p.stdin.write(f"{daemon_token}\n".encode())
        p.stdin.flush()
        return ""

def daemon():
    DBusGMainLoop(set_as_default=True)
    service = GoldwardenDBUSService()
    from gi.repository import GLib, GObject as gobject
    gobject.MainLoop().run()

def run_daemon(token):
    global daemon_token
    daemon_token = token
    thread = Thread(target=daemon)
    thread.start()