go-download
===========

Chunked HTTP Download Manager

![Example](https://rawgithub.com/snikch/go-download/master/example.gif)

## About

This is a hobby project, not stable software (yet).

## Done

- [x] Basic fixed size chunked http downloading
- [x] Download configuration based on resource location
- [x] Resume downloads

## In Progress

- [ ] Authenticated providers
- [ ] Define process communication protocol
- [ ] Implement RPC over unix sockets
- [ ] Create CLI GUI


## TODO
- [ ] Split GUI(s) from core daemon
- [ ] Move core processing to daemonized process
- [ ] Customisable settings
- [ ] Settings persistence
- [ ] Monitor clipboard for url like things
- [ ] Define post processing pipeline via plugs (unzip, unrar, mp4 to itunes etc )