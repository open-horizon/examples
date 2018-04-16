## Installing Horizon Software on Your Device

Interested parties are encouraged to try Horizon on their own machines.

Currently **ARM (32bit and 64bit)**, **x86**, and **IBM Power8** architectures are supported.  To run Horizon on your own machine you will normally install our Debian Linux package.

Instructions below will guide you through installation of Horizon based on your variety of machine.  The links below will jump ahead to the instructions for that machine type:

* [Raspberry Pi3 or Pi2]()
* [NVIDIA Jetson TX2 or TX1]()
* [x86 machines (or VMs)]()
* [IBM Power8 machines]()

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

**Now continue to the "Registering Your Horizon Machine" section below.**

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

**Now continue to the "Registering Your Horizon Machine" section below.**

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

**Now continue to the "Registering Your Horizon Machine" section below.**

### IBM Power8 Machines:

Horizon is moving toward support for cognitive and deep learning supercomputing platforms including the IBM Power8 Architecture with NVIDIA's PASCAL GPU.  Horizon has been tested and is running on a Power S814, 8 core machine, with 1TB RAM.

Begin by following **these setup instructions for IBM Power8**, then...

**Now continue to the "Registering Your Horizon Machine" section below.**

### Additional Platforms

Horizon is always extending the system to run on additional platforms. For more details on current beta development including participation with us, or just to get up to speed on our latest, contact us at our Discourse Forum.

### Registering Your Horizon Machine

Once you have installed the Horizon software (using the instructions above) it will take a few minutes to get up and running.  It will usually take less than 10 minutes.  Once Horizon is running you can register and add your machine onto the Horizon network.

#### Using the CLI to Register on Horizon

Optionally, it is recommended that you begin by verifying that Horizon has had time to come up and is "active" by running this command:

```
systemctl status horizon.service
```

The output of that command should contain lines similar to this:

```
 * horizon.service - Service for Horizon control system (cf. https://bluehorizon.network)
   Loaded: loaded (/etc/systemd/system/horizon.service; enabled)
   Active: active (running) since Thu 2017-07-06 13:35:54 UTC; 36min ago
```

Optionally, you may also wish to verify the version of Horizon on your machine:

```
dpkg -l | grep horizon
  ii  bluehorizon                     2.13.9~ppa~raspbian.jessie   armhf        Configuration for horizon package to use Blue Horizon backend
  ii  bluehorizon-ui                  2.13.9~ppa~raspbian.jessie   armhf        Web UI content for Bluehorizon instance of the Horizon platform
  ii  horizon                         2.13.9~ppa~raspbian.jessie   armhf        The open source version of the Horizon platform
```

The version should be 2.13.9 or later.

First you need to register your edge machine as a "node" in the Horizon Exchange.  You need a user account in the Horizon Exchange to do this, but if you are just using the "public" organization, an account will be created for you automatically when you register your new node.  You can make up any node ID, any node token (i.e., a password for this node), any username, and any password, and you must provide an email address as well, then run the "hzn exchange node create ..." command as shown below:

```
export NODE_ID=_____
export NODE_TOKEN=_____
export USERNAME=_____
export PASSWORD=_____
export EMAIL="_____@____.com"
hzn exchange node create -u $USERNAME:$PASSWORD -e $EMAIL -n $NODE_ID:$NODE_TOKEN -o public
```

To register your edge machine you will need to select the deployment pattern that you wish to run on this machine.  You can list the available patterns for the "public" organization by passing your Horizon Exchange credentials in the command shown below:

```
hzn exchange pattern list -o public -u $USERNAME:$PASSWORD
```

Select an appropriate pattern from this list based upon the variety of edge machine you will be registering, and the hardware peripherals it has attached.  Depending upon the pattern you select, you may also need to provide configuration details required by that pattern.  For example, to run the public pattern "netspeed-arm" (appropriate only for ARM architecture machines) you must provide the location information it requires.  To do this you would create a file, say ./input.json containing something  like the following.  This example registers a machine near San Jose, California, USA (at GPS coordinates 37.0, -121.0) and specifies that the location of the machine should always be obfuscated somewhere within a radius of 5.0 km from the actual location before being shared:

```
{
  "global": [
    {
      "type": "LocationAttributes",
      "variables": {
        "lat": 37.0,
        "lon": -121.0,
        "use_gps": false,
        "location_accuracy_km": 5.0
      }
    }
  ]
}
```

In the current directory, create the input.json file shown above (but with the values shown in bold above changed to reflect your own location).  This input file is specific to the netspeed-arm pattern.  When that file is ready, use the command shown below to deploy this pattern onto your edge machine

```
export PATTERN=netspeed-arm
hzn register -n $NODE_ID:$NODE_TOKEN -f ./input.json public $PATTERN
```
Once the hzn register command succeeds without error, then the Horizon AgBots will be able to find it, and will begin trying to establish agreements with it.  You can check the configuration and current status of Horizon on your edge machine with the commands shown below and very soon you should see an agreement showing up, then shortly after that the docker containers for that agreement should start running:

```
hzn agreement list
hzn service list
hzn workload list
```
Note that you can get more detailed instructions for the hzn command by passing "--help", e.g.:

```
hzn --help
```
You can also get detailed instructions for hzn subcommands (like "register") by passing "--help", e.g.:

```
hzn register --help
```

Your machine Is now registered on Horizon!

### Welcome to Horizon!

Once your machine registration and pattern assignment is completes, the Horizon AgBots discover your edge machine, propose agreements, and monitor their progress once your machine accepts them.  You should never again need to connect directly to the machine to insteall or update its software.  The "Horizon Insights" (in the pattern you selected) will now begin contracting with your device through the Horizon Exchange.  Once your device is in agreement, Horizon Microservices and Workloads (running in Docker containers) will be downloaded and will begin to run. Your device will soon be running the Horizon Insights you specified.

Soon your machine will be visible on the Horizon Unified Map.  Once you see your "dot" appear on that map, tap on it, and then select any of the Horizon Insight icons that appear (one for each of the Horizon Insights in your pattern).  When you select a Horizon Insight icon you will be taken to a page for that Horizon Insight where you can, for example, see visualizations of the data being sent from your machine, or interact directly with the hardware on your edge machine (depending upon the particular Horizon Insights you have selected).

### Questions, Troubleshooting

If you have any difficulties with any of these steps, you may find these other Horizon douments useful:

* [Edge-Developer-Quickstart-Guide]()
* [Edge-Quick-Start-Guide.md]()
* [Edge-Service-Development-Guidelines]()
* [Frequently-Asked-Questions]()
* [Troubleshooting]()
