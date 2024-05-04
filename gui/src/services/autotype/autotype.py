import sys
import os

is_linux = sys.platform == 'linux'
is_wayland = os.environ.get('XDG_SESSION_TYPE') == 'wayland'

def autotype(text):
    if is_linux and is_wayland:
        from .libportal_autotype import autotype_libportal
        autotype_libportal(text)
    elif is_linux:
        from .x11autotype import type
        type(text)
    else:
        from .pyautogui_autotype import autotype_pyautogui
        autotype_pyautogui(text)