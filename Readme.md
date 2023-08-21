## Goldwarden

Goldwarden is a Bitwarden compatible CLI tool written in Go. It focuses on features for Desktop integration, and enhanced security measures that other tools do not provide, such as:

- Support for SSH Agent (Git signing and SSH login)
- Support for injecting environment variables into the environment of a cli command
- System wide autofill
- Biometric authentication (via Polkit) for each credential access
- Vault content is held encrypted in memory and only briefly decrypted when needed
- Kernel level memory protection for keys (via the memguard library)
- Additional measures to protect against memory dumps
- Passwordless login (Approval of other login)
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

On Arch linux, or other distributions with access to the AUR, simply:
```
yay -S goldwarden
```
should be enough to install goldwarden on your system.

Otherwise:
Download the latest release binary and put it into a location you want to have it in, f.e `/usr/bin`.
Then run `goldwarden setup polkit`.
Optionally run: `goldwarden setup systemd` and `goldwarden setup browserbiometrics`.

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

Login
```
goldwarden login --email <email>
```

Create an ssh key
```
goldwarden ssh add --name <name>
```

Run a command with injected environment variables
```
goldwarden run -- <command>
```

Autofill
```
goldwarden autofill --layout <keyboard-layout>
```
(Create a hotkey for this depending on your desktop environment)

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

### Autofill
[goldwarden_autofill.webm](https://github.com/quexten/goldwarden/assets/11866552/6ac7cdc2-0cd7-42fd-9fd0-cfff26e2ceee)

The autofill feature is a bit experimental. It autotypes the password via uinput. This needs a keyboardlayout to map the letters to 
keycodes. Currently supported are qwerty and dvorak.
`goldwarden autofill --layout qwerty`
`goldwarden autofill --layout dvorak`

You can bind this to a hotkey in your desktop environment (i.e i3/sway config file, Gnome custom shortcuts, etc).

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
GOLDWARDEN_DEVICE_UUID= # set device uuid
GOLDWARDEN_AUTH_METHOD= # set auth method (password/passwordless)
GOLDWARDEN_AUTH_USER= # set auth user
GOLDWARDEN_AUTH_PASSWORD= # set auth password
GOLDWARDEN_SILENT_LOGGING=true # disable logging
GOLDWARDEN_SYSTEM_AUTH_DISABLED=true # disable system auth (biometrics / approval)
```

### Building

To build, you will need libfido2-dev. And a go toolchain. 

Additionally, if you want the autofill feature you will need some dependencies. Everything from https://gioui.org/doc/install linux and wl-clipboard (or xclipboard) should be installed.

Run:
```
go install github.com/quexten/goldwarden@latest
go install -tags autofill github.com/quexten/goldwarden@latest
```

or:
```
go build
go build -tags autofill
```

### Design
The tool is split into CLI and daemon, which communicate via a unix socket.

The vault is never written to disk and is only kept in encrypted form in memory, it is re-downloaded upon startup. The encryption keys are stored in secure enclaves (using the memguard library) and only decrypted briefly when needed. This protects from memory dumps. Vault entries are also only decrypted when needed.

When entries change, the daemon gets notified via websockets and updates automatically.

The sensitive parts of the config file are encrypted using a pin. The key is derrived using argon2, and the encryption used is chacha20poly1305. The config is also only held in memory in encrypted form and decrypted using key stored in kernel secured memory when needed.

When accessing a vault entry, the daemon will authenticate against a polkit policy. This allows using biometrics.

### Future Plans
Some things that I consider adding (depending on time and personal need):
- Regular cli managment (add, delete, update, of logins / secure notes)
- Scripts to properly set up the policies/systemd/etc.

If you have other interesting ideas, feel free to open an issue. I can't
promise that I will implement it, but I'm open to suggestions.

### Unsuported
Some things that are unsupported and not likely to develop myself:
- MacOS / BSD support (should not be too much work, most things should work out of the box, some adjustments for pinentry and polkit would be needed)
- Windows support (probably a lot of work, unix sockets don't really exist, and pinentry / polkit would have to be implemented otherwise. There might be go libraries for that, but I don't know)
- Send support
- Attachments
- Credit cards / Identities
