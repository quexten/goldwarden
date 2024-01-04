#!/usr/bin/python
import sys
import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc

from gi.repository import Gtk, Adw, GLib, Gdk
import goldwarden
from threading import Thread
import subprocess
import components
import os

class SettingsWinvdow(Gtk.ApplicationWindow):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.stack = Gtk.Stack()
        self.stack.set_transition_type(Gtk.StackTransitionType.SLIDE_LEFT_RIGHT)
        self.set_child(self.stack)

        self.preferences_page = Adw.PreferencesPage()
        self.preferences_page.set_title("General")
        self.stack.add_named(self.preferences_page, "preferences_page")

        self.preferences_group = Adw.PreferencesGroup()
        self.preferences_group.set_title("Services")
        self.preferences_page.add(self.preferences_group)

        self.ssh_row = Adw.ActionRow()
        self.ssh_row.set_title("SSH Daemon")
        self.ssh_row.set_subtitle("Getting status...")
        self.ssh_row.set_icon_name("emblem-default")
        self.preferences_group.add(self.ssh_row)

        self.goldwarden_daemon_row = Adw.ActionRow()
        self.goldwarden_daemon_row.set_title("Goldwarden Daemon")
        self.goldwarden_daemon_row.set_subtitle("Getting status...")
        self.goldwarden_daemon_row.set_icon_name("emblem-default")
        self.preferences_group.add(self.goldwarden_daemon_row)

        self.login_with_device = Adw.ActionRow()
        self.login_with_device.set_title("Login with device")
        self.login_with_device.set_subtitle("Waiting for requests...")
        self.preferences_group.add(self.login_with_device)

        self.status_row = Adw.ActionRow()
        self.status_row.set_title("DBUS Service")
        self.status_row.set_subtitle("Listening")
        self.preferences_group.add(self.status_row)

        self.shortcut_preferences_group = Adw.PreferencesGroup()
        self.shortcut_preferences_group.set_title("Shortcuts")
        self.preferences_page.add(self.shortcut_preferences_group)

        self.autofill_row = Adw.ActionRow()
        self.autofill_row.set_title("Autofill Shortcut")
        self.autofill_row.set_subtitle("Unavailable, please set up a shortcut in your desktop environment (README)")
        self.shortcut_preferences_group.add(self.autofill_row)

        self.autofill_icon = components.StatusIcon()
        self.autofill_icon.set_icon("dialog-warning", "warning")
        self.autofill_row.add_prefix(self.autofill_icon)

        self.copy_username_shortcut_row = Adw.ActionRow()
        self.copy_username_shortcut_row.set_title("Copy Username Shortcut")
        self.copy_username_shortcut_row.set_subtitle("U")
        self.shortcut_preferences_group.add(self.copy_username_shortcut_row)

        self.copy_password_shortcut_row = Adw.ActionRow()
        self.copy_password_shortcut_row.set_title("Copy Password Shortcut")
        self.copy_password_shortcut_row.set_subtitle("P")
        self.shortcut_preferences_group.add(self.copy_password_shortcut_row)

        self.vault_status_preferences_group = Adw.PreferencesGroup()
        self.vault_status_preferences_group.set_title("Vault Status")
        self.preferences_page.add(self.vault_status_preferences_group)
        
        self.status_row = Adw.ActionRow()
        self.status_row.set_title("Vault Status")
        self.status_row.set_subtitle("Locked")
        self.vault_status_preferences_group.add(self.status_row)

        self.vault_status_icon = components.StatusIcon()
        self.vault_status_icon.set_icon("dialog-error", "error")
        self.status_row.add_prefix(self.vault_status_icon)

        self.last_sync_row = Adw.ActionRow()
        self.last_sync_row.set_title("Last Sync")
        self.last_sync_row.set_subtitle("Never")
        self.last_sync_row.set_icon_name("emblem-synchronizing-symbolic")
        self.vault_status_preferences_group.add(self.last_sync_row)

        self.websocket_connected_row = Adw.ActionRow()
        self.websocket_connected_row.set_title("Websocket Connected")
        self.websocket_connected_row.set_subtitle("False")
        self.vault_status_preferences_group.add(self.websocket_connected_row)

        self.websocket_connected_status_icon = components.StatusIcon()
        self.websocket_connected_status_icon.set_icon("dialog-error", "error")
        self.websocket_connected_row.add_prefix(self.websocket_connected_status_icon)
        
        self.login_row = Adw.ActionRow()
        self.login_row.set_title("Vault Login Entries")
        self.login_row.set_subtitle("0")
        self.login_row.set_icon_name("dialog-password-symbolic")
        self.vault_status_preferences_group.add(self.login_row)

        self.notes_row = Adw.ActionRow()
        self.notes_row.set_title("Vault Notes")
        self.notes_row.set_subtitle("0")
        self.notes_row.set_icon_name("emblem-documents-symbolic")
        self.vault_status_preferences_group.add(self.notes_row)

        self.action_preferences_group = Adw.PreferencesGroup()
        self.action_preferences_group.set_title("Actions")
        self.preferences_page.add(self.action_preferences_group)
        
        self.autotype_button = Gtk.Button()
        self.autotype_button.set_label("Autotype")
        self.autotype_button.set_margin_top(10)
        self.autotype_button.connect("clicked", lambda button: subprocess.Popen(["python3", "/app/bin/autofill.py"], start_new_session=True))
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
        def set_pin():
            set_pin_thread = Thread(target=goldwarden.enable_pin)
            set_pin_thread.start()
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
        
        def update_labels():
            GLib.timeout_add(1000, update_labels)
            
            pin_set = goldwarden.is_pin_enabled()
            status = goldwarden.get_vault_status()
            runtimeCfg = goldwarden.get_runtime_config()
            if runtimeCfg != None:
                self.ssh_row.set_subtitle("Listening at "+runtimeCfg["SSHAgentSocketPath"])
                self.goldwarden_daemon_row.set_subtitle("Listening at "+runtimeCfg["goldwardenSocketPath"])

            if status != None:
                if pin_set:
                    self.unlock_button.set_sensitive(True)
                else:
                    self.unlock_button.set_sensitive(False)
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

        GLib.timeout_add(1000, update_labels)
        self.set_default_size(400, 700)
        self.set_title("Goldwarden")


        #add title buttons
        self.title_bar = Gtk.HeaderBar()
        self.set_titlebar(self.title_bar)

class MyApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def on_activate(self, app):
        self.settings_win = SettingsWinvdow(application=app)
        self.settings_win.present()

def show_login():
    dialog = Gtk.Dialog(title="Goldwarden")
    preference_group = Adw.PreferencesGroup()
    preference_group.set_title("Config")
    preference_group.set_margin_top(10)
    preference_group.set_margin_bottom(10)
    preference_group.set_margin_start(10)
    preference_group.set_margin_end(10)

    dialog.get_content_area().append(preference_group)

    api_url_entry = Adw.EntryRow()
    api_url_entry.set_title("API Url")
    # set value
    api_url_entry.set_text("https://vault.bitwarden.com/api")
    preference_group.add(api_url_entry)

    identity_url_entry = Adw.EntryRow()
    identity_url_entry.set_title("Identity Url")
    identity_url_entry.set_text("https://vault.bitwarden.com/identity")
    preference_group.add(identity_url_entry)

    notification_url_entry = Adw.EntryRow()
    notification_url_entry.set_title("Notification URL")
    notification_url_entry.set_text("https://notifications.bitwarden.com/")
    preference_group.add(notification_url_entry)

    auth_preference_group = Adw.PreferencesGroup()
    auth_preference_group.set_title("Authentication")
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
        goldwarden.set_api_url(api_url_entry.get_text())
        goldwarden.set_identity_url(identity_url_entry.get_text())
        goldwarden.set_notification_url(notification_url_entry.get_text())
        goldwarden.set_client_id(client_id_entry.get_text())
        goldwarden.set_client_secret(client_secret_entry.get_text())
        def login():
            res = goldwarden.login_with_password(email_entry.get_text(), "password")
            def handle_res():
                print("handle res", res)
                if res == "ok":
                    dialog.close()
                    print("ok")
                elif res == "badpass":
                    bad_pass_diag = Gtk.MessageDialog(transient_for=dialog, modal=True, message_type=Gtk.MessageType.ERROR, buttons=Gtk.ButtonsType.OK, text="Bad password")
                    bad_pass_diag.connect("response", lambda dialog, response: bad_pass_diag.close())
                    bad_pass_diag.present()
            GLib.idle_add(handle_res)

        login_thread = Thread(target=login)
        login_thread.start()

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