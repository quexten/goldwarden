import os
import subprocess

is_flatpak = os.path.exists("/.flatpak-info")

def register_autostart(autostart: bool):
    if is_flatpak:
        print("Running in flatpak, registering with background portal for autostart.")
        try:
            subprocess.Popen(["python3", "/app/bin/src/linux/flatpak/autostart.py"], start_new_session=True)
        except:
            pass


def set_status(status: str):
    if is_flatpak:
        try:
            subprocess.Popen(["python3", "/app/bin/src/linux/flatpak/status.py", status], start_new_session=True)
        except:
            pass
