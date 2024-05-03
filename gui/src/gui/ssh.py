#!/usr/bin/env python3
import sys
import gi

gi.require_version('Gtk', '4.0')
gi.require_version('Adw', '1')

from gi.repository import Gtk, Adw, GLib, Gdk, Gio
from ..services import goldwarden
from threading import Thread
from .template_loader import load_template
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

if __name__ == "__main__":
    app = GoldwardenSSHSetupGuideApp(application_id="com.quexten.Goldwarden.sshsetup")
    app.run(sys.argv)