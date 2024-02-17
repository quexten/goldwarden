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

        self.add_ssh_key_group = Adw.PreferencesGroup()
        self.add_ssh_key_group.set_title("Add an SSH Key")
        self.preferences_page.add(self.add_ssh_key_group)

        self.add_ssh_key_row = Adw.ActionRow()
        self.add_ssh_key_row.set_subtitle("flatpak run --command=goldwarden com.quexten.Goldwarden ssh add --name MY_KEY_NAME")
        self.add_ssh_key_row.set_subtitle_selectable(True)
        self.add_ssh_key_group.add(self.add_ssh_key_row)

        self.ssh_socket_path_group = Adw.PreferencesGroup()
        self.ssh_socket_path_group.set_title("SSH Socket Path")
        self.ssh_socket_path_group.set_description("Add this to your your enviorment variables")
        self.preferences_page.add(self.ssh_socket_path_group)

        self.ssh_socket_path_row = Adw.ActionRow()
        self.ssh_socket_path_row.set_subtitle("export SSH_AUTH_SOCK=/home/$USER/.var/app/com.quexten.Goldwarden/data/ssh-auth-sock")
        self.ssh_socket_path_row.set_subtitle_selectable(True)
        self.ssh_socket_path_group.add(self.ssh_socket_path_row)

        self.git_signing_group = Adw.PreferencesGroup()
        self.git_signing_group.set_title("Git Signing")
        self.git_signing_group.set_description("Check the wiki for more information")
        self.preferences_page.add(self.git_signing_group)

        self.set_default_size(400, 700)
        self.set_title("Goldwarden SSH Setup")

app = MyApp(application_id="com.quexten.Goldwarden.sshsetup")
app.run(sys.argv)