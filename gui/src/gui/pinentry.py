import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc
import time
from gi.repository import Gtk, Adw, GLib, Notify, Gdk
from threading import Thread
from .template_loader import load_template
import sys
import os

class GoldwardenPinentryApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def on_activate(self, app):
        self.load()
        self.window.present()

    def load(self):
        builder = load_template("pinentry.ui")
        self.window = builder.get_object("window")
        self.message_label = builder.get_object("message")
        self.message_label.set_label(self.message)
        
        self.cancel_button = builder.get_object("cancel_button")
        self.cancel_button.connect("clicked", self.on_cancel_button_clicked)
        self.approve_button = builder.get_object("approve_button")  
        self.approve_button.connect("clicked", self.on_approve_button_clicked)

        self.password_entry = builder.get_object("password_entry")
        self.password_entry.set_placeholder_text("Enter your password")

        self.window.set_application(self)

    def on_approve_button_clicked(self, button):
        print(self.password_entry.get_text(), flush=True)
        os._exit(0)

    def on_cancel_button_clicked(self, button):
        print("", flush=True)
        os._exit(0)

if __name__ == "__main__":
    app = GoldwardenPinentryApp(application_id="com.quexten.Goldwarden.pinentry")
    message = sys.stdin.readline()
    app.message = message
    app.run(sys.argv)