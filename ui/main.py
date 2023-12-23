#!/usr/bin/python
import sys
import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
import gc

from gi.repository import Gtk, Adw, Gdk, Graphene, Gsk, Gio, GLib, GObject
import monitors.dbus_autofill_monitor
import goldwarden
import clipboard
import time
from threading import Thread

class MainWindow(Gtk.ApplicationWindow):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.stack = Gtk.Stack()
        self.stack.set_transition_type(Gtk.StackTransitionType.SLIDE_LEFT_RIGHT)
        self.set_child(self.stack)

        self.box = Gtk.Box()
        self.box.set_orientation(Gtk.Orientation.VERTICAL)
        self.stack.add_named(self.box, "box")


        self.text_view = Adw.EntryRow()
        self.text_view.set_title("Search")
        # on type func
        def on_type(entry):
            if len(entry.get_text()) > 1:
                self.history_list.show()
            else:
                self.history_list.hide()

            while self.history_list.get_first_child() != None:
                self.history_list.remove(self.history_list.get_first_child())

            self.filtered_logins = list(filter(lambda i: entry.get_text().lower() in i["name"].lower(), self.logins))
            if len( self.filtered_logins) > 10:
                 self.filtered_logins =  self.filtered_logins[0:10]
            self.starts_with_logins = list(filter(lambda i: i["name"].lower().startswith(entry.get_text().lower()), self.logins))
            self.other_logins = list(filter(lambda i: i not in self.starts_with_logins ,  self.filtered_logins))
            self.filtered_logins = None

            for i in self.starts_with_logins  + self.other_logins :
                action_row = Adw.ActionRow()
                action_row.set_title(i["name"])
                action_row.set_subtitle(i["username"])
                action_row.set_icon_name("dialog-password")
                action_row.set_activatable(True)
                action_row.password = i["password"]
                action_row.username = i["username"]
                self.history_list.append(action_row)
            self.starts_with_logins = None
            self.other_logins = None
        self.text_view.connect("changed", lambda entry: on_type(entry))
        self.box.append(self.text_view)
    
        self.history_list = Gtk.ListBox()
        # margin'
        self.history_list.set_margin_start(10)
        self.history_list.set_margin_end(10)
        self.history_list.set_margin_top(10)
        self.history_list.set_margin_bottom(10)
        self.history_list.hide()

        keycont = Gtk.EventControllerKey()
        def handle_keypress(controller, keyval, keycode, state, user_data):
            if keycode == 36:
                print("enter")
                self.hide()
                def do_autotype(username, password):
                    time.sleep(0.5)
                    goldwarden.autotype(username, password)
                    GLib.idle_add(lambda: self.show())
                autotypeThread = Thread(target=do_autotype, args=(self.history_list.get_selected_row().username, self.history_list.get_selected_row().password,))
                autotypeThread.start()
                print(self.history_list.get_selected_row().get_title())
            if keyval == 112:
                print("copy password")
                clipboard.write(self.history_list.get_selected_row().password)
            elif keyval == 117:
                print("copy username")
                clipboard.write(self.history_list.get_selected_row().username)
                
        keycont.connect('key-pressed', handle_keypress, self)
        self.add_controller(keycont)

        self.history_list.get_style_context().add_class("boxed-list")
        self.box.append(self.history_list)
        self.set_default_size(700, 700)
        self.set_title("Goldwarden")

        def on_close(window):
            while self.history_list.get_first_child() != None:
                self.history_list.remove(self.history_list.get_first_child())
            window.hide()
            gc.collect()
            return True
        self.connect("close-request", on_close)

    def show(self):
        for i in range(0, 5):
            action_row = Adw.ActionRow()
            action_row.set_title("aaa")
            action_row.set_subtitle("Test")
            self.history_list.append(action_row)
        

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

        self.autofill_row = Adw.ActionRow()
        self.autofill_row.set_title("Autofill Shortcut")
        self.autofill_row.set_subtitle("Unavailable, please set up a shortcut in your desktop environment (README)")
        self.preferences_group.add(self.autofill_row)

        self.status_row = Adw.ActionRow()
        self.status_row.set_title("DBUS Service")
        self.status_row.set_subtitle("Listening")
        self.preferences_group.add(self.status_row)

        self.status_row = Adw.ActionRow()
        self.status_row.set_title("Vault Status")
        self.status_row.set_subtitle("Locked")
        self.preferences_group.add(self.status_row)

        self.login_row = Adw.ActionRow()
        self.login_row.set_title("Vault Login Entries")
        self.login_row.set_subtitle("0")
        self.preferences_group.add(self.login_row)

        self.notes_row = Adw.ActionRow()
        self.notes_row.set_title("Vault Notes")
        self.notes_row.set_subtitle("0")
        self.preferences_group.add(self.notes_row)

        self.login_button = Gtk.Button()
        self.login_button.set_label("Login")
        self.login_button.connect("clicked", lambda button: show_login())
        self.login_button.set_sensitive(False)
        self.login_button.set_margin_top(10)
        self.preferences_group.add(self.login_button)
    
        self.set_pin_button = Gtk.Button()
        self.set_pin_button.set_label("Set Pin")
        def set_pin():
            set_pin_thread = Thread(target=goldwarden.enable_pin)
            set_pin_thread.start()
        self.set_pin_button.connect("clicked", lambda button: set_pin())
        self.set_pin_button.set_margin_top(10)
        self.set_pin_button.set_sensitive(False)
        self.preferences_group.add(self.set_pin_button)

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
        self.preferences_group.add(self.unlock_button)

        self.logout_button = Gtk.Button()
        self.logout_button.set_label("Logout")
        self.logout_button.set_margin_top(10)
        self.logout_button.connect("clicked", lambda button: goldwarden.purge())
        self.preferences_group.add(self.logout_button)
        
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
        if hasattr(self, "win") is False:
            app.hold()
            self.win = MainWindow(application=app)
            self.win.set_hide_on_close(True)
            self.win.hide()
            self.settings_win = SettingsWinvdow(application=app)
            self.settings_win.set_hide_on_close(True)
        self.settings_win.present()

app = MyApp(application_id="com.quexten.Goldwarden")
def on_autofill():
    logins = goldwarden.get_vault_logins()
    if logins == None:
        return
    app.win.logins = logins
    app.win.show()
    app.win.present()
monitors.dbus_autofill_monitor.on_autofill = on_autofill

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

