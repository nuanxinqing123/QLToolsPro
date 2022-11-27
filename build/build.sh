# xgo -out QLToolsPro --targets=linux/amd64 .
xgo -out QLToolsPro --targets=windows/*,linux/amd64,linux/arm64,linux/amd64 ../
upx QLToolsPro-*