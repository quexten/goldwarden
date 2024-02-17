import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc
import time
from gi.repository import Gtk, Adw, GLib, Notify, Gdk
from threading import Thread
import sys
import os

message = sys.stdin.readline()

class MyApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def on_activate(self, app):
        self.pinentry_window = MainWindow(application=app)
        self.pinentry_window.present()
        self.app = app

class MainWindow(Gtk.ApplicationWindow):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.stack = Gtk.Stack()
        self.stack.set_transition_type(Gtk.StackTransitionType.SLIDE_LEFT_RIGHT)
        self.set_child(self.stack)

        box = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing=6)
        self.stack.add_child(box)

        label = Gtk.Label(label=message)
        box.append(label)

        # Create a button box for cancel and approve buttons
        button_box = Gtk.Box(spacing=6)
        box.append(button_box)

        # Cancel button
        cancel_button = Gtk.Button(label="Cancel")
        cancel_button.set_hexpand(True)  # Make the button expand horizontally
        def on_cancel_button_clicked(button):
            print("false", flush=True)
            os._exit(0)
        cancel_button.connect("clicked", on_cancel_button_clicked)
        button_box.append(cancel_button)

        # Approve button
        approve_button = Gtk.Button(label="Approve")
        approve_button.set_hexpand(True)  # Make the button expand horizontally
        def on_approve_button_clicked(button):
            print("true", flush=True)
            os._exit(0)
        approve_button.connect("clicked", on_approve_button_clicked)
        button_box.append(approve_button)

        self.set_default_size(700, 200)
        self.set_title("Goldwarden Approval")

app = MyApp(application_id="com.quexten.Goldwarden.pinentry")
app.run(sys.argv)