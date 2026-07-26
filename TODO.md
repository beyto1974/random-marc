# TODO

- [ ] Flag for record type/material (books vs. serials vs. scores) — would change 008 layout per type.
- [ ] More faithful 008 fixed-field per material type (illustrations, target audience, form-of-item codes currently left blank).
- [ ] `-lang` flag to vary the 008/041 language code instead of hardcoding `eng`.
- [ ] Fuzz test for `genmarc.Record` (`go test -fuzz`).
- [ ] golangci-lint config.
