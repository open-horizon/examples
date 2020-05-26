##Configuring horizon node on RPi.

Install pre-requirements

```
sudo apt-get install jq
```

Get horizon binaries

```
arch=$(dpkg --print-architecture)
dist=buster
version=2.24.18
wget http://pkg.bluehorizon.network/linux/raspbian/pool/main/h/horizon/bluehorizon_${version}~ppa~raspbian.${dist}_all.deb
wget http://pkg.bluehorizon.network/linux/raspbian/pool/main/h/horizon/horizon-cli_${version}~ppa~raspbian.${dist}_${arch}.deb
wget http://pkg.bluehorizon.network/linux/raspbian/pool/main/h/horizon/horizon_${version}~ppa~raspbian.${dist}_${arch}.deb
```

And install them

```
sudo -s
dpkg -i horizon-cli_${version}~ppa~raspbian.${dist}_${arch}.deb
dpkg -i horizon_${version}~ppa~raspbian.${dist}_${arch}.deb
dpkg -i bluehorizon_${version}~ppa~raspbian.${dist}_all.deb

exit
```