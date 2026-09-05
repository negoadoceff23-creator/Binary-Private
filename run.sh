#!/bin/sh
clear

echo "========================================="
echo "        BINARY BYPASS - ADB CONNECT"
echo "========================================="
echo ""

# 1. Checar se já existe um dispositivo ADB conectado
CONNECTED=$(adb devices 2>/dev/null | grep -v "List" | grep "device$" | head -n 1 | awk '{print $1}')

if [ -n "$CONNECTED" ]; then
    echo "[✓] Dispositivo ADB ja conectado/pareado: $CONNECTED"
else
    echo "[?] Nenhum dispositivo conectado. Vamos parear:"
    echo ""
    printf "IP:PORTA de pareamento: "
    read PAIR_ADDR
    printf "Codigo de pareamento: "
    read PAIR_CODE

    echo ""
    echo "[⚙] Pareando..."
    adb pair "$PAIR_ADDR" "$PAIR_CODE"

    echo ""
    printf "IP:PORTA de conexao: "
    read CONNECT_ADDR

    echo ""
    echo "[⚙] Conectando..."
    adb connect "$CONNECT_ADDR"
fi

echo ""
echo "[⚙] Baixando BinaryBypass do GitHub..."

# Baixa com failover e valida se o binario veio correto
curl -fL -k -o binarybypass "https://raw.githubusercontent.com/negoadoceff23-creator/Binary-Private/main/binarybypass"

if [ $? -eq 0 ] && [ -s binarybypass ]; then
    chmod +x ./binarybypass
    echo "[✓] Download concluido!"
    echo ""
    echo "[🚀] Iniciando BinaryBypass..."
    sleep 1
    ./binarybypass
else
    echo "[!] Falha ao baixar do GitHub."
    echo "[⚙] Tentando executar versao local existente..."
    if [ -f ./binarybypass ]; then
        chmod +x ./binarybypass
        ./binarybypass
    else
        echo "[✗] Erro: binarybypass nao encontrado localmente."
    fi
fi
