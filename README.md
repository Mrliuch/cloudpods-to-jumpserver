# Cloudpods虚拟机同步至Jumpserver
该项目通过Cloudpods的Webhook能力实现创建和删除虚拟机时通过Webhook发送信息至本程序，程序根据收到的事件信息判断是向Jumperserver中创建或者是删除对应的资产。从而实现资产的安全管理

# 依赖项
依赖Redis，用来存储Cloudpods中的虚拟机信息和Jumperserver中的资产信息的对应关系

# 启动服务
1、修改config.yaml文件，配置Redis服务器地址，以及Jumperserver的API地址等

2、node_id 为Jumperserver中资产节点的ID，虚拟机会自动同步至该节点下

# 后续优化
1、增加其他虚拟化平台或公有云私有云平台的Webhook

2、除Webhook外的其他渠道实现自动同步

3、支持自定义同步到自定义节点