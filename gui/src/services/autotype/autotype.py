import sys
import os

is_linux = sys.platform == 'linux'
is_wayland = os.environ.get('XDG_SESSION_TYPE') == 'wayland'

def autotype(text):
    print("autotypeing, is_linux: {}, is_wayland: {}".format(is_linux, is_wayland))
    if is_linux and is_wayland:
        from .libportal_autotype import autotype_libportal
        autotype_libportal(text)

    from .pyautogui_autotype import autotype_pyautogui
    autotype_pyautogui(text)