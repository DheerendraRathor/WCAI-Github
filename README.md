# WCAI-Github

A Github parser for [WCAI](https://github.com/DheerendraRathor/WCAI) Project


## Installation
1. Install Go 1.8.3
2. Install Go dependency manager [dep](https://github.com/golang/dep) or [vg](https://github.com/GetStream/vg) (Whatever suits your buds)
3. Run `dep ensure` or `vg ensure`
4. Copy `config.sample.json` to `config.json` and modify `config.json` with appropriate values


## Usages
1. Run `go run main.go --help`
    - Ex: `go run main.go --update=lang`  (Possible values for update are: fork, lang, views, clones, fetch, limit)
