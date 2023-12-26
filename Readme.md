## Goldwarden

Goldwarden is a Bitwarden compatible desktop integration written in Go. It focuses on providing useful desktop features that the official tools 
do not (yet) have or are not willing to add, and enhanced security measures that other tools do not provide, such as:

- Support for SSH Agent (Git signing and SSH login)
- Support for injecting environment variables into the environment of a cli command
- System wide autotype
- Biometric authentication (via Polkit) for each credential access
- Vault content is held encrypted in memory and only briefly decrypted when needed
- Kernel level memory protection for keys (via the memguard library)
- Additional measures to protect against memory dumps
- Passwordless login (Both logging in, and approving logins)
- Fido2 (Webauthn) support
- more to come...?

The current goal is not to provide a full featured Bitwarden CLI, but to provide specific features that are not available in other tools.
If you want an officially supported way to manage your Bitwarden vault, you should use the Bitwarden CLI (or a regular client).
If you are looking to manage secrets for machine to machine communication, you should use bitwarden secret manager or something like 
hashicorp vault.

### Requirements
Right now, Goldwarden is only tested on Linux. It should be possible to port to mac / bsd, I'm open to PRs.
On Linux, you need at least a working Polkit installation and a pinentry agent. Both X11 and Wayland are supported for autofill, albeit only Wayland is tested.

### Installation

#### Flatpak (WIP)
There is a flatpak that includes a small UI, autotype functionality and autostarting of the daemon.
**Not yet on flathub**

#### CLI
On Arch linux, or other distributions with access to the AUR, simply:
```
yay -S goldwarden
```
should be enough to install goldwarden on your system.

For deb/rpm, download the deb/rpm from the latest release on GitHub and install it using your package manager.

On other distributions, Mac and Windows, you can download it from the latest release on GitHub and put it into a location you want to have it in, f.e `/usr/bin`.
Then run `goldwarden setup polkit`.
Optionally run: `goldwarden setup systemd` and `goldwarden setup browserbiometrics`.

Alternatively, you can build it yourself.
```
go install github.com/quexten/goldwarden@latest
```
### Usage

Start the daemon (this is done by systemd automatically, when set up with `goldwarden setup systemd`):
```
goldwarden daemon
```

Set a pin
```
goldwarden set pin
```

Optionally set the api/identity url for a custom bitwarden server:
```
goldwarden config set-api-url https://my.bitwarden.domain/api
```

```
goldwarden config set-identity-url https://my.bitwarden.domain/identity
```

```
goldwarden config set-notifications-url https://my.bitwarden.domain/notifications
```

Login
```
goldwarden login --email <email>
```

##### Login-with-device
```
goldwarden login --email <email> --passwordless
```

Create an ssh key
```
goldwarden ssh add --name <name>
```

Run a command with injected environment variables
```
goldwarden run -- <command>
```

#### Autofill (Flatpak Only)

To set up a shortcut (CTRL+u) on Gnome:
```
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/gwautofill/']"
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/gwautofill/ name 'Goldwarden Autofill'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/gwautofill/ command 'dbus-send --type=method_call --dest=com.quexten.Goldwarden.autofill /com/quexten/Goldwarden com.quexten.Goldwarden.Autofill.autofill'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/gwautofill/ binding '<Ctrl>u'
```

On other Desktop environments this will work differently, just makes sure that:
`dbus-send --type=method_call --dest=com.quexten.Goldwarden.autofill /com/quexten/Goldwarden com.quexten.Goldwarden.Autofill.autofill`
gets called.

This will be changed once desktop environments implement the global hotkey portal.

#### SSH Agent
[goldwarden_ssh.webm](https://github.com/quexten/goldwarden/assets/11866552/9058f734-60e0-4dd3-b9f8-1d77f7cf4c65)


The SSH agent listens on a socket on `~/.goldwarden-ssh-agent.sock`. This can be used f.e by doing:
```
SSH_AUTH_SOCK=~/.goldwarden-ssh-agent.sock ssh-add -l
```

Beware that some applications do ssh requests in the background, so you might need to set the env variable in your shell config.

To add a key to your vault, run:
```
goldwardens ssh add --name "my key"
```

Alternatively, use one of the gui clients. Create an ed25519 key:
```
ssh-keygen -t ed25519 -f ./id_ed25519
```
Then create a secure note in bitwarden:
```
custom-type: ssh-key
private-key: <contents of id_ed25519> (hidden field)
public-key: <contents of id_ed25519.pub>
```

Then add the private key to bitwarden. The public key can be added to your github account f.e.

##### Git Signing
[goldwarden_git.webm](https://github.com/quexten/goldwarden/assets/11866552/f47dcd93-789d-4fcc-954b-d43d9033e213)

To use the SSH agent for git signing, you need to add the following to your git config:
```
[user]
        email = <your email>
        name = <your name>
        signingKey = <your public key>
[commit]
        gpgsign = true
[gpg]
        format = ssh
```

### Environment Variables
[goldwarden_run_restic.webm](https://github.com/quexten/goldwarden/assets/11866552/9a342df5-feec-4174-a0e9-6a399c2feb65)

Goldwarden can inject environment variables into the environment of a cli command.

First, create a secure note in bitwarden, and add the following custom fields (using restic as an example):
```
custom-type: env
executable: name_of_executable
# env variables
AWS_ACCESS_KEY_ID: <your access key>
AWS_SECRET_ACCESS_KEY: <your secret key> (hidden)
RESTIC_PASSWORD: <your restic password> (hidden)
# optional
RESTIC_REPOSITORY: <your restic repository>
...
```

Then, run the command:
```
goldwarden run -- <command>
```
I.e
```
goldwarden run -- restic backup
```

You can also alias the commands, such that every time you run them, the environment variables are injected:
```
alias restic="goldwarden run -- restic"
```

And then just run the command as usual:
```
restic backup
```

### Autotype based Autofill (Flatpak Only)
[goldwarden_autofill.webm](https://github.com/quexten/goldwarden/assets/11866552/6ac7cdc2-0cd7-42fd-9fd0-cfff26e2ceee)

You can bind this to a hotkey in your desktop environment (i.e i3/sway config file, Gnome custom shortcuts, etc).
```
dbus-send --type=method_call --dest=com.quexten.goldwarden /com/quexten/goldwarden com.quexten.goldwarden.Autofill.autofill
```
#### XDG-RemoteDesktop-Portal

By default, the remote desktop portal is used. As long as your desktop environment handle this (KDE and Gnome do, wlroots does not yet)
this enables autotyping without having to modify permissions.

### Login with device
Approving other devices works out of the box and is enabled by default. If the agent is unlocked, you will be prompted
to approve the device. If you want to log into goldwarden using another device, add the "--passwordless" parameter to the login command.


### Environment Variables
```
GOLDWARDEN_WEBSOCKET_DISABLED=true # disable websocket
GOLDWARDEN_PIN_REQUIREMENT_DISABLED=true # disable pin requirement
GOLDWARDEN_DO_NOT_PERSIST_CONFIG=true # do not persist config
GOLDWARDEN_API_URI=https://my.bitwarden.domain/api # set api uri
GOLDWARDEN_IDENTITY_URI=https://my.bitwarden.domain/identity # set identity uri
GOLDWARDEN_SINGLE_PROCESS=true # run in single process mode, i.e no daemon
GOLDWARDEN_DEVICE_UUID= # set device uuid (beware, this is used as a salt for the configuration encryption)
GOLDWARDEN_AUTH_METHOD= # set auth method (password/passwordless)
GOLDWARDEN_AUTH_USER= # set auth user
GOLDWARDEN_AUTH_PASSWORD= # set auth password
GOLDWARDEN_SILENT_LOGGING=true # disable logging
GOLDWARDEN_SYSTEM_AUTH_DISABLED=true # disable system auth (biometrics / approval)
```

### Building

To build, you will need libfido2-dev. And a go toolchain. 

Run:
```
go install github.com/quexten/goldwarden@latest
go install -tags autofill github.com/quexten/goldwarden@latest
```

or:
```
go build
```

### Design
The tool is split into CLI and daemon, which communicate via a unix socket.

The vault is never written to disk and is only kept in encrypted form in memory, it is re-downloaded upon startup. The encryption keys are stored in secure enclaves (using the memguard library) and only decrypted briefly when needed. This protects from memory dumps. Vault entries are also only decrypted when needed.

When entries change, the daemon gets notified via websockets and updates automatically.

The sensitive parts of the config file are encrypted using a pin. The key is derrived using argon2, and the encryption used is chacha20poly1305. The config is also only held in memory in encrypted form and decrypted using key stored in kernel secured memory when needed.

When accessing a vault entry, the daemon will authenticate against a polkit policy. This allows using biometrics.

By default, credential entry is cached for 10 minutes. During this time, a parent program can invoke goldwarden multiple times, but biometrics are only confirmed the first time. Since this is per parent-program, this means that invokations from 2 tty's would independently each ask for biometrics confirmation the first time.

### Future Plans
Some things that I consider adding (depending on time and personal need):
- Regular cli managment (add, delete, update, of logins / secure notes)
- Installers
- (MacOS & Windows support tracked in https://github.com/quexten/goldwarden/issues/4)

If you have other interesting ideas, feel free to open an issue. I can't
promise that I will implement it, but I'm open to suggestions.

### Unsuported
Some things that are unsupported and not likely to develop myself:
- Send support
- Attachments
- Credit cards / Identities
