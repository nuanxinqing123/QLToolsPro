Docker 启动命令
```shell
docker run --restart=always -itd --name QLToolsPro -v $PWD/config:/QLToolsPro/config -v $PWD/logs:/QLToolsPro/logs -v $PWD/plugin:/QLToolsPro/plugin -p 6600:6600 nuanxinqing123/qltools-pro:latest
```