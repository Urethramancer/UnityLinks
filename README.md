# UnityLinks [![Build Status](https://travis-ci.org/Urethramancer/UnityLinks.svg)](https://travis-ci.org/Urethramancer/UnityLinks)
This is a server and link parser to get the individual official package links for releases of [Unity](https://unity3d.com). Supported links at the moment are only for 5.0 and up, but others might work. See usage below.

An example server exists at http://unity.grimdork.net, and is updated roughly on the day a new patch or regular version is released. Betas are not tracked.

## Why?
The installer sometimes fails behind firewalls, while direct web downloads, which is what the installer actually uses, work fine in a browser. Sometimes the downloads simply fail before they're done, and start over from zero when you retry. Getting the direct links lets you pick and choose, continue broken downloads and package it up for wider internal use, something the installer makes more complicated than it should.

## Warning
The current code is an ugly mess because gcfg couldn't parse the INI files produced for the Unity installer in every case. A very naïve parsing method is used due to the regularity of these files. Do not read - your eyes will bleed.

## Dependencies
1. [web.go](https://github.com/hoisie/web)
2. [str](github.com/mgutz/str)

## Build
Clone the repo and build with Go. Latest generally works (see build status link in the title above).

## Usage
### Server
To run the server, just run the binary on any Unix-like system, preferably some form of Linux. Any required folders will be created. Use whichever method you prefer to keep it running (Docker, LCX, fancy daemon scripts, wishful thinking).

You can bind it to a specific address with the *-a* flag and a port of choice with the *-p* flag.

### Updating
Its other use is a little more complicated, as it needs a couple of external pieces of information.

You will need to copy a link to a Unity version you want the package links for. It will generally look something like this for the Mac version:

> http://beta.unity3d.com/download/d64ba7d31ce9/UnityDownloadAssistant-5.3.6p5.dmg

Or like this for the Windows version:
>http://beta.unity3d.com/download/d64ba7d31ce9/UnityDownloadAssistant-5.3.6p5.exe

There are two pieces of information you want from this: The hash and the version.

In this case the hash is *d64ba7d31ce9* and the version is *5.3.6p5*.

Create a file with the name of the version in the updates directory in the same folder as UnityLinks, containing only the hash on one line:

```sh
echo d64ba7d31ce9 > updates/5.3.6p5
```

Repeat for any number of versions you'd like to extract links for.

Then update the versions:

```sh
./UnityLinks -u
```
or

```sh
./unitylinks -u
```

The running server should now have an extra page of links for each version, if all went well. The update method uses a few different links out of the possible ones Unity Technologies have available, which may help when the regular installer fails.

## Licence
MIT.
