# wtk---Web-Toolkit
A collection of web developer tools, packed into an Assembly wrapper.

# IMPORTANT
This project might not fit your needs and most certainly will not work on Windows unless started from within `wsl`. It uses POSIX compliant syscalls, but there still might be undefined behavior on UNIX based systems.  
It was tested solely on Linux machines

Building:
`make build_and_link` can be used to build the **wtk** file.

`wtk mfv .` can be used as an example usage. **wtk** is obviously the wrapper call, **mfv** is the app to be called, . is the input - in this case its the current directory.

## Warning:
Currently developed and tested on Linux.  
  
The project uses POSIX compliant interfaces and mostly Linux-specific behavior.  
Windows is only expected to work through WSL.

## Needed Packages
1. `go` the crawler is written in Go and as such requires it installed in order to build the application yourself  
2. `gcc` for obvious reasons  
3. `fasm` to compile the final executable  
  
The repo includes already built files for the sake of shipping with precompiled binaries. The main executable is located in the `build` directory, and its called `wtk`. **It might require you to give it execute privileges - `chmod +x wtk`**.

## Better descriptions of each application and benchmarks

- ### wtk <app> <input>
wtk is the main executable, it is where each other app can be called from. It is written in **Assembly** and acts as the delegator of the whole application

- ### tk mfv <input>
mfv is the Missing File Validator program. Can be used with a directory or a single file as input. If a directory is provided it will recursively walk all the sub-directories, collect image files along the way while also upon encountering text files like html as one example, will look for references to files in the filesystem and finally produce an output of all the missing files.
  - Provides a list of all the files that contain missing content
  - Provides the full path to each file that references missing content
  - Provides a list of all the missing content under each file that references them
  - Provides a total of all the missing files found
  - Colored output for easier browsing

```
tk tm wtk mfv .
Scan Complete
File Path: ./index.html
        Not Found: ./content/image.png
        Not Found: ./content/image.gif
        Not Found: ./content/image.jpg
        Not Found: ./content/image.jpeg
        Not Found: ./content/image.jpe



Total files missing: 5

-----------------------------------
RAM used: 3.203 MB
Execution time: 0.052840 seconds
-----------------------------------
```

- ### wtk crawler <input>
crawler is the Crawler program. It accepts a URL as input. It creates an SQLite database in `~/.local/share/wtk/db/crawler.db` upon starting, and begins crawling the website. It uses a simple state machine in order to see what is left to crawl. Upon encountering a link, it checks whether the host is the same, and if it is, it pushes it into the database for later, if not however, the crawler discards it. The reason for that being, not to leak into unwanted websites and target only the host you provided.  
  
For example if provided with `https://my-website.com`, the crawler will never go to `https://your-website.com` even if both have references to each other.  
The crawler collects the hosts, compares uniqueness in order not to crawl the same pages more than once, self restarts from where it left off if forcefully stopped with `Ctrl + c` or any other way, and stores all the hosts, status codes, headers and html, plus other things.

- ### wtk mirror <input>
mirror is the Mirror program. Currently not implemented yet.

- ### wtk httpserver <input>
httpserver is the HttpServer program. Currently not implemented yet.
