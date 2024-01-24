# kubeutils项目的agent
用于协助kubeutils进行集群操作

# 构建可执行程序
```
go build -o kubeagent main.go
```

# Dockerfile
```
# 使用 alpine 作为基础镜像
FROM alpine:latest
#创建所需目录
RUN mkdir -p /home/kubeAgent/conf && mkdir /root/.kube/
# 将可执行程序复制到镜像中
COPY kubeagent /home/kubeAgent
RUN chmod +x /home/kubeAgent/kubeagent
# 将配置文件复制到镜像中的 conf 目录
COPY conf.ini /home/kubeAgent/conf/
# 设置工作目录为 /home/kubeServer
WORKDIR /home/kubeAgent
#暴露端口
EXPOSE 8081
# 指定容器启动时要运行的命令
CMD ["./kubeagent"]
```

# 配置文件
```
[server]
; 部署方式，有"docker"和"k8s"两种选择
model = docker
server_addr = 1.1.1.1   #kubeutils端地址和端口
port = 8999

[agent]
cluster_name = test
agent_addr = 2.2.2.2   #本地地址和端口（实际暴漏出去的）
port = 8081
```

# 部署方式
## docker部署
```
docker run -d --name kubeagent \
        -v conf.ini:/home/kubeAgent/conf/conf.ini \
        -v /root/.kube/config:/root/.kube/config \
        -p 8081:8081 \
        kubeagent:v1.0
```

## K8s部署
```
在kubeutilsweb页面进行点击部署
```
