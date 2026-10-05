# PLAN: 113 — NEIGHBOR_RULE_SET

## Решение

Минимальный дифф: headless-правило получает два поля и собирает их апстримными матчерами.
Поиск соседей не трогаем — он создаётся в `Router.Start` на `StartStateInitialize`, раньше, чем
грузятся rule-set (`StartStateStart`), и rule-set не может его запустить. Вместо правки роутера —
требование `route.find_neighbor` с громким отказом.

```go
// route/rule/rule_headless_neighbor_lx.go
func appendNeighborHeadlessItems(ctx, rule, options) error {
	if нет полей { return nil }
	router := service.FromContext[adapter.Router](ctx)
	if router != nil && !router.NeedFindNeighbor() { return E.New("… require `route.find_neighbor`") }
	rule.items += NewSourceMACAddressItem / NewSourceHostnameItem
}
```

Роутер регистрируется в контексте до `Router.Initialize` (`box.go`), поэтому проверка работает
и для `inline`/`local` rule-set, которые разбираются в `Initialize`, и для `remote`, которые
разбираются на старте и при обновлении.

## Точки касания upstream (все под `lx:begin neighbor-rule-set`)

| Файл | Что | Почему нельзя иначе |
|------|-----|---------------------|
| `option/rule_set.go` | поля `SourceMACAddress`, `SourceHostname` в `DefaultHeadlessRule` | структура опций |
| `route/rule/rule_headless.go` | вызов `appendNeighborHeadlessItems` в `NewDefaultHeadlessRule` | сборщик правила |
| `common/srs/binary.go` | отказ в `writeDefaultRule` | иначе `rule-set compile` молча выбросит условие |

`DefaultHeadlessRule.IsValid` и `isNonIPCIDRHeadlessRule` сравнивают через `reflect.DeepEqual`
с пустым значением — новые поля учитываются без правок.

## Отвергнуто

- **Двоичная запись и версия 6.** Номера `ruleItem*` и `RuleSetVersion*` выдаёт апстрим; занять их
  у себя — разойтись с будущим апстримным `.srs`. Конфигуратор использует формат `source`.
- **Автозапуск поиска соседей** (`ContainsNeighborRule` в метаданных + дозапуск в `Router.Start`):
  ещё два апстримных файла и логика старта роутера. Для апстрим-PR — да, для форка — нет.
- **Сопоставление MAC → IP на стороне конфигуратора** (`source_ip_cidr` в rule-set): работает без
  ядра, но держит вторую копию таблицы соседей, требует доступа к ней из контейнера конфигуратора
  и даёт окно, когда IP ушедшего устройства достаётся новому.

## Риски

- Мерж апстрима, переписавший `NewDefaultHeadlessRule` или `writeDefaultRule`: маркеры переносятся
  руками, тесты 113 ловят потерю.
- Апстрим введёт поля с тем же JSON-именем, но другой семантикой — тогда это DROP по условию снятия.
