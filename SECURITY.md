Reporting a Vulnerability
If you discover a security vulnerability in CDLytics, please do not open a public GitHub issue.
Instead, report it privately by emailing:
sujuzin@proton.me
Please include:

A description of the vulnerability
Steps to reproduce
Potential impact
Any suggested fixes if applicable

You can expect an acknowledgment within 48 hours and a resolution timeline within 7 days depending on severity.
Scope
The following are in scope for security reports:

Authentication and session management
Account data exposure or unauthorized access
SQL injection or database vulnerabilities
Cross-site scripting (XSS)
Cross-site request forgery (CSRF)
Insecure direct object references
API endpoint vulnerabilities
AWS infrastructure misconfigurations

The following are out of scope:

Denial of service attacks
Social engineering
Issues in third-party dependencies outside our control

Security Practices
CDLytics follows these security practices:

All user passwords are hashed and never stored in plaintext
Authentication tokens are short-lived and rotated
Database access is restricted via AWS RDS security groups
All traffic is served over HTTPS via CloudFront
Environment variables and secrets are never committed to source control
AWS IAM roles follow the principle of least privilege

Disclosure Policy
We follow a coordinated disclosure policy. Once a vulnerability is resolved, we are happy to credit the reporter in our release notes if they wish.
