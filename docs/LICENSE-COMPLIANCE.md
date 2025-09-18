# License Compliance

This document explains the license compliance scanning setup for the Coruscant project.

## Overview

The project uses **ScanCode Toolkit** for comprehensive license compliance scanning to ensure all dependencies and code files comply with open source licensing requirements.

## Scanning Coverage

### What Gets Scanned
- **Source Code**: All Go modules, configuration files, scripts
- **Dependencies**: Go module dependencies (go.mod/go.sum analysis)
- **Infrastructure**: Dockerfiles, GitHub Actions workflows, YAML configs
- **Documentation**: README files, license files, copyright notices

### Report Formats
- **SPDX JSON**: Industry-standard machine-readable format
- **HTML**: Human-readable reports for manual review
- **CSV**: Tabular format for spreadsheet analysis

## Automated Scanning

### GitHub Actions Workflow
The `license-compliance.yml` workflow runs automatically on:
- Pull requests to main branch
- Pushes to main branch
- Weekly schedule (Saturdays at 4 AM UTC)
- Manual workflow dispatch

### Scan Categories
1. **Repository-wide**: Complete repository scan with detailed reports
2. **Module-specific**: Focused scans for Go modules (golib, kyber)
3. **Infrastructure**: Docker and GitHub Actions configurations

## Local Development

### Quick Commands
```bash
# Run comprehensive license scan
task license:scan

# Run fast license-only scan
task license:scan:fast

# Scan specific modules
task license:scan:golib
task license:scan:kyber

# Clean up scan results
task license:clean
```

### Prerequisites
- Python 3.9+ (for ScanCode Toolkit)
- pip (Python package manager)

## Compliance Workflow

### 1. Review Scan Results
- Download artifacts from GitHub Actions runs
- Check HTML reports for human-readable summaries
- Analyze SPDX JSON files for detailed license data

### 2. License Compatibility
- Verify all detected licenses are compatible with project requirements
- Check for any GPL or other copyleft licenses in dependencies
- Ensure proper attribution requirements are met

### 3. Update Documentation
- Keep LICENSE and NOTICE files updated
- Add third-party attributions as required
- Update dependency documentation

### 4. Policy Compliance
- Follow organization's open source policy
- Document any exceptions or special cases
- Maintain compliance records

## Common License Types

### Typically Acceptable
- MIT License
- Apache License 2.0
- BSD licenses (2-clause, 3-clause)
- ISC License

### Review Required
- GPL licenses (may require source disclosure)
- LGPL licenses (linking restrictions)
- Custom or proprietary licenses

### Typically Not Acceptable
- Copyleft licenses in proprietary projects
- Licenses with commercial use restrictions
- Unclear or missing license information

## Troubleshooting

### ScanCode Installation Issues
```bash
# Update pip and install with all features
pip install --upgrade pip
pip install scancode-toolkit[full]
```

### Large Repository Scans
- Use `license:scan:fast` for quicker license-only scanning
- Focus on specific modules with targeted scans
- Exclude unnecessary directories (vendor/, node_modules/)

### False Positives
- Review scan results manually for accuracy
- Check license detection confidence scores
- Verify ambiguous license detections

## Integration with CI/CD

The license compliance pipeline integrates with existing CI/CD:
- Runs independently (no build blocking)
- Uploads detailed reports as artifacts
- Provides compliance summary for review
- Supports both automated and manual review workflows

## Resources

- [ScanCode Toolkit Documentation](https://scancode-toolkit.readthedocs.io/)
- [SPDX License List](https://spdx.org/licenses/)
- [Open Source License Compatibility](https://www.gnu.org/licenses/license-compatibility.html)
- [GitHub License Guide](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/licensing-a-repository)