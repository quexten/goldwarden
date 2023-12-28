from gi.repository import Gtk

def status_icon_ok(icon_name):
    imagebox = Gtk.Box()
    imagebox.set_orientation(Gtk.Orientation.VERTICAL)
    imagebox.set_halign(Gtk.Align.CENTER)
    imagebox.set_valign(Gtk.Align.CENTER)
    image = Gtk.Image()
    image.get_style_context().add_class("status-icon")
    image.get_style_context().add_class("ok-icon")
    image.set_from_icon_name(icon_name)
    imagebox.append(image)
    return imagebox

def status_icon_error(icon_name):
    imagebox = Gtk.Box()
    imagebox.set_orientation(Gtk.Orientation.VERTICAL)
    imagebox.set_halign(Gtk.Align.CENTER)
    imagebox.set_valign(Gtk.Align.CENTER)
    image = Gtk.Image()
    image.get_style_context().add_class("status-icon")
    image.get_style_context().add_class("error-icon")
    image.set_from_icon_name(icon_name)
    imagebox.append(image)
    return imagebox

def status_icon_warning(icon_name):
    imagebox = Gtk.Box()
    imagebox.set_orientation(Gtk.Orientation.VERTICAL)
    imagebox.set_halign(Gtk.Align.CENTER)
    imagebox.set_valign(Gtk.Align.CENTER)
    image = Gtk.Image()
    image.get_style_context().add_class("status-icon")
    image.get_style_context().add_class("warning-icon")
    image.set_from_icon_name(icon_name)
    imagebox.append(image)

    return imagebox

class StatusIcon(Gtk.Box):
    def __init__(self):
        super().__init__()
        self.icon_name = None
        self.status = None

    def set_icon(self, icon_name, status):
        if self.icon_name == icon_name and self.status == status:
            return
        self.icon_name = icon_name
        self.status = status
        
        while self.get_first_child() != None:
            self.remove(self.get_first_child())
        
        if status == "ok":
            self.append(status_icon_ok(icon_name))
        elif status == "error":
            self.append(status_icon_error(icon_name))
        elif status == "warning":
            self.append(status_icon_warning(icon_name))
        else:
            raise Exception("Invalid status", status)