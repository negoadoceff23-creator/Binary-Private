#!/data/data/com.termux/files/usr/bin/bash
clear

echo "========================================="
echo "        BINARY BYPASS - ADB CONNECT"
echo "========================================="
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

echo ""
echo "[⚙] Baixando BinaryBypass do GitHub..."
curl -L -s -o binarybypass https://raw.githubusercontent.com/negoadoceff23-creator/Binary-Private/main/binarybypass
chmod +x ./binarybypass

echo ""
echo "[🚀] Iniciando BinaryBypass..."
sleep 1
./binarybypass
