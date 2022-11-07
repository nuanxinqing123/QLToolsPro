xgo -out QLToolsPro_1.0 --targets=windows/*,linux/* ../
rm -rf QLToolsPro_1.0-linux-riscv64
rm -rf QLToolsPro_1.0-linux-s390x
rm -rf QLToolsPro_1.0-linux-mips64
rm -rf QLToolsPro_1.0-linux-mips64le
upx QLToolsPro_1.0-*