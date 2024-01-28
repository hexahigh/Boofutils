# Boofutils
![Build and release](https://github.com/hexahigh/boofutils/actions/workflows/build&release.yml/badge.svg)
[![Release](https://img.shields.io/github/release/hexahigh/boofutils.svg)](https://github.com/hexahigh/boofutils/releases)
[![License](https://img.shields.io/github/license/hexahigh/boofutils)](https://github.com/hexahigh/boofutils/blob/main/LICENSE)
[![Downloads](https://img.shields.io/github/downloads/hexahigh/boofutils/total.svg)](https://github.com/hexahigh/boofutils/releases)
![Go report card](https://goreportcard.com/badge/github.com/hexahigh/boofutils)<br>
A utility program im working on.
<br>
Booftutils is very unigue.
In contrary to other utility programs, Boofutils uses a sh*t ton of memory.
![alot of memory](https://pomf2.lain.la/f/zxi1cpji.png)
## Download
If you are using a debian based distro you can install boofutils using my apt repository.
Simply run these commands:
```bash
echo "deb [signed-by=/usr/share/keyrings/boofdev.apt.pub] https://apt.080609.xyz stable main" | sudo tee -a /etc/apt/sources.list.d/boofdev.list && sudo wget -q -O /usr/share/keyrings/boofdev.apt.pub https://apt.080609.xyz/pgp-key.public
sudo apt update
sudo apt install -y boofutils
```

### Building from source:
You will need go version 1.21 or higher. If you do not have it installed then i recommend using [GVM.](https://github.com/moovweb/gvm)

You will also need to install libasound and pkg-config.
 If you are using a Debian based distro, you can install them with `sudo apt install -y libasound2-dev pkg-config git`.

 After you have installed them you can simply clone the repo and build it. These commands should be the same on all operating systems.
```bash
git clone https://github.com/hexahigh/boofutils
cd boofutils
go build
```
If you are using Windows or don't want to build from source, you can download the latest release below.

You can download the latest nightly version [here.](https://github.com/hexahigh/boofutils/releases/tag/latest_auto) (Recommended)<br>
Or, you can download the latest release [here.](https://github.com/hexahigh/boofutils/releases/latest)
