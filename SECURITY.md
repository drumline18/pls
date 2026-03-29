# Security Policy

## Supported versions

Until the project has tagged stable releases, security fixes land on the main development branch.

Once public releases begin, supported versions should be documented here.

## Reporting a vulnerability

For now, please report vulnerabilities privately to the maintainer instead of opening a public issue.

When the public repo exists, replace this section with the preferred private contact path.

## Scope

Because `pls` suggests and can optionally execute shell commands, security-sensitive areas include:
- command generation rules
- risk classification
- confirmation gating
- execution paths
- config loading and precedence
- provider credential handling

Please avoid posting weaponized proofs of concept in public issues.
