# openEuler_1w

华为openEuler 1w+ 软件包迁移项目
## apps-辅助快速测试的脚本

可以在[src](./apps/src)目录下找到这些脚本的源码，重新编码后通过交叉编译构建出arrch64下的可执行文件。或者可以直接在[bin](./apps/bin)找到可执行文件

1. autosumbit
- 功能:  
通过传入参数，自动修改job.yaml并且自动生成日志文件，映射repo与job id
- 用法:
```shell
./auto-submit job.yaml repo_addr
```
- 例子:
```shell
./auto-submit jobs/wq.yaml https://opentuna.cn/epel/7/source/tree/Packages/b/bodhi-2.11.0-3.el7.src.rpm
```


2. 正在赶来的路上..........
## srcrpm-欧拉rpm包反编译镜像站
增加反编译镜像包，直接通过pr提交，提交并审核通过后会，我会更新到下面的服务器。
http://39.104.160.208:30984/openEuler_1w/srcrpm/
