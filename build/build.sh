# xgo -out QLToolsPro --targets=linux/amd64 .
xgo -out QLToolsPro --targets=linux/amd64,linux/arm64,linux/arm-7 ../
upx QLToolsPro-*