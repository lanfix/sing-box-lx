# TASKS: 113 — NEIGHBOR_RULE_SET

- [x] 1. FEATURE 018-NEIGHBOR_RULE_SET, SPEC / PLAN
- [x] 2. `route/rule/rule_headless_neighbor_lx.go`: матчеры MAC/имени, отказ без `find_neighbor`
- [x] 3. `option/rule_set.go`: поля под маркером
- [x] 4. `route/rule/rule_headless.go`: вызов под маркером
- [x] 5. `common/srs/binary.go`: отказ компиляции под маркером
- [x] 6. Тесты `route/rule/rule_headless_neighbor_lx_test.go`, `common/srs/binary_neighbor_lx_test.go`
- [x] 7. `lx-test/config/neighbor_rule_set.json` для `lx-check`
- [x] 8. `go vet`, `go test` с тегами lx; `make -f Makefile.lx lx-build`; `sing-box check`; `rule-set compile`
- [x] 9. e2e в docker-сети: два клиента с заданными MAC, смена `local` rule-set без перезапуска
- [x] 10. Доки: §14 в `docs-lx/lx-config{,.ru}.md`; строки в `SPECS/README.md`, `SPECS/FEATURES/README.md`, watchlist 005
- [ ] 11. Полевой прогон на роутере с конфигуратором
- [ ] 12. PR в `SagerNet/sing-box` (`testing`) с полной формой
