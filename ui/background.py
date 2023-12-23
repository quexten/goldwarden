import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc
from gi.repository import Gtk, Adw, GLib, Gio
from random import randint
import time


def receive_autostart(self, *args):
    print("autostart enabled..!?")
    print(args)

def request_autostart():
    bus = Gio.bus_get_sync(Gio.BusType.SESSION, None)
    proxy = Gio.DBusProxy.new_sync(
        bus,
        Gio.DBusProxyFlags.NONE,
        None,
        'org.freedesktop.portal.Desktop',
        '/org/freedesktop/portal/desktop',
        'org.freedesktop.portal.Background',
        None,
    )

    token = 0 + randint(10000000, 20000000)
    options = {
        'handle_token': GLib.Variant('s', f'com/quexten/Goldwarden/{token}'),
        'reason': GLib.Variant('s', ('Autostart Goldwarden in the background.')),
        'autostart': GLib.Variant('b', True),
        'commandline': GLib.Variant('as', ['main.py', '--hidden']),
        'dbus-activatable': GLib.Variant('b', False),
    }

    try:
        request = proxy.RequestBackground('(sa{sv})', "", options)
        if request is None:
            raise Exception(
                "Registering with background portal failed."
            )

        bus.signal_subscribe(
            'org.freedesktop.portal.Desktop',
            'org.freedesktop.portal.Request',
            'Response',
            request,
            None,
            Gio.DBusSignalFlags.NO_MATCH_RULE,
            receive_autostart,
            None,
        )
    except Exception as e:
        print(e)

request_autostart()

loop = GLib.MainLoop()
loop.run()
