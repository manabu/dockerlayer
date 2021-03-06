[![wercker status](https://app.wercker.com/status/e65cc8012a3f1a833c621a7e53f7fd2e/m/master "wercker status")](https://app.wercker.com/project/byKey/e65cc8012a3f1a833c621a7e53f7fd2e)

# dockerlayer
Display some docker layer information.

Output Style inspired by [[Proposal]: docker diff between image layers · Issue #12641 · docker/docker](https://github.com/docker/docker/issues/12641)

# Download

All releases are

* [Releases · manabu/dockerlayer](https://github.com/manabu/dockerlayer/releases)

## Latest Stable release

* [Release 0.1.3 · manabu/dockerlayer](https://github.com/manabu/dockerlayer/releases/tag/0.1.3)

## Latest Dev version

* [Release 0.1.4-dev · manabu/dockerlayer](https://github.com/manabu/dockerlayer/releases/tag/0.1.4-dev)




# Before Use

I need more investigate about Image format and others.

I develop with comparing with ***docker history*** command.

And I assume Image format v1

# Usage

```
docker pull manabu/dockerlayer:0.1.1
go run dockerlayer.go manabu/dockerlayer:0.1.1
```

## Sample output

```
0c0915d4ce93  /bin/sh -c echo "addenv" >> /fake.txt
0c0915d4ce93  /bin/sh -c #(nop)  ENV DUMMYVALUE=dummy
0c0915d4ce93  /bin/sh -c #(nop)  CMD ["/bin/bash"]
C fake.txt 7 0():0() 100644
79db2ad5a477  /bin/sh -c echo "fake" > /fake.txt
79db2ad5a477  /bin/sh -c #(nop)  CMD ["/bin/dash"]
79db2ad5a477  /bin/sh -c #(nop)  CMD ["/bin/sh"]
A fake.txt 5 0():0() 100644
ee6a5d644ec9  /bin/sh -c rm dummy.txt
D dummy.txt -6 0():0() 0
bc0c2c4963c5  /bin/sh -c echo "hello" > /dummy.txt
bc0c2c4963c5  /bin/sh -c #(nop) CMD ["/bin/bash"]
A dummy.txt 6 0():0() 100644
1fb1c66fac26  /bin/sh -c sed -i 's/^#\s*\(deb.*universe\)$/\1/g' /etc/apt/sources.list
C etc/ 0 0():0() 40755
C etc/apt/ 0 0():0() 40755
C etc/apt/sources.list -12 0():0() 100644
fd4f25b1c446  /bin/sh -c rm -rf /var/lib/apt/lists/*
C var/ 0 0():0() 40755
C var/lib/ 0 0():0() 40755
C var/lib/apt/ 0 0():0() 40755
C var/lib/apt/lists/ 0 0():0() 40755
D var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_trusty_Release -58512 0():0() 100444
D var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_trusty_Release.gpg -933 0():0() 100444
D var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_trusty_main_binary-amd64_Packages -8234934 0():0() 100444
D var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_trusty_main_i18n_Translation-en -4149211 0():0() 100444
D var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_trusty_restricted_binary-amd64_Packages -184141 0():0() 100444
D var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_trusty_restricted_i18n_Translation-en -21217 0():0() 100444
D var/lib/apt/lists/lock 0 0():0() 100444
D var/lib/apt/lists/partial 0 0():0() 100444
7a46bd958bc8  /bin/sh -c set -xe                && echo '#!/bin/sh' > /usr/sbin/policy-rc.d     && echo 'exit 101' >> /usr/sbin/policy-rc.d     && chmod +x /usr/sbin/policy-rc.d               && dpkg-divert --local --rename --add /sbin/initctl     && cp -a /usr/sbin/policy-rc.d /sbin/initctl    && sed -i 's/^exit.*/exit 0/' /sbin/initctl             && echo 'force-unsafe-io' > /etc/dpkg/dpkg.cfg.d/docker-apt-speedup             && echo 'DPkg::Post-Invoke { "rm -f /var/cache/apt/archives/*.deb /var/cache/apt/archives/partial/*.deb /var/cache/apt/*.bin || true"; };' > /etc/apt/apt.conf.d/docker-clean   && echo 'APT::Update::Post-Invoke { "rm -f /var/cache/apt/archives/*.deb /var/cache/apt/archives/partial/*.deb /var/cache/apt/*.bin || true"; };' >> /etc/apt/apt.conf.d/docker-clean   && echo 'Dir::Cache::pkgcache ""; Dir::Cache::srcpkgcache "";' >> /etc/apt/apt.conf.d/docker-clean              && echo 'Acquire::Languages "none";' > /etc/apt/apt.conf.d/docker-no-languages          && echo 'Acquire::GzipIndexes "true"; Acquire::CompressionTypes::Order:: "gz";' > /etc/apt/apt.conf.d/docker-gzip-indexes               && echo 'Apt::AutoRemove::SuggestsImportant "false";' > /etc/apt/apt.conf.d/docker-autoremove-suggests
C etc/ 0 0():0() 40755
C etc/apt/ 0 0():0() 40755
C etc/apt/apt.conf.d/ 0 0():0() 40755
A etc/apt/apt.conf.d/docker-autoremove-suggests 44 0():0() 100644
```

## about filter

If you want to filter with file or directory name

```
go run dockerlayer.go sources
```

### current spec

1. Match is executed by ```regexp.MatchString```
2. ```/var/run/docker``` is stored with ```var/run/docker``` in tar file.


## about .wh. file

* [docker/v1.md at master · docker/docker](https://github.com/docker/docker/blob/master/image/spec/v1.md)

* [distribution/manifest-v2-2.md at master · docker/distribution](https://github.com/docker/distribution/blob/master/docs/spec/manifest-v2-2.md)

* [opencontainers/image-spec: OCI Image Format](https://github.com/opencontainers/image-spec)

* [Layering of Docker images – Thomas Uhrig](http://tuhrig.de/layering-of-docker-images/)


# Discussion about docker diffi

* [diffi command to inspect changes on image filesystem by ashwinphatak · Pull Request #12919 · docker/docker](https://github.com/docker/docker/pull/12919)

* [[Proposal]: docker diff between image layers · Issue #12641 · docker/docker](https://github.com/docker/docker/issues/12641)

# Windows building error


To fix this

```
--> windows/amd64 error: exit status 1
Stderr: ../../fsouza/go-dockerclient/client_windows.go:16:2: cannot find package "github.com/Microsoft/go-winio" in any of:
	/goroot/src/github.com/Microsoft/go-winio (from $GOROOT)
	/gopath/src/github.com/Microsoft/go-winio (from $GOPATH)
../../spf13/cobra/command_win.go:9:2: cannot find package "github.com/inconshreveable/mousetrap" in any of:
	/goroot/src/github.com/inconshreveable/mousetrap (from $GOROOT)
	/gopath/src/github.com/inconshreveable/mousetrap (from $GOPATH)
```

Get

```
go get github.com/inconshreveable/mousetrap
go get -d github.com/Microsoft/go-winio
```


# TODO

- [ ] Investigate image format
- [ ] Improve output format
- [ ] Build command
- [ ] Add test
- [ ] Support JSON format
