# TODO??!?!? for now using golang implementation
from ..goldwarden import autotype

def libportal_autotype(text):
    print("autotypeing with libportal")
    goldwarden.autotype(text)

# import dbus
# import dbus.mainloop.glib
# from dbus.mainloop.glib import DBusGMainLoop

# from gi.repository import GLib

# import random
# import time

# step = 0

# def typestring(text):
#     step = 0
#     handle = ""

#     def handler(*args, **kwargs):
#       global step
#       if step == 0:
#         handle_xdp_session_created(*args, **kwargs)
#       elif step == 1:
#         handle_xdp_devices_selected(*args, **kwargs)
#       elif step == 2:
#         handle_session_start(*args, **kwargs)
#       else:
#         print(args, kwargs)
#       step += 1

#     def handle_session_start(code, results, object_path):
#       global handle

#       if code != 0:
#         raise Exception("Could not start session")
      
#       for sym in text:
#         if sym == "\t":
#             inter.NotifyKeyboardKeycode(handle, {}, 15, 1)
#             time.sleep(0.001)
#             inter.NotifyKeyboardKeycode(handle, {}, 15, 0)
#             time.sleep(0.001)
#         else:
#             inter.NotifyKeyboardKeysym(handle, {}, ord(sym), 1)
#             time.sleep(0.001)
#             inter.NotifyKeyboardKeysym(handle, {}, ord(sym), 0)
#             time.sleep(0.001)

#       bus

#     def handle_xdp_devices_selected(code, results, object_path):
#       global handle

#       if code != 0:
#         raise Exception("Could not select devices")
      
#       start_options = {
#           "handle_token": "krfb" + str(random.randint(0, 999999999))
#       }
#       reply = inter.Start(handle, "", start_options)
#       print(reply)  

#     def handle_xdp_session_created(code, results, object_path):
#       global handle

#       if code != 0:
#         raise Exception("Could not create session")
#       print(results)
#       handle = results["session_handle"]

#       # select sources for the session
#       selection_options = {
#           "types": dbus.UInt32(7),  # request all (KeyBoard, Pointer, TouchScreen)
#           "handle_token": "krfb" + str(random.randint(0, 999999999))
#       }
#       selector_reply = inter.SelectDevices(handle, selection_options)
#       print(selector_reply)

#     def main():
#       global bus
#       global inter
#       loop = GLib.MainLoop()
#       dbus.mainloop.glib.DBusGMainLoop(set_as_default=True)
#       bus = dbus.SessionBus()
#       obj = bus.get_object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
#       inter = dbus.Interface(obj, "org.freedesktop.portal.RemoteDesktop")

#       bus.add_signal_receiver(
#         handler,
#         signal_name="Response",
#         dbus_interface="org.freedesktop.portal.Request",
#         bus_name="org.freedesktop.portal.Desktop",
#         path_keyword="object_path")

#       print(inter)
#       result = inter.CreateSession({
#         "session_handle_token": "sessionhandletoken"
#       })
#       print(result)
#       loop.run()

#     main()
