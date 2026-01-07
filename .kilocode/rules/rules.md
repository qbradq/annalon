# Agent Rules

## MCP Server Usage

You have access to multiple MCP services. Use the tools they provide. Keep to
the following rules.

### Context 7

Always prefer Context 7 documentation over training data when possible. Do not
wait for me to say "use context7", just use it anytime you need package
documentation. Use training data as a last resort.

### GoPls

- Do not use `go build` or similar commands to check for errors. Use `gopls`
  MCP tool `go_diagnostics` instead.
- Prefer `gopls` MCP tool `go_search` for finding symbols in the workspace.
- Prefer `gopls` MCP tool `go_symbol_references` for searching for symbol
  references.
- Use `gopls` MCP tool `go_workspace` to get a summary of the Go workspace.
