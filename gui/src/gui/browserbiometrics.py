import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc
import time
from gi.repository import Gtk, Adw, GLib, Notify, Gdk
from threading import Thread
import sys
import os
from . import components

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

        # vertical box
        self.box = Gtk.Box()
        self.box.set_orientation(Gtk.Orientation.VERTICAL)
        self.set_child(self.box)
        
        self.stack = Gtk.Stack()
        self.stack.set_transition_type(Gtk.StackTransitionType.SLIDE_LEFT_RIGHT)
        self.box.append(self.stack)

        self.preferences_page = Adw.PreferencesPage()
        self.preferences_page.set_title("General")
        self.stack.add_named(self.preferences_page, "preferences_page")

        self.register_browser_biometrics_group = Adw.PreferencesGroup()
        self.register_browser_biometrics_group.set_title("Register Browser Biometrics")
        self.register_browser_biometrics_group.set_description("Run the following command in your terminal to set up the browser biometrics integration")
        self.preferences_page.add(self.register_browser_biometrics_group)

        self.setup_command_row = Adw.ActionRow()
        self.setup_command_row.set_subtitle("flatpak run --filesystem=home --command=goldwarden com.quexten.Goldwarden setup browserbiometrics")
        self.setup_command_row.set_subtitle_selectable(True)
        self.register_browser_biometrics_group.add(self.setup_command_row)

        self.set_default_size(700, 400)
        self.set_title("Goldwarden Browser Biometrics Setup")

app = MyApp(application_id="com.quexten.Goldwarden.browserbiometrics")
app.run(sys.argv)