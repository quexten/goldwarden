import gi
gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')
gi.require_version('Notify', '0.7')
import gc
import time
from gi.repository import Gtk, Adw, GLib, Notify, Gdk
from ..services import goldwarden
from threading import Thread
from .template_loader import load_template
import sys
import os
from ..services import totp
Notify.init("Goldwarden")

class GoldwardenQuickAccessApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.logins = []
        self.filtered_logins = []
        self.query = ""
        self.connect('activate', self.on_activate)

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

    def update(self):
        self.update_list()
        self.render_list()

    def autotype(self, text):
        goldwarden.autotype(text)
        time.sleep(0.1)
        os._exit(0)

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
        if len(self.filtered_logins) > 1:
            self.results_list.set_visible(True)
            while self.results_list.get_first_child() != None:
                self.results_list.remove(self.results_list.get_first_child())
            self.status_page.set_visible(False)
        else:
            self.results_list.set_visible(False)
            self.status_page.set_visible(True)

        for i in self.filtered_logins:
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

    def on_type(self, entry):
        search_query = entry.get_text()
        self.query = search_query
        self.update()
    
    def update_logins(self):
        logins = goldwarden.get_vault_logins()
        print(logins)
        if logins == None:
            os._exit(0)
            return
        self.logins = logins
        self.update()
    
if __name__ == "__main__":
    # todo add proper method to debug this
    # token = sys.stdin.readline()
    token = "Test"
    goldwarden.create_authenticated_connection(token)
    app = GoldwardenQuickAccessApp(application_id="com.quexten.Goldwarden.quickaccess")
    app.run(sys.argv)