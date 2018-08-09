从rms接口获取网络、服务器相关的信息。按照返回信息的类型组织模块。

文件列表
- common.go: baidu_rms package的公用资源定义
- vip.go: 查询vip信息
- rms_task.go: 查询rms任务单进度
- rms_token.go: 从open.rms.baidu.com的获取API使用所需的token
- vip_bns.go: 查询VIP的bns信息，open.rms.baidu.com的API，需要token
- vip_bns_strategy.go: 查询VIP的bns感知策略及半自动过单的阈值信息
- vip_domain.go: 查询VIP的域名信息
- vip_rs.go: 查询VIP的RS信息
- vip_ttm.go: 查询VIP的端口是否开启透传信息
- vip_user.go: 查询VIP的负责人信息
- vip_info.go: 查询VIP的使用期限、负载均衡策略、值班表、bns及bns感知策略信息
- noah_pdb.go: 查询noah id到Pdb id的对应关系
- rs.go: 查询rs信息
- vip_isp.go: 查询VIP的运营商信息
