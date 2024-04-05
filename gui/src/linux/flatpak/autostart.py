import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
from gi.repository import GLib, Gio
from random import randint
import sys
from threading import Timer


def receive_autostart(self, *args):
    print("autostart enabled..!?")
    print(args)
    sys.exit(0)

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
        'commandline': GLib.Variant('as', ['goldwarden_ui_main.py', '--hidden']),
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

thread = Timer(10, os._exit, [0])
thread.start()

loop = GLib.MainLoop()
loop.run()
