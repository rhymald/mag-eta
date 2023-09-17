# Welcome!

> MAG-eta started as 7th attempt!

1. Alpha: Not published - primitives
2. Beta: Not published - primitives, fighting mechanics
3. Gamma: [here](https://github.com/rhymald/mag-gamma/tree/MBF-elemental-state-refactoring) - primitives and character
4. Delta: [here](https://github.com/rhymald/mag-delta/tree/N33-player-refactoring) - fighting mechanics, block tree
5. Epsilon: [here](https://github.com/rhymald/mag-epsilon/tree/N3G-character) - interactive CLI, trying transactional
6. Zeta: [here](https://github.com/rhymald/mag-zeta/tree/N7S-world) - successfully transactional, with movements across world grid
7. Eta: current repo - will be transactional, block tree

# Tests: 

|When|Resources|Total objects|Read|Comment|
|:-|:-:|:-:|:-:|-:|
|N9G: 1.5-D index, no panic on empty trace|4C 6G|32000|20ms|Unknown|
|N9G: 1.5-D index, no panic on empty trace|4C 6G|51725|40ms|Some objects cant move, - not enough CPU|
|N9G: 1.5-D index, anic on empty trace|4C 6G|7675-7862|No trace|CPU used for 335%|
|N9G: 1.5-D index, anic on empty trace|7C 1G|???|No trace|CPU used for 335%|

# How-to

List all funcs and types:
```bash
grep -r "^\(func\)\|^\(type\)" . | grep Dot
```

Delete and cleanups: 
```bash
sudo docker rm $(sudo docker ps -a -f status=exited -q) && sudo docker rmi $(sudo docker images -a -q)
sudo sh -c 'truncate -s 0 /var/lib/docker/containers/*/*-json.log'
set -eu ; LANG=en_US.UTF-8 snap list --all | awk '/disabled/{print $1, $3}' | while read snapname revision; do ; snap remove "$snapname" --revision="$revision" ; done
```

## Build

Container: 

```bash
docker buildx build .
docker tag 33415c rhymald/mag:latest
docker push rhymald/mag:latest
```

## Run

```bash
sudo docker compose up --build
```
