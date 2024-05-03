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
from ..services import goldwarden

goldwarden.create_authenticated_connection(None)

class GoldwardenLoginApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def on_activate(self, app):
        self.load()
        self.window.present()

    def load(self):
        builder = load_template("login.ui")
        self.window = builder.get_object("window")
        self.window.set_application(self)
        self.email_row = builder.get_object("email_row")
        self.client_id_row = builder.get_object("client_id_row")
        self.client_secret_row = builder.get_object("client_secret_row")
        self.server_row = builder.get_object("server_row")
        self.login_button = builder.get_object("login_button")
        self.login_button.connect("clicked", lambda x: self.on_login())
        
    def on_login(self):
        email = self.email_row.get_text()
        client_id = self.client_id_row.get_text()
        client_secret = self.client_secret_row.get_text()
        server = self.server_row.get_text()
        goldwarden.set_url(server)
        if client_id != "":
            goldwarden.set_client_id(client_id)
        if client_secret != "":
            goldwarden.set_client_secret(client_secret)
        goldwarden.login_with_password(email, "")

if __name__ == "__main__":
    settings = Gtk.Settings.get_default()
    settings.set_property("gtk-error-bell", False)

    app = GoldwardenLoginApp(application_id="com.quexten.Goldwarden.login")
    app.run(sys.argv)