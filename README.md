# terraform-provider-dare

Data at rest encryption terraform provider

## Dare binary

```bash
bin/dare encrypt -input=./foo -master-key=$(bin/dare keygen)
```

## Env vars

- `DARE_DEMO_SIZE`
- `DARE_DEMO_PATH`
- `DARE_MASTER_KEY`