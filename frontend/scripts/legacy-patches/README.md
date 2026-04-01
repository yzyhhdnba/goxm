# Legacy Patch Scripts

This folder stores one-off JavaScript patch scripts that were used during frontend API migration and quick local fixes.

These scripts are **not** part of the normal build/runtime path.

## Why moved here

- Keep the repository root clean.
- Group frontend-only maintenance scripts with the frontend project.
- Avoid confusion between runnable project entrypoints and temporary patch utilities.

## Notes

- Most scripts do regex-based file rewriting.
- Re-running old scripts may overwrite newer code changes.
- Prefer updating source files directly for future changes.
