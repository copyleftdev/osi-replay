# osi-replay

A modular Go project for:
- Capturing packets (OSI Layers 2+)
- Replaying PCAPs on a chosen interface
- Transforming / Sanitizing PCAPs
- Rewriting IP/MAC addresses

See `cmd/` subfolders for individual tools:
- `capture`
- `replay`
- `transform`
- `rewriter`

`pkg/` contains the reusable library logic.
