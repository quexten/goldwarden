#!/usr/bin/env python3
import sys
import gi

gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')

from gi.repository import Gtk, Adw, GLib, Gdk, Gio
from ..services import goldwarden
from threading import Thread
from .resource_loader import load_template
import subprocess
import os

def run_window(name, token):
    gui_path = os.path.dirname(os.path.realpath(__file__))
    cwd = os.path.abspath(os.path.join(gui_path, os.pardir, os.pardir))
    print(f"Running window {name} with path {cwd}")
    p = subprocess.Popen(["python3", "-m", "src.gui." + name], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=cwd, start_new_session=True)
    p.stdin.write(f"{token}\n".encode())
    p.stdin.flush()


class GoldwardenSettingsApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def on_activate(self, app):
        self.load()
        self.update()
        self.window.present()
        GLib.timeout_add(100, self.update)

    def load(self):
        builder = load_template("settings.ui")
        self.window = builder.get_object("window")
        self.window.set_application(self)
        self.stack = builder.get_object("stack")

        self.set_pin_status_box = builder.get_object("set_pin_status")
        self.set_pin_button = builder.get_object("set_pin_button")
        self.set_pin_button.connect("clicked", lambda x: goldwarden.enable_pin())

        self.unlock_status_box = builder.get_object("unlock_status")
        self.unlock_button = builder.get_object("unlock_button")
        self.unlock_button.connect("clicked", lambda x: goldwarden.unlock())
        
        self.login_status_box = builder.get_object("login_status")
        self.login_button = builder.get_object("login_button")
        self.login_button.connect("clicked", lambda x: run_window("login", self.token))

        self.settings_view = builder.get_object("settings_view")
        self.lock_button = builder.get_object("lock_button")
        self.lock_button.connect("clicked", lambda x: goldwarden.lock())
        self.logout_button = builder.get_object("logout_button")
        self.logout_button.connect("clicked", lambda x: goldwarden.purge())
        self.update_pin_button = builder.get_object("update_pin_button")
        self.update_pin_button.connect("clicked", lambda x: goldwarden.enable_pin())
        self.quickaccess_button = builder.get_object("quickaccess_button")
        self.quickaccess_button.connect("clicked", lambda x: run_window("quickaccess", self.token))
        self.last_sync_row = builder.get_object("last_sync_row")
        self.websocket_connected_row = builder.get_object("websocket_connected_row")
        self.logins_row = builder.get_object("logins_row")
        self.notes_row = builder.get_object("notes_row")

        self.menu_button = builder.get_object("menu_button")
        menu = Gio.Menu.new()
        self.popover = Gtk.PopoverMenu() 
        self.popover.set_menu_model(menu)
        self.menu_button.set_popover(self.popover)

        action = Gio.SimpleAction.new("shortcuts", None)
        action.connect("activate", lambda action, parameter: run_window("shortcuts", self.token))
        self.window.add_action(action)
        menu.append("Keyboard Shortcuts", "win.shortcuts") 

        action = Gio.SimpleAction.new("ssh", None)
        action.connect("activate", lambda action, parameter: run_window("ssh", self.token))
        self.window.add_action(action)
        menu.append("SSH Agent", "win.ssh")

        action = Gio.SimpleAction.new("browserbiometrics", None)
        action.connect("activate", lambda action, parameter: run_window("browserbiometrics", self.token))
        self.window.add_action(action)
        menu.append("Browser Biometrics", "win.browserbiometrics")

        action = Gio.SimpleAction.new("about", None)
        action.connect("activate", lambda action, parameter: self.show_about())
        self.window.add_action(action)
        menu.append("About", "win.about")
    
    def update(self):
        self.render()
        return True

    def render(self):
        pin_set = goldwarden.is_pin_enabled()
        status = goldwarden.get_vault_status()
        runtimeCfg = goldwarden.get_runtime_config()
        if status == None:
            is_daemon_running = goldwarden.is_daemon_running()
            if not is_daemon_running:
                self.status_row.set_subtitle("Daemon not running")
                self.vault_status_icon.set_icon("dialog-error", "error")
            return

        logged_in = status["loggedIn"]
        unlocked = not status["locked"]
        if not pin_set:
            self.stack.set_visible_child(self.set_pin_status_box)
            return
        if not unlocked:
            self.stack.set_visible_child(self.unlock_status_box)
            return
        if not logged_in:
            self.stack.set_visible_child(self.login_status_box)
            return
        self.stack.set_visible_child(self.settings_view)

        self.last_sync_row.set_subtitle(status["lastSynced"])
        self.websocket_connected_row.set_subtitle("Yes" if status["websocketConnected"] else "No")
        self.logins_row.set_subtitle(str(status["loginEntries"]))
        self.notes_row.set_subtitle(str(status["noteEntries"]))

    def show_about(self):
        dialog = Adw.AboutWindow(transient_for=app.get_active_window()) 
        dialog.set_application_name("Goldwarden") 
        dialog.set_version(goldwarden.version())
        dialog.set_developer_name("Bernd Schoolmann (Quexten)") 
        dialog.set_license_type(Gtk.License(Gtk.License.MIT_X11)) 
        dialog.set_comments("A Bitwarden compatible password manager") 
        dialog.set_website("https://github.com/quexten/goldwarden") 
        dialog.set_issue_url("https://github.com/quexten/goldwarden/issues") 
        dialog.add_credit_section("Contributors", ["Bernd Schoolmann"]) 
        dialog.set_copyright("© 2024 Bernd Schoolmann") 
        dialog.set_developers(["Bernd Schoolmann"]) 
        dialog.set_application_icon("com.quexten.Goldwarden")
        dialog.set_visible(True)

if __name__ == "__main__":
    settings = Gtk.Settings.get_default()
    settings.set_property("gtk-error-bell", False)

    token = sys.stdin.readline().strip()

    goldwarden.create_authenticated_connection(None)
    app = GoldwardenSettingsApp(application_id="com.quexten.Goldwarden.settings")
    app.token = token
    app.run(sys.argv)
