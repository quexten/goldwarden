#!/usr/bin/env python3
import sys
import gi

gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')

from gi.repository import Gtk, Adw, GLib, Gdk, Gio
from ..services import goldwarden
from threading import Thread
from .resource_loader import load_template, load_json
import subprocess
from . import components
import os

class GoldwardenSSHSetupGuideApp(Adw.Application):
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.connect('activate', self.on_activate)

    def on_activate(self, app):
        self.load()
        self.window.present()

    def load(self):
        builder = load_template("ssh.ui")
        self.window = builder.get_object("window")
        self.window.set_application(self)
        commands = load_json("commands")
        self.add_ssh_key_row = builder.get_object("add_ssh_key_row")
        self.add_ssh_key_row.set_subtitle(commands["add-ssh-key"])
        self.ssh_socket_path_row = builder.get_object("ssh_socket_path_row")
        self.ssh_socket_path_row.set_subtitle(commands["ssh-socket-path"])

if __name__ == "__main__":
    app = GoldwardenSSHSetupGuideApp(application_id="com.quexten.Goldwarden.sshsetup")
    app.run(sys.argv)