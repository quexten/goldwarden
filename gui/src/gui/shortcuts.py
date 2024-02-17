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

def add_action_row(parent, title, subtitle, icon=None):
    row = Adw.ActionRow()
    row.set_title(title)
    row.set_subtitle(subtitle)
    if icon != None:
        row.set_icon_name(icon)
    parent.add(row)
    return row

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

        self.global_preferences_group = Adw.PreferencesGroup()
        self.global_preferences_group.set_title("Global Shortcuts")
        self.preferences_page.add(self.global_preferences_group)

        self.autofill_row = Adw.ActionRow()
        self.autofill_row.set_title("Autofill Shortcut")
        self.autofill_row.set_subtitle("Not implemented - check the wiki for manual setup")
        self.global_preferences_group.add(self.autofill_row)

        self.autofill_icon = components.StatusIcon()
        self.autofill_icon.set_icon("dialog-warning", "warning")
        self.autofill_row.add_prefix(self.autofill_icon)

        self.quickaccess_preferences_group = Adw.PreferencesGroup()
        self.quickaccess_preferences_group.set_title("Quick Access Shortcuts")
        self.preferences_page.add(self.quickaccess_preferences_group)

        add_action_row(self.quickaccess_preferences_group, "Copy Username Shortcut", "CTRL + U")
        add_action_row(self.quickaccess_preferences_group, "Autotype Username Shortcut", "CTRL + ALT + U")
        add_action_row(self.quickaccess_preferences_group, "Copy Password Shortcut", "CTRL + P")
        add_action_row(self.quickaccess_preferences_group, "Autotype Password Shortcut", "CTRL + ALT + P")
        add_action_row(self.quickaccess_preferences_group, "Copy TOTP Shortcut", "CTRL + T")
        add_action_row(self.quickaccess_preferences_group, "Autotype TOTP Shortcut", "CTRL + ALT + T")
        add_action_row(self.quickaccess_preferences_group, "Launch URI Shortcut", "CTRL+L")
        add_action_row(self.quickaccess_preferences_group, "Launch Web Vault Shortcut", "CTRL+V")
        add_action_row(self.quickaccess_preferences_group, "Focus Search Shortcut", "F")
        add_action_row(self.quickaccess_preferences_group, "Quit Shortcut", "Esc")

        self.set_default_size(700, 700)
        self.set_title("Goldwarden Shortcuts")

app = MyApp(application_id="com.quexten.Goldwarden.shortcuts")
app.run(sys.argv)