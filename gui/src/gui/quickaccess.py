import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
# gi.require_version('Notify', '0.7')
import gc
import time
from gi.repository import Gtk, Adw, GLib, Gdk
from ..services import goldwarden
from ..services.autotype import autotype
from threading import Thread
from .resource_loader import load_template
import sys
import os
from ..services import totp

class GoldwardenQuickAccessApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.logins = []
        self.filtered_logins = []
        self.query = ""
        self.connect('activate', self.on_activate)
        self.selected_index = 0

    def on_activate(self, app):
        self.load()
        self.window.present()
        thread = Thread(target=self.update_logins)
        thread.start()

    def load(self):
        builder = load_template("quickaccess.ui")
        self.window = builder.get_object("window")
        self.results_list = builder.get_object("results_list")
        self.status_page = builder.get_object("status_page")
        self.text_view = builder.get_object("search_row")
        self.text_view.connect("changed", self.on_type)
        self.window.set_application(self)

        evk = Gtk.EventControllerKey.new()
        evk.set_propagation_phase(Gtk.PropagationPhase.CAPTURE)
        evk.connect("key-pressed", self.key_press)
        self.window.add_controller(evk)  

    def key_press(self, event, keyval, keycode, state):
        if keyval == Gdk.KEY_Escape:
            os._exit(0)

        if keyval == Gdk.KEY_Tab:
            return True

        if keyval == Gdk.KEY_Up:
            self.selected_index = self.selected_index - 1
            if self.selected_index < 0:
                self.selected_index = 0
            self.render_list()
            return True
        elif keyval == Gdk.KEY_Down:
            self.selected_index = self.selected_index + 1
            if self.selected_index >= len(self.filtered_logins):
                self.selected_index = len(self.filtered_logins) - 1
            self.render_list()
            return True

        if self.selected_index >= len(self.filtered_logins) or self.selected_index < 0:
            self.selected_index = 0

        auto_type_combo = state & Gdk.ModifierType.CONTROL_MASK and state & Gdk.ModifierType.SHIFT_MASK
        copy_combo = state & Gdk.ModifierType.CONTROL_MASK and not state & Gdk.ModifierType.SHIFT_MASK

        if not len(self.filtered_logins) > 0:
            return

        # totp code
        if keyval == Gdk.KEY_t or keyval == Gdk.KEY_T:
            if self.filtered_logins[self.selected_index]["totp"] == "":
                return
            if auto_type_combo:
                self.run_autotype(totp.totp(self.filtered_logins[self.selected_index]["totp"]))
            if copy_combo:
                self.set_clipboard(totp.totp(self.filtered_logins[self.selected_index]["totp"]))

        if keyval == Gdk.KEY_u or keyval == Gdk.KEY_U:
            if auto_type_combo:
                self.run_autotype(self.filtered_logins[self.selected_index]["username"])
            if copy_combo:
                self.set_clipboard(self.filtered_logins[self.selected_index]["username"])
        
        if keyval == Gdk.KEY_p or keyval == Gdk.KEY_P:
            if auto_type_combo:
                self.run_autotype(self.filtered_logins[self.selected_index]["password"])
            if copy_combo:
                self.set_clipboard(self.filtered_logins[self.selected_index]["password"])

        if (keyval == Gdk.KEY_l or keyval == Gdk.KEY_L) and auto_type_combo:
            Gtk.show_uri(None, self.results_list.get_selected_row().uri, Gdk.CURRENT_TIME)

        if (keyval == Gdk.KEY_v or keyval == Gdk.KEY_V) and auto_type_combo:
            self.set_clipboard(self.filtered_logins[self.selected_index]["uri"])
            environment = goldwarden.get_environment()
            if environment == None:
                return
            item_uri = environment["vault"] + "#/vault?itemId=" + self.results_list.get_selected_row().uuid
            Gtk.show_uri(None, item_uri, Gdk.CURRENT_TIME)

        if keyval == Gdk.KEY_Return:
            if auto_type_combo:
                self.run_autotype(f"{self.filtered_logins[self.selected_index]['username']}\t{self.filtered_logins[self.selected_index]['password']}")

    def update(self):
        self.update_list()
        self.render_list()
        return True

    def run_autotype(self, text):
        def perform_autotype(text):
            GLib.idle_add(self.window.hide)
            time.sleep(2)
            autotype.autotype(text)
            time.sleep(0.1)
            os._exit(0)
        thread = Thread(target=perform_autotype, args=(text,))
        thread.start()

    def set_clipboard(self, text):
        Gdk.Display.get_clipboard(Gdk.Display.get_default()).set_content(
            Gdk.ContentProvider.new_for_value(text)
        )

        def kill():
            time.sleep(0.5)
            os._exit(0)
        thread = Thread(target=kill)
        thread.start()

    def update_list(self):
        if self.query == "":
            self.filtered_logins = []
            return
        
        self.filtered_logins = list(filter(lambda i: self.query.lower() in i["name"].lower(), self.logins))

        self.starts_with_logins = list(filter(lambda i: i["name"].lower().startswith(self.query.lower()), self.filtered_logins))
        self.other_logins = list(filter(lambda i: i not in self.starts_with_logins, self.filtered_logins))
        self.filtered_logins = self.starts_with_logins + self.other_logins
        if len(self.filtered_logins) > 7:
            self.filtered_logins = self.filtered_logins[0:7]

    def render_list(self):
        if len(self.filtered_logins) > 0:
            self.results_list.set_visible(True)
            while self.results_list.get_first_child() != None:
                self.results_list.remove(self.results_list.get_first_child())
            self.status_page.set_visible(False)
        else:
            self.results_list.set_visible(False)
            self.status_page.set_visible(True)

        for i in self.filtered_logins:
            action_row = Adw.ActionRow()
            if "name" in i:
                action_row.set_title(i["name"])
            else:
                action_row.set_title("[no name]")
            if "username" in i:
                action_row.set_subtitle(i["username"])
                action_row.username = i["username"]
            else:
                action_row.set_subtitle("[no username]")
                action_row.username = "[no username]"
            if "password" in i:
                action_row.password = i["password"]
            else:
                action_row.password = "[no password]"
            if "uri" in i:
                action_row.uri = i["uri"]
            else:
                action_row.uri = "[no uri]"
            if "uuid" in i:
                action_row.uuid = i["uuid"]
            else:
                action_row.uuid = "[no uuid]"
            if "totp" in i:
                action_row.totp = i["totp"]
            else:
                action_row.totp = ""
            action_row.set_icon_name("dialog-password")
            action_row.set_activatable(True)
            self.results_list.append(action_row)
        
        # select the nth item
        if len(self.filtered_logins) > 0:
            self.results_list.select_row(self.results_list.get_row_at_index(self.selected_index))
            self.results_list.set_focus_child(self.results_list.get_row_at_index(self.selected_index))
        
        self.starts_with_logins = None
        self.other_logins = None

    def on_type(self, entry):
        search_query = entry.get_text()
        self.query = search_query
        self.update()
    
    def update_logins(self):
        logins = goldwarden.get_vault_logins()
        if logins == None:
            os._exit(0)
            return
        self.logins = logins
        self.update()
    
if __name__ == "__main__":
    settings = Gtk.Settings.get_default()
    settings.set_property("gtk-error-bell", False)

    token = sys.stdin.readline()
    goldwarden.create_authenticated_connection(token)
    app = GoldwardenQuickAccessApp(application_id="com.quexten.Goldwarden.quickaccess")
    app.run(sys.argv)