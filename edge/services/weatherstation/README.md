## ZEUS - Horizon WAVE Weather POC 

### Repo contents: Docker container build scripts and supporting files
* Docker Build script (Alpine Linux Image)
* Required files
    * Python libraries not available using apk add
    * Weewx python scripts / package
    * Python scripts for config file modification w/ IBM account info
    * Supporting weewx install and required linux sys files
    * start.sh file (runs weewx as a service)

### Usage:
* (Option 1) Connect weather station to ARM-based device)
    * (RPi 2/3 or XU3/4 build device)
    * USB connect to Pi/XU4 (currently supported) OR wifi connect (TODO)
* (Option 2) Follow steps below and enable weewx weather station simulator last
* Git Clone files below onto ARM-based device
* cd into directory ./zeus
* Build docker container image using:

    > docker build -t zeus:test --force-rm .
    
* (Optional, for option 2)
    * Start weather container with all those env variables / open bash shell
    
            > docker run --name zeus -it --privileged -v /dev/bus/usb:/dev/bus/usb -e MTN_CONTRACT=sdjsdds -e MTN_LON=-118.2250 -e MTN_LAT=33.8967 -e MTN_PWS_WU_ID=KCACOMPT9 -e MTN_PWS_ALT='150,foot' -e MTN_PWS_LOC='Compton, CA Next to Roscoes' -e MTN_PWS_ST_TYPE=FineOffsetUSB -e MTN_PWS_MODEL=WS2080 -e MTN_PWS_UNITS=us -e MTN_PWS_WU_RPDF=True zeus /bin/bash
    * Run weewx config utility (see [wee_config](http://www.weewx.com/docs/usersguide.htm#wee_config_utility)) for help
        * cd to weewx bin directory 
        
            > cd /home/weewx/bin
        * To test without a weather station: run wee_config and change driver to simulator
        
            >./wee_config --reconfigure --driver=weewx.drivers.simulator --no-prompt
        * verify configuration was saved
        * Edit the file /home/weewx/bin/weewxd and comment out the active Horizon HA proxy lines (lines 18 & 19)
        
            > #These two URLs can be commented to run weewx as normal (no redirect)
            > #wr.StdWunderground.rapidfire_url="http://169.53.229.90:9000/rapidfire/weatherstation/updateweatherstation.php"
            > #wr.StdWunderground.archive_url="http://192.53.229.90:9000/archive/weatherstation/updateweatherstation.php"

    * Restart weewx service:
    
        > /etc/init.d/weewx reload
        
    * Confirm weewx is operational by enabling debug mode (syslog disabled in Alpine
    Linux.... debug info will be printed to stderr)

        > ./wee_config --reconfigure --debug
        
        > /etc/init.d/weewx reload  (if already running)
        
        > /etc/init.d/weewx start   (if not already running)
        
* **TODO: Continued development**
    * Update device reg page to prompt for and enable PWS credentials, register
    PWS on IBMBlueHorizon weatherunderground.com account (Ling)
    * Update this container with appropriate environment variables according to
    current convention (MTN_LAT, MTN_LON) and new for PWS (MTN_PWS_WU_ID). More
    definition in weewx_mod.py and on wiki (link at bottom)

### Files
|                       |                                                     |
|-----------------------|-----------------------------------------------------|
Dockerfile              | Docker build file                                   |
answers.txt             | File w/ answers to automated install script for weewx. Values are specific to those expected by weewx_mod.py |
init-functions          | Debian sys file, expected by weewxd when run as service. (Lacks robust testing: copied into /lib/lsb/ as script expects, works for this alpine linux install) |
pyusb-1.0.0b2.tar.gz    | python-usb supporting libraries package (not available using apk add in alpine) |
start.sh                | bash script to be run in weather, container; checks MTN contract, modifies weewx config, and starts weewx service |
vars.sh                 | See explanation for 'init-functions' |
weewx-3.4.0.tar.gz      | Weewx Python scripts and installation script; core code for weather POC. Supports USB-enabled and some wifi-enabled PWS (personal weather stations.) |
weewx.conf.sample       | Sample weewx.conf file for reference (unused in build) |
weewx_mod.py            | Python script; modifies weewx_conf file with values specific to this Pi's weather station (must coordinate with PWS registered at wunderground.com by device reg. page backend (to-do). |
weewxd                  | Weewx script, run as a service, which replaces default provided in weewx*.gz file. This file has been modified with URL pointing to Horizon HA Proxy on IBM's Softlayer machines. Proxy modifies outgoing weewx HTTP API request to WU with replacement password of IBMBlueHorizon WU account. |
                      
                        
### References:
* [Weather POC Wiki](https://repo.hovitos.engineering/MTN/mtn/wikis/wave_weather_POC)
* pywws: http://pywws.readthedocs.org/en/latest/index.html
* weewx: http://www.weewx.com/  (open source, many supported PWS's)