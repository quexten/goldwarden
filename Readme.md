## Goldwarden

Goldwarden is a Bitwarden compatible desktop client. It focuses on providing useful desktop features that the official tools 
do not (yet) have or are not willing to add, and enhanced security measures that other tools do not provide, such as:

- Support for SSH Agent (Git signing and SSH login)
- System wide autotype (Gnome, KDE only for now)
- Biometric authentication (via Polkit) for each credential access
- Implements Bitwarden browser-extension biometrics on Linux
- Support for injecting environment variables into the environment of a cli command
- Vault content is held encrypted in memory and only briefly decrypted when needed
- Kernel level memory protection for keys (via the memguard library)
- Additional measures to protect against memory dumps
- Passwordless login (Both logging in, and approving logins)
- Fido2 (Webauthn) support
- more to come...?

The aim is not to replace the official clients, but to complement by implementing the missing features.

### Requirements
Right now, Goldwarden is only tested on Linux. Somewhat feature-stripped builds for Mac and Windows are available too, but untested.
Autotype is currently implemented via the remotedesktop portal. This is supported on KDE and Gnome, but not yet on wl-root based environments.

### Installation

#### Flatpak
There is a flatpak that includes a small UI, autotype functionality and autostarting of the daemon.
[<img width='240' alt='Download on Flathub' src='https://flathub.org/assets/badges/flathub-badge-en.png' />](https://flathub.org/apps/details/com.quexten.Goldwarden)

<img src='https://github.com/quexten/goldwarden/assets/11866552/5d36ed8c-46f1-4444-adb0-f4ca1d0433c5' width='700'>

#### CLI
##### Arch (AUR)
On Arch linux, or other distributions with access to the AUR, simply:
```
yay -S goldwarden
```
should be enough to install goldwarden on your system.

##### Deb / RPM
For deb/rpm, download the deb/rpm from the latest release on GitHub and install it using your package manager.

#### NixOS
https://github.com/NixOS/nixpkgs/pull/278362

##### Github Binary Releases
On other distributions, Mac and Windows, you can download it from the latest release on GitHub and put it into a location you want to have it in, f.e `/usr/bin`.

##### Compiling
Alternatively, you can build it yourself.
```
go install github.com/quexten/goldwarden@latest
```

### Setup and Usage
To get started, follow the instructions provided in the wiki https://github.com/quexten/goldwarden/wiki/Getting-Started.
For instructions on specific features, also consult the wiki page for the feature.
