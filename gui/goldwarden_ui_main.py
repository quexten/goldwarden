#!/usr/bin/env python3

import platform

if platform.system() == 'Darwin':
    import src.macos.main as macos_main
    macos_main.main()
elif platform.system() == 'Linux':
    import src.linux.main as linux_main
    linux_main.main()
else:
    print("Unsupported OS " + platform.system() + "... exiting...")