功能，生成所有二代公司钱包，线上充值上线后的广告户充值记录
1.相关表
budget_transfers 充值记录
recharge_records_ad 充值记录对应的广告户记录
business_center表，budget_transfers的bc_id对应该表bc_id，其他表的country对应该表的name
注意条件country 使用like

wallet 钱包
wallet_recharge 钱包流水
wallet_recharge_ad 广告户流水

2.生成钱包
先 生成所有 wallet 钱包数据，每个公司每个bc一条，除了 CN_USD_01和CN_USD_02这个是一条

3.数据版本，线上充值上线后的广告户充值记录
budget_transfers 数据有两个版本
旧版本一个是 budget_transfers 一个表的数据，
    表示对于一条 advertise_id 的直接充值操作，这里拆分为钱包的充值和使用，并对广告户充值的操作
新版本 budget_transfers对应 recharge_records_ad 有数据，
    对应了一次充值，recharge_records_ad 多条广告户充值记录

4.数据迁移
现在要将旧表的数据迁移到新表，钱包流水记录要根据之前的数据产生

budget_transfers 的没有对应recharge_records_ad，数据表示 advertise_id 对应 wallet_recharge 钱包流水 充值和使用两条数据，wallet_recharge_ad 广告户流水一条充值数据
budget_transfers 的 有对应为 wallet_recharge 的一次充值
recharge_records_ad 充值记录对应的广告户记录对应一条 wallet_recharge_ad 广告户流水记录，
    还对应 wallet_recharge 钱包流水使用数据和充值
    （相当于 wallet_recharge 的一次充值对应这里的多个广告户的多次使用）
budget_transfers状态通过pay_status对应，如果是0，也代表成功（旧版本默认都成功）

5.新旧表状态对应
budget_transfers 的 pay_status 对应wallet_recharge状态 status
1待审核 对应 1待审核
2已充值（或者0旧版没有pay_status）对应 3已完成
3充值中对应 2充值中
4审核拒绝对应 6审核拒绝
5充值失败 对应 4充值失败
6部分失败 对应 4充值失败

budget_transfers 的 pay_status 对应wallet_recharge_ad 状态 status
2已充值（或者0旧版没有pay_status）对应 3已完成
3充值中对应 2充值中
4审核拒绝对应 4充值失败
5充值失败 对应 4充值失败
6部分失败 对应 4充值失败

recharge_records_ad 的 recharge_status 对应 wallet_recharge 的 status
待充值 1 对应 2 充值中
充值成功 2 对应 3 已完成
充值失败 3 对应 4 充值失败
已关闭 4 对应 4 充值失败

recharge_records_ad 的 recharge_status 对应 wallet_recharge_ad 的 status
待充值 1 对应 2 充值中
充值成功 2 对应 3 已完成
充值失败 3 对应 4 充值失败
已关闭 4 对应 4 充值失败

钱包账单类型的对应
广告账户(充值) 4， 根据金额判断
广告账户(退款) 5，
