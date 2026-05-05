# patchwork

Lightweight utility to apply and track ad-hoc patches across multiple git repos.

---

## Installation

```bash
go install github.com/yourusername/patchwork@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/patchwork.git && cd patchwork && go build -o patchwork .
```

---

## Usage

Define your repos and patches in a `patchwork.yaml` file, then apply them all at once:

```yaml
repos:
  - path: ../service-a
  - path: ../service-b

patches:
  - name: fix-logging
    file: patches/fix-logging.patch
```

```bash
# Apply all patches across all repos
patchwork apply

# Check which patches have been applied
patchwork status

# Roll back a specific patch
patchwork revert --patch fix-logging
```

Patchwork tracks applied patches in a local `.patchwork` state file so you never apply the same patch twice.

---

## Commands

| Command         | Description                              |
|-----------------|------------------------------------------|
| `apply`         | Apply all pending patches to all repos   |
| `status`        | Show patch status across repos           |
| `revert`        | Revert a named patch from all repos      |
| `validate`      | Dry-run to check if patches apply cleanly|

---

## License

MIT © 2024 yourusername