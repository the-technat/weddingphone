# weddingphone

A guestbook that's not a book but a phone

## Concept

Have you every been at a wedding and filled out a guestbook? How would you describe this experience? Wouln'd a phone, where you can record your guestbook entry on the phone be a nice alternative? That's exactly what the weddingphone is trying to do.

## Development

### gokrazy 

Setup a testing instance of gokrazy by inserting an SD card into your machine and installing the following tools:

```console
yay -S go
go install github.com/gokrazy/tools/cmd/gokr-packer@latest
```

If you want to reach your instance over wifi, run the following:

```console
echo '{"ssid": "Secure WiFi", "psk": "secret"}' > extrafiles/github.com/gokrazy/wifi/etcd/wifi.json
```

Figure out which device your SD card is, mine is `/dev/mmcblk0` for this quickstart.

Then your very first instance bootstrap can be done using the following command:

```console
gokr-packer \
  -tls=self-signed \
  -overwrite=/dev/mmcblk0 \
  -hostname weddingphone \
  -serial_console=disabled \
  github.com/gokrazy/fbstatus \
  github.com/gokrazy/serial-busybox \
  github.com/gokrazy/breakglass \ 
  github.com/gokrazy/wifi \
  github.com/the-technat/weddingphone
```

Note: remove the `-serial_console=disabled` If you want your primary console to be serial. 

#### Update

**Make sure your instance is rechable using the hostname `weddingphone`.*

Once bootstrapped and accessable over WiFi (see the web interface of gokrazy), you can update programms on the instance like so:

```console
gokr-packer \ 
  -tls=self-signed \ 
  -update=yes \ 
  -hostname weddingphone \
  -serial_console=disabled \
  github.com/the-technat/weddingphone
```

