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
- A [tailscale](https://tailscale.com) account to connect to your raspberry pi from everywhere
  - of course a serial to usb cable that can be attached to the raspberry pi (something [like this](https://www.pi-shop.ch/usb-to-ttl-serial-kable-debug-console-kable-fuer-den-raspberry-pi)) does the trick too
  - or use a monitor and HDMI cable

### gokrazy

As mentioned in the intro, we are using gokrazy for this project (no Raspberry Pi OS or other Pi friendly linux distro). This has multiple advantages which you can read more about [here](https://gokrazy.org/).

So to bootstrap your Raspberry Pi using gokrazy, insert your SD card into the computer and find it's device path. This quicksatrt uses `/dev/mmcblk0` as device path (Change your's in `Makefile` accordingly).

Next you need to decide whether you want to use WiFi or Ethernet (for the rest of the Pis life time it will be managed over the network).

For WiFi run the following on the command line within the cloned git repo:

```console
echo '{"ssid": "Secure WiFi", "psk": "secret"}' > extrafiles/github.com/gokrazy/wifi/etc/wifi.json
```

Note: if you get an error 'path does not exist', run `mkdir -p extrafiles/github.com/gokrazy/wifi/etc/`.

Next you need to know that I'm using [tailscale](https://gokrazy.org/packages/tailscale/) to connect to the raspberry pi using the hostname `weddingphone` from everywhere. If you want to use that too, get yourself an account at [tailscale.com](https://tailscale.com) and generate an [auth key](https://login.tailscale.com/admin/settings/keys) for the raspberry pi which you insert at bootstrap time like so:

```console
mkdir -p flags/tailscale.com/tailscale/
cat > flags/tailscale.com/cmd/tailscale/flags.txt <<EOF
up
--auth-key=tskey-AAAAAAAAAAAA-AAAAAAAAAAAAAAAAAAAAAA
EOF
```

If you don't want to use tailscale, just make sure you can reach the pi using the hostname `weddingphone` somehow.

Then you can bootstrap the SD card with gokrazy:

```console
# Using wifi and tailscale
make card=/dev/mmcblk0 overwrite

# Without tailscale
make card=/dev/mmcblk0 overwrite-no-tailscale
```

Put the SD card into your raspberry, plug in power and access the Web stats using `http://weddingphone`

If you want to update the weddingphone use `make update` (over the network).

Of course you can also do a full upgrade of the entire gokrazy instance, see [here](https://github.com/gokrazy/gokrazy#updating-your-installation) for instructions on that.
