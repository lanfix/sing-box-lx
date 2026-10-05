# IMPLEMENTATION_REPORT: 113 — NEIGHBOR_RULE_SET

## Файлы

| Род | Файл |
|-----|------|
| new | `route/rule/rule_headless_neighbor_lx.go` |
| new | `route/rule/rule_headless_neighbor_lx_test.go` |
| new | `common/srs/binary_neighbor_lx_test.go` |
| new | `lx-test/config/neighbor_rule_set.json` |
| `// lx:` upstream | `option/rule_set.go` (+4), `route/rule/rule_headless.go` (+6), `common/srs/binary.go` (+5) |
| deps | — |
| wiring | — (без build-тега) |
| docs | `docs-lx/lx-config{,.ru}.md` §14, `SPECS/` |

## DoD

```sh
make -f Makefile.lx lx-build LX_VERSION=1.14.2-lx.11-mac.1-dev   # go1.26.8, полный набор тегов — ок
go vet  -tags "$(make -f Makefile.lx -s lx-print-tags)" ./route/... ./common/srs/ ./option/   # ок
go test -tags "$(make -f Makefile.lx -s lx-print-tags)" -ldflags=-checklinkname=0 \
  ./route/... ./common/srs/ ./option/ ./cmd/sing-box/                                         # ок
./sing-box check -c lx-test/config/neighbor_rule_set.json                                     # ок
```

- Без `route.find_neighbor`: `check` →
  ``initialize router: parse rule-set[0]: parse rule_set.rules.[0]: `source_mac_address` and `source_hostname` rule items in rule-set require `route.find_neighbor` ``.
- `rule-set compile` → `` `source_mac_address` and `source_hostname` rule items are only supported in source rule-set ``.

## e2e (docker, сеть `172.30.0.0/24`)

Ядро `172.30.0.2` (mixed-inbound, `find_neighbor: true`, `local` rule-set `devices` → `reject`,
`final: direct`), клиенты `c1` `02:00:00:00:00:0a` и `c2` `02:00:00:00:00:0b`, запросы через
HTTP-прокси.

| Фаза | rule-set | c1 | c2 |
|------|----------|----|----|
| 1 | `02:00:00:00:00:0a` | `found neighbor …:0a` → `rule_set=devices => reject` | `found neighbor …:0b` → direct |
| 2 (правка файла, тот же процесс) | `02-00-00-00-00-0B` | direct | `rule_set=devices => reject` |

В журнале одна строка `sing-box started`. Перечитывание с MAC без `find_neighbor` →
`router: reload rule-set devices: … require route.find_neighbor`, прежний набор остаётся.

## Зона конфликтов при мерже

`option/rule_set.go` (`DefaultHeadlessRule`), `route/rule/rule_headless.go`
(`NewDefaultHeadlessRule`, перед AdGuard-блоком), `common/srs/binary.go` (`writeDefaultRule`, начало).

## Вне скоупа

Двоичный `.srs`, версия rule-set, автозапуск поиска соседей — в апстрим-PR. Полевой прогон на роутере.
