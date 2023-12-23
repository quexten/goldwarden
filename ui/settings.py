#!/usr/bin/python
import sys
import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc

from gi.repository import Gtk, Adw, GLib
import goldwarden
from threading import Thread

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
        self.ssh_row.set_subtitle("Listening at ~/.goldwarden-ssh-agent.sock")
        self.preferences_group.add(self.ssh_row)

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

        self.login_row = Adw.ActionRow()
        self.login_row.set_title("Vault Login Entries")
        self.login_row.set_subtitle("0")
        self.vault_status_preferences_group.add(self.login_row)

        self.notes_row = Adw.ActionRow()
        self.notes_row.set_title("Vault Notes")
        self.notes_row.set_subtitle("0")
        self.vault_status_preferences_group.add(self.notes_row)

        self.action_preferences_group = Adw.PreferencesGroup()
        self.action_preferences_group.set_title("Actions")
        self.preferences_page.add(self.action_preferences_group)
        
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
        
        def update_labels():
            pin_set = goldwarden.is_pin_enabled()
            status = goldwarden.get_vault_status()
            locked = status["locked"]
            self.login_button.set_sensitive(pin_set and not locked)
            self.set_pin_button.set_sensitive(not pin_set or not locked)
            self.status_row.set_subtitle(str("Unlocked" if not locked else "Locked"))
            self.login_row.set_subtitle(str(status["loginEntries"]))
            self.notes_row.set_subtitle(str(status["noteEntries"]))
            self.unlock_button.set_sensitive(True)
            self.unlock_button.set_label("Unlock" if locked else "Lock")
            GLib.timeout_add(1000, update_labels)

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

app = MyApp(application_id="com.quexten.Goldwarden")

def show_login():
    dialog = Gtk.Dialog(title="Goldwarden")
    preference_group = Adw.PreferencesGroup()
    preference_group.set_title("Config")
    dialog.get_content_area().append(preference_group)

    api_url_entry = Adw.EntryRow()
    api_url_entry.set_title("API Url")
    # set value
    api_url_entry.set_text("https://api.bitwarden.com/")
    preference_group.add(api_url_entry)

    identity_url_entry = Adw.EntryRow()
    identity_url_entry.set_title("Identity Url")
    identity_url_entry.set_text("https://identity.bitwarden.com/")
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

    dialog.add_button("Login", Gtk.ResponseType.OK)
    def on_save(res):
        if res != Gtk.ResponseType.OK:
            return
        goldwarden.set_api_url(api_url_entry.get_text())
        goldwarden.set_identity_url(identity_url_entry.get_text())
        goldwarden.set_notification_url(notification_url_entry.get_text())
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

app.run(sys.argv)