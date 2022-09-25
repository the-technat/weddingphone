# The consrv hostname resolves to the device’s Tailscale IP address,
# once Tailscale is set up.
PACKER := gokr-packer -hostname=weddingphone

PKGS_TAILSCALE := \
	github.com/gokrazy/breakglass \
	github.com/gokrazy/timestamps \
	github.com/gokrazy/serial-busybox \
	github.com/gokrazy/stat/cmd/gokr-webstat \
	github.com/gokrazy/stat/cmd/gokr-stat \
	github.com/gokrazy/mkfs \
	github.com/gokrazy/wifi \
	tailscale.com/cmd/tailscaled \
	tailscale.com/cmd/tailscale \
  github.com/the-technat/weddingphone

PKGS := \
	github.com/gokrazy/breakglass \
	github.com/gokrazy/timestamps \
	github.com/gokrazy/serial-busybox \
	github.com/gokrazy/stat/cmd/gokr-webstat \
	github.com/gokrazy/stat/cmd/gokr-stat \
	github.com/gokrazy/wifi \
  github.com/the-technat/weddingphone

all:

.PHONY: update overwrite

update:
	${PACKER} -update=yes ${PKGS_TAILSCALE}

update-no-tailscale:
	${PACKER} -update=yes ${PKGS}

overwrite:
	${PACKER} -overwrite=${card} ${PKGS_TAILSCALE}

overwrite-no-tailscale:
	${PACKER} -overwrite=${card} ${PKGS}
