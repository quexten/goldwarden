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
        bus_name = dbus.service.BusName('com.quexten.goldwarden', bus=dbus.SessionBus())
        dbus.service.Object.__init__(self, bus_name, '/com/quexten/goldwarden')

    @dbus.service.method('com.quexten.goldwarden.Autofill')
    def autofill(self):
        on_autofill()
        return ""

DBusGMainLoop(set_as_default=True)
service = GoldwardenDBUSService()
