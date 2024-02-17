#!/usr/bin/python
import sys
import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc

from gi.repository import Gtk, Adw, GLib, Gdk, Gio
from ..services import goldwarden
from threading import Thread
import subprocess
from . import components
import os

root_path = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir))
token = sys.stdin.readline()
goldwarden.create_authenticated_connection(None)

def quickaccess_button_clicked():
    p = subprocess.Popen(["python3", "-m", "src.gui.quickaccess"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
    p.stdin.write(f"{token}\n".encode())
    p.stdin.flush()

def shortcuts_button_clicked():
    p = subprocess.Popen(["python3", "-m", "src.gui.shortcuts"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
    p.stdin.write(f"{token}\n".encode())
    p.stdin.flush()

def ssh_button_clicked():
    p = subprocess.Popen(["python3", "-m", "src.gui.ssh"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, cwd=root_path, start_new_session=True)
    p.stdin.write(f"{token}\n".encode())
    p.stdin.flush()

def add_action_row(parent, title, subtitle, icon=None):
    row = Adw.ActionRow()
    row.set_title(title)
    row.set_subtitle(subtitle)
    if icon != None:
        row.set_icon_name(icon)
    parent.add(row)
    return row

class SettingsWinvdow(Gtk.ApplicationWindow):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        # vertical box
        self.box = Gtk.Box()
        self.box.set_orientation(Gtk.Orientation.VERTICAL)
        self.set_child(self.box)
        
        def set_pin():
            set_pin_thread = Thread(target=goldwarden.enable_pin)
            set_pin_thread.start()

        self.banner = Adw.Banner()
        self.banner.set_title("No pin set, please set it now")
        self.banner.set_button_label("Set Pin")
        self.banner.connect("button-clicked", lambda banner: set_pin())
        self.box.append(self.banner)

        self.stack = Gtk.Stack()
        self.stack.set_transition_type(Gtk.StackTransitionType.SLIDE_LEFT_RIGHT)
        self.box.append(self.stack)

        self.preferences_page = Adw.PreferencesPage()
        self.preferences_page.set_title("General")
        self.stack.add_named(self.preferences_page, "preferences_page")

        self.action_preferences_group = Adw.PreferencesGroup()
        self.action_preferences_group.set_title("Actions")
        self.preferences_page.add(self.action_preferences_group)
        
        self.autotype_button = Gtk.Button()
        self.autotype_button.set_label("Quick Access")
        self.autotype_button.set_margin_top(10)
 
        self.autotype_button.connect("clicked", lambda button: quickaccess_button_clicked())
        self.autotype_button.get_style_context().add_class("suggested-action")
        self.action_preferences_group.add(self.autotype_button)

        self.login_button = Gtk.Button()
        self.login_button.set_label("Login")
        self.login_button.connect("clicked", lambda button: show_login())
        self.login_button.set_sensitive(False)
        self.login_button.set_margin_top(10)
        self.login_button.get_style_context().add_class("suggested-action")
        self.action_preferences_group.add(self.login_button)
    
        self.set_pin_button = Gtk.Button()
        self.set_pin_button.set_label("Set Pin")
        self.set_pin_button.connect("clicked", lambda button: set_pin())
        self.set_pin_button.set_margin_top(10)
        self.set_pin_button.set_sensitive(False)
        self.set_pin_button.get_style_context().add_class("suggested-action")
        self.action_preferences_group.add(self.set_pin_button)

        self.unlock_button = Gtk.Button()
        self.unlock_button.set_label("Unlock")
        self.unlock_button.set_margin_top(10)
        def unlock_button_clicked():
            action = goldwarden.unlock if self.unlock_button.get_label() == "Unlock" else goldwarden.lock
            unlock_thread = Thread(target=action)
            unlock_thread.start()
        self.unlock_button.connect("clicked", lambda button: unlock_button_clicked())
        # set disabled
        self.unlock_button.set_sensitive(False)
        self.action_preferences_group.add(self.unlock_button)

        self.logout_button = Gtk.Button()
        self.logout_button.set_label("Logout")
        self.logout_button.set_margin_top(10)
        self.logout_button.connect("clicked", lambda button: goldwarden.purge())
        self.logout_button.get_style_context().add_class("destructive-action")
        self.action_preferences_group.add(self.logout_button)

        self.wiki_button = Gtk.LinkButton(uri="https://github.com/quexten/goldwarden/wiki/Flatpak-Configuration")
        self.wiki_button.set_label("Help & Wiki")
        self.wiki_button.set_margin_top(10)
        self.action_preferences_group.add(self.wiki_button)

        self.vault_status_preferences_group = Adw.PreferencesGroup()
        self.vault_status_preferences_group.set_title("Vault Status")
        self.preferences_page.add(self.vault_status_preferences_group)
        
        self.status_row = add_action_row(self.vault_status_preferences_group, "Vault Status", "Locked")

        self.vault_status_icon = components.StatusIcon()
        self.vault_status_icon.set_icon("dialog-error", "error")
        self.status_row.add_prefix(self.vault_status_icon)

        self.last_sync_row = add_action_row(self.vault_status_preferences_group, "Last Sync", "Never", "emblem-synchronizing-symbolic")
        self.websocket_connected_row = add_action_row(self.vault_status_preferences_group, "Websocket Connected", "False")

        self.websocket_connected_status_icon = components.StatusIcon()
        self.websocket_connected_status_icon.set_icon("dialog-error", "error")
        self.websocket_connected_row.add_prefix(self.websocket_connected_status_icon)

        self.login_row = add_action_row(self.vault_status_preferences_group, "Vault Login Entries", "0", "dialog-password-symbolic")
        self.notes_row = add_action_row(self.vault_status_preferences_group, "Vault Notes", "0", "emblem-documents-symbolic")
 
        self.header = Gtk.HeaderBar()
        self.set_titlebar(self.header)

        action = Gio.SimpleAction.new("shortcuts", None)
        action.connect("activate", lambda action, parameter: shortcuts_button_clicked())
        self.add_action(action)
        menu = Gio.Menu.new()
        menu.append("Keyboard Shortcuts", "win.shortcuts") 
        self.popover = Gtk.PopoverMenu() 
        self.popover.set_menu_model(menu)

        action = Gio.SimpleAction.new("ssh", None)
        action.connect("activate", lambda action, parameter: ssh_button_clicked())
        self.add_action(action)
        menu.append("SSH Agent", "win.ssh")
        
        self.hamburger = Gtk.MenuButton()
        self.hamburger.set_popover(self.popover)
        self.hamburger.set_icon_name("open-menu-symbolic")
        self.header.pack_start(self.hamburger)


        def update_labels():
            pin_set = goldwarden.is_pin_enabled()
            status = goldwarden.get_vault_status()
            print("status", status)
            runtimeCfg = goldwarden.get_runtime_config()
            if runtimeCfg != None:
                self.ssh_row.set_subtitle("Listening at "+runtimeCfg["SSHAgentSocketPath"])
                self.goldwarden_daemon_row.set_subtitle("Listening at "+runtimeCfg["goldwardenSocketPath"])

            if status != None:
                if pin_set:
                    self.unlock_button.set_sensitive(True)
                    self.banner.set_revealed(False)
                else:
                    self.unlock_button.set_sensitive(False)
                    self.banner.set_revealed(True)
                logged_in = status["loggedIn"]
                if logged_in and not status["locked"]:
                    self.preferences_group.set_visible(True)
                    self.shortcut_preferences_group.set_visible(True)
                    self.autotype_button.set_visible(True)
                    self.login_row.set_sensitive(True)
                    self.notes_row.set_sensitive(True)
                    self.websocket_connected_row.set_sensitive(True)
                else:
                    self.preferences_group.set_visible(False)
                    self.shortcut_preferences_group.set_visible(False)
                    self.autotype_button.set_visible(False)
                    self.websocket_connected_row.set_sensitive(False)
                    self.login_row.set_sensitive(False)
                    self.notes_row.set_sensitive(False)

                locked = status["locked"]
                self.login_button.set_sensitive(pin_set and not locked)
                self.set_pin_button.set_sensitive(not pin_set or not locked)
                self.autotype_button.set_sensitive(not locked)
                self.status_row.set_subtitle(str("Logged in" if (logged_in and not locked) else "Logged out") if not locked else "Locked")
                if locked or not logged_in:
                    self.vault_status_icon.set_icon("dialog-warning", "warning")
                else:
                    self.vault_status_icon.set_icon("emblem-default", "ok")
                if not logged_in:
                    self.logout_button.set_sensitive(False)
                else:
                    self.logout_button.set_sensitive(True)
                self.login_row.set_subtitle(str(status["loginEntries"]))
                self.notes_row.set_subtitle(str(status["noteEntries"]))
                self.websocket_connected_row.set_subtitle("Connected" if status["websocketConnected"] else "Disconnected")
                if status["websocketConnected"]:
                    self.websocket_connected_status_icon.set_icon("emblem-default", "ok")
                else:
                    self.websocket_connected_status_icon.set_icon("dialog-error", "error")
                self.last_sync_row.set_subtitle(str(status["lastSynced"]))
                if status["lastSynced"].startswith("1970") or status["lastSynced"].startswith("1969"):
                    self.last_sync_row.set_subtitle("Never")
                self.unlock_button.set_label("Unlock" if locked else "Lock")
            else:
                is_daemon_running = goldwarden.is_daemon_running()
                if not is_daemon_running:
                    self.status_row.set_subtitle("Daemon not running")
                    self.vault_status_icon.set_icon("dialog-error", "error")
            
            GLib.timeout_add(5000, update_labels)

        GLib.timeout_add(1000, update_labels)
        self.set_default_size(400, 700)
        self.set_title("Goldwarden")

class MyApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def on_activate(self, app):
        self.settings_win = SettingsWinvdow(application=app)
        self.settings_win.present()

def show_login():
    dialog = Gtk.Dialog(title="Goldwarden")

    auth_preference_group = Adw.PreferencesGroup()
    auth_preference_group.set_title("Authentication")
    auth_preference_group.set_margin_top(10)
    auth_preference_group.set_margin_bottom(10)
    auth_preference_group.set_margin_start(10)
    auth_preference_group.set_margin_end(10)
    dialog.get_content_area().append(auth_preference_group)

    email_entry = Adw.EntryRow()
    email_entry.set_title("Email")
    email_entry.set_text("")
    auth_preference_group.add(email_entry)

    client_id_entry = Adw.EntryRow()
    client_id_entry.set_title("Client ID (optional)")
    client_id_entry.set_text("")
    auth_preference_group.add(client_id_entry)

    client_secret_entry = Adw.EntryRow()
    client_secret_entry.set_title("Client Secret (optional)")
    client_secret_entry.set_text("")
    auth_preference_group.add(client_secret_entry)

    dialog.add_button("Login", Gtk.ResponseType.OK)
    def on_save(res):
        if res != Gtk.ResponseType.OK:
            return
        goldwarden.set_url(url_entry.get_text())
        goldwarden.set_client_id(client_id_entry.get_text())
        goldwarden.set_client_secret(client_secret_entry.get_text())
        def login():
            res = goldwarden.login_with_password(email_entry.get_text(), "password")
            def handle_res():
                if res == "ok":
                    dialog.close()
                elif res == "badpass":
                    bad_pass_diag = Gtk.MessageDialog(transient_for=dialog, modal=True, message_type=Gtk.MessageType.ERROR, buttons=Gtk.ButtonsType.OK, text="Bad password")
                    bad_pass_diag.connect("response", lambda dialog, response: bad_pass_diag.close())
                    bad_pass_diag.present()
            GLib.idle_add(handle_res)

        login_thread = Thread(target=login)
        login_thread.start()

    preference_group = Adw.PreferencesGroup()
    preference_group.set_title("Config")
    preference_group.set_margin_top(10)
    preference_group.set_margin_bottom(10)
    preference_group.set_margin_start(10)
    preference_group.set_margin_end(10)

    dialog.get_content_area().append(preference_group)

    url_entry = Adw.EntryRow()
    url_entry.set_title("Base Url")
    url_entry.set_text("https://vault.bitwarden.com/")
    preference_group.add(url_entry)

    #ok response
    dialog.connect("response", lambda dialog, response: on_save(response))
    dialog.set_default_size(400, 200)
    dialog.set_modal(True)
    dialog.present()

isflatpak = os.path.exists("/.flatpak-info")
pathprefix = "/app/bin/" if isflatpak else "./"
css_provider = Gtk.CssProvider()
css_provider.load_from_path(pathprefix+"style.css")
Gtk.StyleContext.add_provider_for_display(
    Gdk.Display.get_default(),
    css_provider,
    Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION
)

app = MyApp(application_id="com.quexten.Goldwarden.settings")
app.run(sys.argv)