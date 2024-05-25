"""
Script to set the status of the background process.
Run separately so that gtk dependencies don't stick around in memory.
"""
import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
from gi.repository import GLib, Gio
import sys
from threading import Timer

def set_status(message):
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

    options = {
        'message': GLib.Variant('s', message),
    }

    try:
        request = proxy.SetStatus('(a{sv})', options)
        sys.exit(0)
    except Exception as e:
        print(e)
        sys.exit(0)

if len(sys.argv) > 1:
    set_status(sys.argv[1])

thread = Timer(10, sys.exit, [0])
thread.start()

loop = GLib.MainLoop()
loop.run()
