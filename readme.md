# MDBrute

Brute-force memdatabase address search for RoM-Bot

### My anti-virus software flags this as malicious. Is it malware?

No, it is not malware. No data is saved or sent over the internet anywhere for
any purpose, and no modifications are made by this software -- neither in
memory, registry, nor files.

Anti-virus softwares use heuristics to detect unknown malware by looking for
common characteristics, patterns, and techniques.

Due to the way this software works (See notes about admin privileges below),
this may trigger certain anti-virus systems to signal a **false positive**.

#### I still don't trust you

Cool. Look over the source code yourself to see if there's anything suspicious.
If you are confident that the source is clean, install [Golang](https://go.dev/)
and compile it yourself.

## Compatibility
This tool was built for use with the Gameforge client.
It should function out-of-the-box for all client languages.

## Usage
* Log into the game. You must be fully loaded into the game world so that all assets are loaded into the client.
* Run MDBrute.exe from a console with **administrator privileges**.
* A found address should be found and displayed in the console; plug that into RoM-Bot's `addresses.lua` file for the value of `memdatabase.base`

### Why do I need to run this with admin privileges?

MDBrute uses `SeDebugPrivilege` to gain permission to read memory from other
processes. It uses this to _read_ (and only read) data from the Runes of Magic
client. In order to do this, administrator privileges are required.

Any memory that is read from the client is only used for pattern matching.
Specifically, looking for the static address that can be used to read
game information such as skill names/IDs, monster data, etc.


#### What do I do if no results were found?
Try running MDBrute with expanded search range; reduce the start address and increase the end address.
See the below table for information about the defaults for these values.

If this is still not successful, then the game client has changed too much and MDBrute
will need to be updated to be compatible.

## CLI arguments

| Argument      | Default              | Description |
|---------------|----------------------|-------------|
| proc          |                      | The process ID to target. If not given, one will be found automatically. |
| start         | 6422528 (0x00620000) | The lower bounds of where to begin searching; An offset from Client.exe |
| end           | 6619136 (0x00650000) | The upper bounds of where to begin searching; An offset from Client.exe |

**Example:**
```sh
mdbrute.exe --proc=123456 --start=6422528 --end=6619136
```


## License
Public domain. Use it as you wish.
