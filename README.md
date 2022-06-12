# openEuler_1w

华为openEuler 1w+ 软件包迁移项目
## apps-辅助快速测试脚本

可以在[src](./apps/src)目录下找到这些脚本的源码，重新编码后通过交叉编译构建出arrch64下的可执行文件。或者可以直接在[bin](./apps/bin)找到可执行文件

1. autosumbit
- 功能:  
指定`job.yaml`文件，自动修改`job.yaml`的`repo_addr`字段并且自动生成日志文件，映射repo与job id。  
生成日志文件的名字为 `job.yaml.log`
- 用法:  
检查自己的job.yaml与[fmt.yaml](./apps/src/autosubmit/fmt.yaml)的字段是否匹配（**值可以根据情况填写**）。然后运行下面的小程序
```shell
./auto-submit job.yaml
请输入 repo_addr: repo_addr1
请输入 repo_addr: repo_addr2
# ctrl+c 退出
.....
```
查看日志文件，进行提交任务的复查比对。
```shell
cat job.yaml.log
```

2. 正在赶来的路上..........
## srcrpm-欧拉rpm包反编译镜像站
增加反编译镜像包，直接通过pr提交，提交并审核通过后会，我会更新到下面的服务器。
http://39.104.160.208:30984/openEuler_1w/srcrpm/
