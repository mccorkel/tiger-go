# Cursor Rules

## Protected Files
The following files should not be modified without explicit permission:

1. `go.mod`
   - All dependency changes must be approved
   - Version updates require discussion
   - Direct/indirect status changes need approval

2. `frontend/package.json`
   - All npm package changes must be approved
   - Version updates require discussion
   - Dev dependencies changes need approval

## Project Defaults
1. AWS Region: `us-east-1`
   - All AWS services should default to us-east-1
   - Region changes require discussion and approval

## Reason
These files control project dependencies and their versions. Uncontrolled modifications could lead to:
- Incompatible package versions
- Breaking changes
- Security vulnerabilities
- Build failures

## Process
When changes are needed:
1. Request permission
2. Show the exact changes
3. Wait for approval
4. Document changes in versions.md

## Version Documentation
When adding or updating packages:
1. All changes must be documented in versions.md
2. Include package name and version
3. Update both Frontend and Backend Dependencies sections as needed
4. Add relevant notes about the change 