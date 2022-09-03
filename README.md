# weddingphone

A guestbook that's not a book but a phone

## Intro

Have you every been at a wedding and filled out a guestbook? How would you describe this experience? Wouln'd a phone, where you can record your guestbook entry on the phone be a nice alternative? That's exactly what the weddingphone is trying to do.

This repo contains the go programm that runs on a raspberry pi using gokrazy. The raspberry pi has a microphone and speaker attached that are nicely wired into an old phone.

## Development

### Design principles

- It's a personal project, we currently don't do TDD
- Save recorded audio files locally and if there is networking connectivity upload them to S3
- Log to stdout, currently there is no need to log to files
- Record audio directly within go and buffer to memory to avoid too many writes to disk
- In general, store in memory and do as few writes to disk as possible (SD cards have limited lifetime)

### Development environment

Of course you need `go` to develop. But there are some more tools and hardware you need.

- Ensure you have `go` 1.19 installed
- Ensure you have the go programm `github.com/gokrazy/tools/cmd/gokr-packer@latest` installed and GOBIN in your PATH
- Ensure you have a Raspberry Pi 2/3B/3B+ with a decent power supply and an SD card (size doesn't matter)
- For easier development I recommend you also have a serial to usb cable that can be attached to the raspberry pi (something [like this](https://www.pi-shop.ch/usb-to-ttl-serial-kable-debug-console-kable-fuer-den-raspberry-pi))
  - of course a monitor and HDMI cable does the trick too

### gokrazy

As mentioned in the intro, we are using gokrazy for this project (no Raspberry Pi OS or other Pi friendly linux distro). This has multiple advantages which you can read more about [here](https://gokrazy.org/).

So to bootstrap your Raspberry Pi using gokrazy and our programm, insert your SD card into the computer and find it's device path. This quicksatrt uses `/dev/mmcblk0` as device path.

Next you need to decide whether you want to use WiFi or Ethernet (for the rest of the Pis life time it will be managed over the network).

For WiFi run the following on the command line within the cloned git repo:

```console
echo '{"ssid": "Secure WiFi", "psk": "secret"}' > extrafiles/github.com/gokrazy/wifi/etc/wifi.json
```

Note: if you get an error 'path does not exist', run `mkdir -p extrafiles/github.com/gokrazy/wifi/etc/`.

For Ethernet you can remove the `github.com/gokrazy/wifi` package from the command below.

Make sure your Raspberry Pi can be resolved using the hostname you specified below (e.g mine would be `weddingphone.silver.lan`)

Then your very Pi can be bootstraped using the following command:

```console
gokr-packer \
  -tls=self-signed \
  -overwrite=/dev/mmcblk0 \
  -hostname weddingphone.silver.lan \
  -serial_console=disabled \
  github.com/gokrazy/fbstatus \
  github.com/gokrazy/serial-busybox \
  github.com/gokrazy/breakglass \ 
  github.com/gokrazy/wifi \
  github.com/the-technat/weddingphone
```

Some notes:

- `-tls=self-signed` enables the web UI to be access over HTTPS
- `-overwrite=/dev/mmcblk0` specifies that this is the first bootstrap and that it should go to this SD card
- `-hostname weddingphone.silver.lan` specifies how this device is named and for which domain the certificate is issued
- `-serial_console=disabled` says that the device outputs some status informations on the HDMI port but has no other method of interactin with the device -> remove this if you want to see everything the device outputs on serial (primary) and HDMI (secondary)
- `github.com/gokrazy/fbstatus` -> very minimal status display that can be seen on the HDMI output
- `github.com/gokrazy/serial-busybox` -> very minimal console that can be accessed either via serial (if enabled) or using a monitor/keyboard
- `github.com/gokrazy/breakglass` -> minimal ssh access to the device's serial-busybox console
- `github.com/gokrazy/wifi` -> WiFi support for the device (remove this package if you only want to use Ethernet)
- `github.com/the-technat/weddingphone` -> the actual programm that we want to install

All of the specified programms are stared automatically at boot and do what they are designed for. You can remove any of them if you don't need them. I recommend for a production ready device that you remove the `breakglass` programm but enable the serial console...

#### Work with gokrazy devices

Once you plug in your device, it will boot up and establish a network connection over DHCP. You can then access the web interface using your defined hostname (mine would be [https://weddingphone.silver.lan](https://weddingphone.silver.lan)).
There you can see your running programms, their stdout/stderr as well as any environment variables, some os stats and that's it. Very minimal.

If you want to update the weddingphone programm you can use the following command do to so:

```console
gokr-packer \ 
  -tls=self-signed \ 
  -update=yes \ 
  -hostname weddingphone.silver.lan \
  -serial_console=disabled \
  github.com/the-technat/weddingphone
```

This will compile and push the new version of the weddingphone programm to the Pi and then restart the device.

Of course you can also do a full upgrade of the entire gokrazy instance, see [here](https://github.com/gokrazy/gokrazy#updating-your-installation) for instruction on that.
