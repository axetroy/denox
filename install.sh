#!/usr/bin/env bash
set -e

downloadFolder="${HOME}/Downloads"

mkdir -p downloadFolder

get_arch() {
    a=$(uname -m)
    case ${a} in
    "x86_64" | "amd64" )
        echo "amd64"
        ;;
    "i386" | "i486" | "i586")
        echo "386"
        ;;
    *)
        echo ${NIL}
        ;;
    esac
}

get_os(){
    echo $(uname -s | awk '{print tolower($0)}')
}

main() {
    local os=$(get_os)
    local arch=$(get_arch)
    local dest_file="${downloadFolder}/denox_${os}_${arch}.tar.gz"

    if [[ $# -eq 0 ]]; then
        asset_path=$(
            command curl -sSf https://github.com/axetroy/denox/releases |
                command grep -o "/axetroy/denox/releases/download/.*/denox_${os}_${arch}\\.tar.gz" |
                command head -n 1
        )
        if [[ ! "$asset_path" ]]; then exit 1; fi
        asset_uri="https://github.com${asset_path}"
    else
        asset_uri="https://github.com/axetroy/denox/releases/download/${1}/denox_${os}_${arch}\\.tar.gz"
    fi

    mkdir -p ${downloadFolder}

    echo "[1/3] Download ${asset_uri} to ${downloadFolder}"
    rm -f ${dest_file}
    wget -P "${downloadFolder}" "${asset_uri}"

    echo "[2/3] Install denox to the ${HOME}/bin"
    mkdir -p ${HOME}/bin
    tar -xz -f ${dest_file} -C ${HOME}/bin
    exe=${HOME}/bin/denox
    chmod +x ${exe}

    echo "[3/3] Set environment variables"
    echo "Denox was installed successfully to ${exe}"
    if command -v denox >/dev/null; then
        echo "Run 'denox https://deno.land/std/examples/welcome.ts' to get started"
    else
        echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
        echo "  export PATH=${HOME}/bin:\$PATH"
        echo "Run '$exe --help' to get started"
    fi

    exit 0
}

main