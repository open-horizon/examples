## Installing Horizon Software on Your Device

Interested parties are encouraged to try Horizon on their own machines.

Currently **ARM (32bit and 64bit)**, **x86**, and **IBM Power8** architectures are supported.  To run Horizon on your own machine you will normally install our Debian Linux package.

Instructions below will guide you through installation of Horizon based on your variety of machine.  The links below will jump ahead to the instructions for that machine type:

* [Raspberry Pi3 or Pi2](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#arm-32-bit-raspberry-pi)
* [NVIDIA Jetson TX2 or TX1](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#arm-64-bit-nvidia-jetson)
* [x86 machines (or VMs)](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#x86-machines)
* [IBM Power8 machines](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#ibm-power8-machines)
* [Additional Platforms](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#additional-platforms)

If you wish to bring any other device not listed here onto Horizon, please contact the Horizon community by clicking the Forum tab at the top of this page.

### ARM 32 bit, Raspberry Pi:

#### Hardware Requirements
* Raspberry Pi 3 (recommended)
* Raspberry Pi 2
* MicroSD flash card (32GB recommended)
* An appropriate power supply for your machine.  At least a 2 Amp power supply is recommended (and please note that more than 2 Amps may be required if power-hungry USB peripherals are being powered from your Raspberry Pi).
* An internet connection for your machine, wired or WiFi (and please note that the Pi 2 will require additional hardware for WiFi)
* Sensor hardware (optional, but most Horizon-Insight applications require optional special-purpose sensor hardware.  For example, Software Defined Radio (SDR) Insights require an SDR USB dongle with an antenna. However, Horizon's Netspeed application requires no additional hardware.  Please see our full hardware list to discover appropriate Raspberry Pi sensor hardware options.)

#### Procedure

Flash an appropriate debian-based Linux image (e.g., Rasbian) onto your MicroSD card (the instructions below in this section assume Raspbian for wifi and ssh configuration).  Note that the Raspberry Pi Foundation provides good instructions for flashing MicroSD images from many platforms.
WARNING: Flashing an image onto your MicroSD card will completely erase it (any previous contents will be lost permanently).
Edit your newly flashed image to provide appropriate WPA2 WiFi credentials (unless you will be using a wired network connection):
Find and edit the file wpa_supplicant_cred.txt in the MicroSD card's top level folder
Add your WiFi credentials (i.e., your network's SSID name and your passphrase) using the format shown below.  That is, edit the example text below to replace **YourSSID** and **YourPassphrase** with values appropriate for your network.  PLEASE NOTE: Omit quotes for these fields, even if your SSID or passphrase contain spaces or tabs!

```
wpa_ssid=YourSSID
wpa_passphrase=YourPassphrase
```

Save and close the file.
If you will be working with your Raspberry Pi "headless" (i.e., with no monitor or keyboard) you must also create an empty file named "ssh" in the MicroSD card's top level folder to enable you to ssh into this machine with the default credentials.  To protect your Pi this is disabled by default.
When you have done the steps above, unmount the MicroSD card properly (safely eject) so all of your changes will get written
Insert the MicroSD card into your Raspberry Pi, attach any optional sensor hardware (see above), then connect the power supply
Change The Default Password:

Once your Pi boots up, you are strongly advised to immediately change the default password.
In the Horizon Raspbian flash images, the default account uses login name: "pi", and has password: "raspberry".
Login to this account then use the standard Linux passwd command to change the password, e.g.:

```
$ passwd
Enter new UNIX password: ...
Retype new UNIX password: ...
passwd: password updated successfully
$ 
```

**Install the Horizon Software**

The steps below expect you to have root privileges, so if you are logged in as the pi user, execute this to begin:

```
sudo -s
```

Ensure that you have the current docker version installed (since many distros are set up to run much older docker versions):

```
curl -fsSL get.docker.com | sh
```

Now configure apt to be able to install the latest stable version of Horizon:

```
wget -qO - http://pkg.bluehorizon.network/bluehorizon.network-public.key | apt-key add -
aptrepo=updates
# aptrepo=testing    # or use this for the latest, development version
cat <<EOF > /etc/apt/sources.list.d/bluehorizon.list
deb [arch=$(dpkg --print-architecture)] http://pkg.bluehorizon.network/linux/ubuntu xenial-$aptrepo main
deb-src [arch=$(dpkg --print-architecture)] http://pkg.bluehorizon.network/linux/ubuntu xenial-$aptrepo main
EOF

apt-get update
Install the latest Horizon software:
apt-get install -y horizon bluehorizon bluehorizon-ui
```

[Now continue to the "Registering Your Horizon Machine" section below.](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#registering-your-horizon-machine)

### ARM 64 bit, NVIDIA Jetson:

#### Hardware Requirements

* NVIDIA Jetson TX2 (recommended)
* NVIDIA Jetson TX1
* HDMI Monitor, USB hub, USB keyboard, USB mouse
* Storage: >10GB (SSD recommended)
* An internet connection for your machine (wired or WiFi)
* Sensor hardware (optional, but most Horizon-Insight applications require optional special-purpose sensor hardware.  For example, Horizon Aural Insights require USB sound card, analog microphone, and headphones or speaker hardware, and Horizon Eye Insights require webcam hardware.  However, Horizon's Netspeed application requires no additional hardware.  Please see our full hardware list to discover appropriate Raspberry Pi sensor hardware options.)
 

#### Procedure

Begin by using these links to open source instructions, for TX1, or for TX2, to get the current NVIDIA JetPack installed on your machine.  These links also include instructions for setting up Docker and other prerequisites you will need before installing the Horizon software.
The steps below expect you to have root privileges, so if you are logged in as the nvidia user, execute this to begin:

```
sudo -s
```

Now configure apt to be able to install the latest stable version of Horizon:

```
wget -qO - http://pkg.bluehorizon.network/bluehorizon.network-public.key | apt-key add -

cat <<EOF > /etc/apt/sources.list.d/bluehorizon.list
deb [arch=arm64] http://pkg.bluehorizon.network/linux/ubuntu xenial-updates main
deb-src [arch=arm64] http://pkg.bluehorizon.network/linux/ubuntu xenial-updates main
EOF

apt-get update
```

**Install the latest Horizon software:**

```
apt-get install -y horizon bluehorizon bluehorizon-ui
```

**Start the Horizon service:**

```
systemctl start horizon.service
```
At this point you no longer require root privileges, so can exit this privileged shell and return to running as the nvidia user:

```
exit
```

**Change The Default Password:**

Once your NVIDIA TX2 or TX1 is up, you are strongly advised to immediately change the default password.
In the JetPack installation procedure the default account uses login name: "nvidia", and has password: "nvidia".
Login to this account then use the standard Linux passwd command to change the password, e.g.:

```
$ passwd
Enter new UNIX password: ...
Retype new UNIX password: ...
passwd: password updated successfully
$ 
```

[Now continue to the "Registering Your Horizon Machine" section below.](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#registering-your-horizon-machine)

### x86 Machines:

#### Hardware Requirements

* 64 bit Intel or AMD machine or VM
* An internet connection for your machine (wired or WiFi)
* Sensor hardware (optional, and note that currently only Horizon's Netspeed insight application, which requires no additional sensor hardware is supported on x86 machines, although some of the other Horizon-Insight applications may work with appropriate sensor hardware and configuration)

#### Procedure

Install a recent Debian Linux variant (e.g., ubuntu 16.04, which is used for the instructions below)
Configure apt to be able to install the latest stable version of Horizon:

```
wget -qO - http://pkg.bluehorizon.network/bluehorizon.network-public.key | apt-key add -

cat <<EOF > /etc/apt/sources.list.d/bluehorizon.list
deb [arch=amd64] http://pkg.bluehorizon.network/linux/ubuntu xenial-updates main
deb-src [arch=amd64] http://pkg.bluehorizon.network/linux/ubuntu xenial-updates main
EOF

apt-get update
```

**Install the latest Horizon software:**

```
apt-get install -y horizon bluehorizon bluehorizon-ui
```

**Configure Horizon workload logging (recommended)**

```
cat <<'EOF' > /etc/rsyslog.d/10-horizon-docker.conf
$template DynamicWorkloadFile,"/var/log/workload/%syslogtag:R,ERE,1,DFLT:.*workload-([^\[]+)--end%.log"

:syslogtag, startswith, "workload-" -?DynamicWorkloadFile
& stop
:syslogtag, startswith, "docker/" -/var/log/docker_containers.log
& stop
:syslogtag, startswith, "docker" -/var/log/docker.log
& stop
EOF

service rsyslog restart
```

**Start the Horizon service:**

```
systemctl start horizon.service
```

[Now continue to the "Registering Your Horizon Machine" section below.](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#registering-your-horizon-machine)

### IBM Power8 Machines:

Horizon is moving toward support for cognitive and deep learning supercomputing platforms including the IBM Power8 Architecture with NVIDIA's PASCAL GPU.  Horizon has been tested and is running on a Power S814, 8 core machine, with 1TB RAM.

Begin by following [these setup instructions for IBM Power8](https://staging.bluehorizon.network/documentation/adding-your-power-device), then...

[Now continue to the "Registering Your Horizon Machine" section below.](https://github.com/open-horizon/examples/blob/master/edge/doc/Adding-Devices.md#registering-your-horizon-machine)

### Additional Platforms

Horizon is always extending the system to run on additional platforms. For more details on current beta development including participation with us, or just to get up to speed on our latest, contact us at our Discourse Forum.

### Registering Your Horizon Machine

Once you have installed the Horizon software (using the instructions above) it will take a few minutes to get up and running.  It will usually take less than 10 minutes.  Once Horizon is running you can register and add your machine onto the Horizon network.

Now use one of these guides to get your machine registered:

* [Edge-Quick-Start-Guide.md](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md)
* [Edge-Developer-Quickstart-Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Developer-Quickstart-Guide.md)

### Welcome to Horizon!

Once your machine registration is complete, the Horizon AgBots will discover your edge machine, propose agreements, and monitor their progress once your machine accepts them.  You should never again need to connect directly to the machine to install or update its software.

### Questions, Troubleshooting

If you have any difficulties with any of these steps, you may find these other Horizon douments useful:

* [Edge-Developer-Quickstart-Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Developer-Quickstart-Guide.md)
* [Edge-Quick-Start-Guide.md](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md)
* [Edge-Service-Development-Guidelines](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Service-Development-Guidelines.md)
* [Frequently-Asked-Questions](https://github.com/open-horizon/examples/blob/master/edge/doc/Frequently-Asked-Questions.md)
* [Troubleshooting](https://github.com/open-horizon/examples/blob/master/edge/doc/Troubleshooting.md)
