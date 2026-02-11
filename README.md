<div align="center">
  
# Atomic-Mirai

**Atomic-Mirai is a variant of the Condi-Mirai botnet family, designed for educational and research purposes in understanding IoT security vulnerabilities and DDoS attack vectors.**

  <img width="714" height="420" alt="image" src="https://github.com/user-attachments/assets/d4e83fa7-a9f9-4a22-b424-bdfc2eebd89c" />
</div>


## Build Instructions

### Requirements

- Ubuntu 20.04 or later
- GCC cross-compilers for target architectures
- Go 1.20+
- SQLite3
- Nginx (for serving binaries)

### Compilation

```bash
# Compile bot for all architectures
chmod +x build.sh
./build.sh

# Compile C&C server
go build cnc/*.go
```

### Configuration

1. Configure C&C domain/IP in `bot/table.c` using the XOR tool
2. Set server settings in `assets/server.toml`
3. Initialize database: `sqlite3 assets/data.sql`
4. Place compiled binaries in web server directory

## Disclaimer

**EDUCATIONAL PURPOSE ONLY**

This software is provided strictly for educational, research, and security testing purposes. The authors and contributors:

- Do NOT condone or encourage any illegal activity
- Are NOT responsible for any misuse or damage caused by this software
- Assume NO liability for actions taken by users of this software

**WARNING**: Unauthorized access to computer systems, creation of botnets, and DDoS attacks are ILLEGAL in most jurisdictions and can result in severe criminal penalties including:

#

<div align="center">

[**@Atomic Team**](https://github.com/AtomTM) <br> [**@CirqueiraDev**](https://github.com/CirqueiraDev) 
</div>
