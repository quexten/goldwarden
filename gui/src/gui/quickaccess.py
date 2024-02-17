import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
gi.require_version('Notify', '0.7')
import gc
import time
from gi.repository import Gtk, Adw, GLib, Notify, Gdk
from ..services import goldwarden
from threading import Thread
import sys
import os
from ..services import totp
Notify.init("Goldwarden")

token = "Test"
goldwarden.create_authenticated_connection(token)

def autotype(text):
    time.sleep(2)
    print("Autotyping")
    goldwarden.autotype(text)
    print("Autotyped")
    time.sleep(5)
    os._exit(0)

def set_clipboard(text):
    Gdk.Display.get_clipboard(Gdk.Display.get_default()).set_content(
            Gdk.ContentProvider.new_for_value(text)
        )

    def kill():
        time.sleep(0.5)
        os._exit(0)
    thread = Thread(target=kill)
    thread.start()

class MyApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def update_logins(self):
        logins = goldwarden.get_vault_logins()
        if logins == None:
            os._exit(0)
            return
        self.app.autofill_window.logins = logins

    def on_activate(self, app):
        self.autofill_window = MainWindow(application=app)
        self.autofill_window.logins = []
        self.autofill_window.present()
        self.app = app
        thread = Thread(target=self.update_logins)
        thread.start()

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
    
        def on_type(entry):
            if len(entry.get_text()) > 1:
                self.results_list.show()
            else:
                self.results_list.hide()

            while self.results_list.get_first_child() != None:
                self.results_list.remove(self.results_list.get_first_child())

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
                action_row.uuid = i["uuid"]
                action_row.uri = i["uri"]
                action_row.totp = i["totp"]
                self.results_list.append(action_row)
            self.starts_with_logins = None
            self.other_logins = None
        self.text_view.connect("changed", lambda entry: on_type(entry))
        self.box.append(self.text_view)
    
        self.results_list = Gtk.ListBox()
        # margin'
        self.results_list.set_margin_start(10)
        self.results_list.set_margin_end(10)
        self.results_list.set_margin_top(10)
        self.results_list.set_margin_bottom(10)
        self.results_list.hide()

        keycont = Gtk.EventControllerKey()
        def handle_keypress(cotroller, keyval, keycode, state, user_data):
            ctrl_pressed = state & Gdk.ModifierType.CONTROL_MASK > 0
            alt_pressed = state & Gdk.ModifierType.ALT_MASK > 0

            if keycode == 9:
                os._exit(0)

            if keyval == 65364:
                # focus results
                if self.results_list.get_first_child() != None:
                    self.results_list.get_first_child().grab_focus()
                    self.results_list.select_row(self.results_list.get_first_child())

            if keyval == 113:
                return False

            if keycode == 36:
                self.hide()
                autotypeThread = Thread(target=autotype, args=(f"{self.results_list.get_selected_row().username}\t{self.results_list.get_selected_row().password}",))
                autotypeThread.start()
            if keyval == 112:
                print("pass", ctrl_pressed, alt_pressed)
                if ctrl_pressed and not alt_pressed:
                    set_clipboard(self.results_list.get_selected_row().password)
                if ctrl_pressed and alt_pressed:
                    self.hide()
                    autotype(self.results_list.get_selected_row().password)
            elif keyval == 117:
                if ctrl_pressed and not alt_pressed:
                    set_clipboard(self.results_list.get_selected_row().username)
                if ctrl_pressed and alt_pressed:
                    self.hide()
                    autotype(self.results_list.get_selected_row().username)
            elif keyval == 118:
                if ctrl_pressed and alt_pressed:
                    environment = goldwarden.get_environment()
                    if environment == None:
                        return
                    item_uri = environment["vault"] + "#/vault?itemId=" + self.results_list.get_selected_row().uuid
                    Gtk.show_uri(None, item_uri, Gdk.CURRENT_TIME)
            elif keyval == 108:
                if ctrl_pressed and alt_pressed:
                    Gtk.show_uri(None, self.results_list.get_selected_row().uri, Gdk.CURRENT_TIME)
            elif keyval == 116:
                totp_code = totp.totp(self.resuts_list.get_selected_row().totp)
                if ctrl_pressed and not alt_pressed:
                    set_clipboard(totp_code)
                if ctrl_pressed and alt_pressed:
                    self.hide()
                    autotype(totp_code)
            elif keyval == 102:
                # focus search
                self.text_view.grab_focus()
                
        keycont.connect('key-pressed', handle_keypress, self)
        self.add_controller(keycont)

        self.results_list.get_style_context().add_class("boxed-list")
        self.box.append(self.results_list)
        self.set_default_size(700, 700)
        self.set_title("Goldwarden Quick Access")

app = MyApp(application_id="com.quexten.Goldwarden.autofill-menu")
app.run(sys.argv)
