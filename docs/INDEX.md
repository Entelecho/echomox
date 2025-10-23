# EchoMox Documentation Index

Welcome to the EchoMox documentation! This index will help you find the information you need.

## 🚀 Getting Started

**New to EchoMox?** Start here:

1. **[README.md](../README.md)** - Project overview, features, and quick start guide
2. **[Deployment Guide](DEPLOYMENT.md)** - Complete installation and configuration instructions
3. **[Quick Start Tutorial](#quick-start-tutorial)** - 10-minute setup guide

## 📖 Core Documentation

### Architecture & Design

- **[Technical Architecture](ARCHITECTURE.md)** ⭐ *Recommended*
  - Complete system architecture with mermaid diagrams
  - Component interactions and data flows
  - Email processing pipeline
  - Security and deployment architectures
  - ~1,500 lines with 20+ diagrams

### Deployment & Operations

- **[Deployment Guide](DEPLOYMENT.md)**
  - Prerequisites and system requirements
  - Installation methods (binary, Docker, from source)
  - DNS configuration with examples
  - Monitoring and maintenance
  - Troubleshooting guide
  - Security best practices

### Reservoir Computing

- **[Reservoir Computing Framework](../RESERVOIR_COMPUTING.md)**
  - Overview of the AI/ML framework
  - Echo State Networks (ESN) details
  - Membrane Computing (P-Systems)
  - Affective Computing and emotion detection
  - Mathematical foundations
  - Performance characteristics

- **[Integration Guide](../reservoir/INTEGRATION.md)**
  - How to integrate reservoir computing with existing code
  - Configuration examples
  - Training the ESN
  - Advanced usage patterns
  - Performance tuning

- **[Package Documentation](../reservoir/README.md)**
  - API reference for the reservoir package
  - Code examples
  - Configuration options

## 🎯 Use Case Guides

### For System Administrators

1. Start with [Deployment Guide](DEPLOYMENT.md)
2. Review [Network Architecture](ARCHITECTURE.md#network-architecture)
3. Set up [Monitoring](DEPLOYMENT.md#monitoring-and-maintenance)
4. Configure [Security](ARCHITECTURE.md#security-architecture)

### For Developers

1. Read [Technical Architecture](ARCHITECTURE.md)
2. Explore [Core Components](ARCHITECTURE.md#core-components)
3. Review [Package Documentation](../reservoir/README.md)
4. Check [Integration Guide](../reservoir/INTEGRATION.md)

### For ML/AI Researchers

1. Deep dive into [Reservoir Computing](../RESERVOIR_COMPUTING.md)
2. Study [Echo State Network Architecture](ARCHITECTURE.md#echo-state-network-esn-details)
3. Understand [Membrane Computing](ARCHITECTURE.md#membrane-computing-p-systems)
4. Review [Affective Computing Integration](ARCHITECTURE.md#affective-computing-integration)

## 📊 Diagram Reference

All documentation includes mermaid diagrams for visual understanding:

### System Architecture Diagrams
- High-level system architecture
- Component architecture
- System context (C4 model)

### Processing Flow Diagrams
- Incoming email flow (sequence diagram)
- Outgoing email flow (sequence diagram)
- Spam classification pipeline (flowchart)

### Reservoir Computing Diagrams
- Reservoir computing architecture
- ESN architecture details
- Membrane computing hierarchy
- Affective computing integration

### Data Architecture Diagrams
- Database schema (ER diagram)
- File system layout
- Memory usage breakdown

### Network Diagrams
- Port and protocol layout
- TLS/security flow
- DNS setup visualization

### Deployment Diagrams
- Single server deployment
- High-availability architecture (planned)
- Docker deployment
- Performance tuning flowchart

## 🔗 Quick Links

### External Resources
- [GitHub Repository](https://github.com/Entelecho/echomox)
- [Upstream mox Project](https://github.com/mjl-/mox)
- [mox Website](https://www.xmox.nl)

### Related Documentation
- [develop.txt](../develop.txt) - Development guidelines
- [compatibility.txt](../compatibility.txt) - Compatibility information

## 📝 Documentation Standards

All documentation follows these standards:

- **Markdown Format**: GitHub-flavored markdown
- **Diagrams**: Mermaid.js for all technical diagrams
- **Code Examples**: Syntax-highlighted with language tags
- **Links**: Relative paths for internal links
- **Structure**: Clear hierarchy with table of contents

## 🆘 Getting Help

If you can't find what you're looking for:

1. Check the [Troubleshooting](DEPLOYMENT.md#troubleshooting) section
2. Search the [GitHub Issues](https://github.com/Entelecho/echomox/issues)
3. Review the [FAQ](../README.md#faq---frequently-asked-questions)
4. Join the community chat (#mox on irc.oftc.net)

## 📅 Document Versions

| Document | Version | Last Updated |
|----------|---------|--------------|
| ARCHITECTURE.md | 1.0 | 2025-10-23 |
| DEPLOYMENT.md | 1.0 | 2025-10-23 |
| RESERVOIR_COMPUTING.md | 1.0 | 2025-10-23 |
| README.md | 1.0 | 2025-10-23 |

---

## Quick Start Tutorial

### Prerequisites
```bash
# Check Go version
go version  # Should be >= 1.23

# Check system resources
free -h     # At least 512MB RAM recommended
df -h       # At least 10GB storage
```

### Installation (3 minutes)
```bash
# Create mox user
sudo useradd -m -d /home/mox mox

# Switch to mox user
sudo -u mox -i

# Download and compile
cd /home/mox
GOBIN=$PWD CGO_ENABLED=0 go install github.com/mjl-/mox@latest
```

### Configuration (5 minutes)
```bash
# Run quickstart
./mox quickstart you@example.com

# This will:
# 1. Generate config files
# 2. Create admin password
# 3. Show DNS records to add
# 4. Provide systemd service commands
```

### DNS Setup (2 minutes)
Add the DNS records shown by quickstart to your DNS provider.

### Start Server
```bash
# Install as service
sudo systemctl enable --now mox

# Check status
sudo systemctl status mox

# View logs
sudo journalctl -u mox -f
```

### Access Admin Interface
```bash
# SSH tunnel for security
ssh -L 8080:localhost:80 you@your-server

# Open in browser
http://localhost:8080/admin/
```

### Enable Reservoir Computing
Edit `/home/mox/config/domains.conf`:

```yaml
Accounts:
  youruser:
    ReservoirFilter:
      Enabled: true
      ReservoirWeight: 0.3
      ESNParams:
        ReservoirSize: 100
        SpectralRadius: 0.95
        LeakRate: 0.3
```

Restart the service:
```bash
sudo systemctl restart mox
```

That's it! You now have EchoMox running with reservoir computing enabled.

---

**Need more details?** See the [Deployment Guide](DEPLOYMENT.md) for complete instructions.
