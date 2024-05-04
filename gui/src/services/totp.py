# https://github.com/susam/mintotp
#!/usr/bin/env python3

import base64
import hmac
import struct
import sys
import time

def hotp(key, counter, digits=6, digest='sha1'):
    key = base64.b32decode(key.upper() + '=' * ((8 - len(key)) % 8))
    counter = struct.pack('>Q', counter)
    mac = hmac.new(key, counter, digest).digest()
    offset = mac[-1] & 0x0f
    binary = struct.unpack('>L', mac[offset:offset+4])[0] & 0x7fffffff
    return str(binary)[-digits:].zfill(digits)


def totp(key, time_step=30, digits=6, digest='sha1'):
    if key.startswith('otpauth://'):
        key = key.split('secret=')[1].split('&')[0]
    key = key.replace(' ', '')
    key = key.strip()
    return hotp(key, int(time.time() / time_step), digits, digest)