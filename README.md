# wtk---Web-Toolkit
Assembly wrapper for some web developer tools.

# IMPORTANT
This project might not fit your needs and most certainly will not work on Windows unless started from within `wsl`. It uses POSIX compliant syscalls, but there still might be undefined behavior on UNIX based systems.  
It was tested solely on Linux machines

Building:
`make build_and_link` can be used to build the **wtk** file.

`wtk mfv .` can be used as an example usage. **wtk** is obviously the wrapper call, **mfv** is the app to be called, . is the input - in this case its the current directory.
