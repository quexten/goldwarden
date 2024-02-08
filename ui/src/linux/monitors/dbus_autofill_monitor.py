#Python DBUS Test Server
#runs until the Quit() method is called via DBUS

from gi.repository import Gtk
import dbus
import dbus.service
from dbus.mainloop.glib import DBusGMainLoop
from threading import Thread

on_autofill = lambda: None

class GoldwardenDBUSService(dbus.service.Object):
    def __init__(self):
        bus_name = dbus.service.BusName('com.quexten.Goldwarden.autofill', bus=dbus.SessionBus())
        dbus.service.Object.__init__(self, bus_name, '/com/quexten/Goldwarden')

    @dbus.service.method('com.quexten.Goldwarden.Autofill')
    def autofill(self):
        on_autofill()
        return ""

def daemon():
    DBusGMainLoop(set_as_default=True)
    service = GoldwardenDBUSService()
    from gi.repository import GLib, GObject as gobject
    gobject.MainLoop().run()

def run_daemon():
    thread = Thread(target=daemon)
    thread.start()
    