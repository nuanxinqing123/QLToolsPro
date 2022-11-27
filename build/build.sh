# xgo -out QLToolsPro --targets=linux/amd64 .
xgo -out QLToolsPro --targets=windows/*,linux/* ../
upx QLToolsPro-*