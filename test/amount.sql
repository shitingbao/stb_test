-- =================================================================================
-- 数据清洗脚本：将 budget_transfers 和 recharge_records_ad 数据迁移到新的钱包表结构
-- =================================================================================

BEGIN;

-- =================================================================================
-- 第一步：生成所有 wallet 钱包数据
-- 每个公司每个BC一条记录，CN_USD_01 和 CN_USD_02 合并为一条
-- =================================================================================

-- 情况1：country != '0' 的情况（指定BC），排除 CN_USD_01 和 CN_USD_02
INSERT INTO wallet (company_id, country, amount, created_at, updated_at)
SELECT DISTINCT
    d.dept_id AS company_id,
    bc.country AS country,
    0 AS amount,
    NOW() AS created_at,
    NOW() AS updated_at
FROM sys_dept d
INNER JOIN sys_dept_bc bc ON bc.dept_id = d.dept_id
WHERE d.company_type IN (2, 3)  -- 公司类型：二级代理(2)、直客(3)
    AND d.dept_type = 1         -- 部门类型：1代表二代的公司
    AND d.dept_id IS NOT NULL
    AND d.dept_id > 0
    AND bc.dept_id IS NOT NULL
    AND bc.dept_id > 0
    AND bc.country IS NOT NULL
    AND bc.country != ''
    AND bc.country != '0'
    AND bc.country NOT IN ('CN_USD_01', 'CN_USD_02')
    AND NOT EXISTS (
        SELECT 1 FROM wallet w 
        WHERE w.company_id = d.dept_id 
        AND w.country = bc.country
        AND w.company_id IS NOT NULL
        AND w.country IS NOT NULL
    );

-- 情况2：CN_USD_01 和 CN_USD_02 合并为一条（存储为 'CN_USD_01,CN_USD_02'）
INSERT INTO wallet (company_id, country, amount, created_at, updated_at)
SELECT DISTINCT
    d.dept_id AS company_id,
    'CN_USD_01,CN_USD_02' AS country,
    0 AS amount,
    NOW() AS created_at,
    NOW() AS updated_at
FROM sys_dept d
INNER JOIN sys_dept_bc bc ON bc.dept_id = d.dept_id
WHERE d.company_type IN (2, 3)  -- 公司类型：二级代理(2)、直客(3)
    AND d.dept_type = 1         -- 部门类型：1代表二代的公司
    AND d.dept_id IS NOT NULL
    AND d.dept_id > 0
    AND bc.dept_id IS NOT NULL
    AND bc.dept_id > 0
    AND bc.country IS NOT NULL
    AND bc.country IN ('CN_USD_01', 'CN_USD_02')
    AND NOT EXISTS (
        SELECT 1 FROM wallet w 
        WHERE w.company_id = d.dept_id 
        AND w.country = 'CN_USD_01,CN_USD_02'
        AND w.company_id IS NOT NULL
        AND w.country IS NOT NULL
    )
GROUP BY d.dept_id;

-- 情况3：country = '0' 的情况（全量BC）
INSERT INTO wallet (company_id, country, amount, created_at, updated_at)
SELECT DISTINCT
    d.dept_id AS company_id,
    CASE
        WHEN bc_all.name IN ('CN_USD_01', 'CN_USD_02') THEN 'CN_USD_01,CN_USD_02'
        ELSE bc_all.name
    END AS country,
    0 AS amount,
    NOW() AS created_at,
    NOW() AS updated_at
FROM sys_dept d
INNER JOIN sys_dept_bc bc_zero ON bc_zero.dept_id = d.dept_id 
    AND bc_zero.country = '0'
    AND bc_zero.dept_id IS NOT NULL
    AND bc_zero.dept_id > 0
    AND bc_zero.country IS NOT NULL
CROSS JOIN business_center bc_all
WHERE d.company_type IN (2, 3)  -- 公司类型：二级代理(2)、直客(3)
    AND d.dept_type = 1         -- 部门类型：1代表二代的公司
    AND d.dept_id IS NOT NULL
    AND d.dept_id > 0
    AND bc_all.name IS NOT NULL
    AND bc_all.name != ''
    AND bc_all.name != '0'
    AND NOT EXISTS (
        SELECT 1 FROM wallet w 
        WHERE w.company_id = d.dept_id 
        AND w.country = CASE
            WHEN bc_all.name IN ('CN_USD_01', 'CN_USD_02') THEN 'CN_USD_01,CN_USD_02'
            ELSE bc_all.name
        END
        AND w.company_id IS NOT NULL
        AND w.country IS NOT NULL
    );

-- =================================================================================
-- 第二步：迁移旧版本 budget_transfers 数据（没有对应 recharge_records_ad）
-- 一条 budget_transfers → wallet_recharge 两条记录（充值和使用）
-- 一条 budget_transfers → wallet_recharge_ad 一条记录
-- =================================================================================

-- 2.1 生成钱包充值和使用记录（一条 budget_transfers 生成两条记录：充值为正，使用为负）
INSERT INTO wallet_recharge (
    wallet_id, from_user_id, advertiser_id, amount, recharge_type, 
    status, request_id, remark, created_at, updated_at, approved_by, approved_at, is_pay
)
SELECT
    w.id AS wallet_id,
    bt.from_user_id,
    bt.target_id AS advertiser_id,
    bt.amount AS amount,  -- 充值记录：使用原始金额（正数）
    4 AS recharge_type,  -- 广告账户(充值)
    CASE bt.pay_status
        WHEN 0 THEN 3  -- pay_status = 0 旧版没有pay_status，代表成功 → 已完成
        WHEN 1 THEN 1  -- 1待审核 → 1待审核
        WHEN 2 THEN 3  -- 2已充值 → 3已完成
        WHEN 3 THEN 2  -- 3充值中 → 2充值中
        WHEN 4 THEN 6  -- 4审核拒绝 → 6审核拒绝
        WHEN 5 THEN 4  -- 5充值失败 → 4充值失败
        WHEN 6 THEN 4  -- 6部分失败 → 4充值失败
        ELSE 0  -- 默认值，方便比对数据
    END AS status,
    bt.request_id,
    bt.remark AS remark,
    bt.created_at,
    bt.updated_at,
    bt.approved_by,
    bt.approved_at,
    CASE 
        WHEN bt.pay_status = 0 THEN 1  -- pay_status = 0 旧版没有pay_status，代表成功，标记为已支付
        WHEN bt.pay_status = 2 THEN 1  -- pay_status = 2 已充值，标记为已支付
        ELSE 0 
    END AS is_pay
FROM budget_transfers bt
INNER JOIN business_center bc_bt ON bc_bt.bc_id = bt.bc_id  -- 通过 bc_id 找到对应的 name
INNER JOIN wallet w ON w.company_id = bt.company_id
    AND (
        (bc_bt.name IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%CN_USD%')
        OR (bc_bt.name NOT IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%' || bc_bt.name || '%')
    )
WHERE bt.target_id != '' AND bt.target_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM recharge_records_ad rra 
        WHERE rra.request_id = bt.request_id
    )
    AND NOT EXISTS (
        SELECT 1 FROM wallet_recharge wr 
        WHERE wr.request_id = bt.request_id AND wr.recharge_type = 4
    )
UNION ALL
SELECT
    w.id AS wallet_id,
    bt.from_user_id,
    bt.target_id AS advertiser_id,
    -bt.amount AS amount,  -- 使用记录：金额取负
    5 AS recharge_type,  -- 广告账户(退款/使用)
    CASE bt.pay_status
        WHEN 0 THEN 3  -- pay_status = 0 旧版没有pay_status，代表成功 → 已完成
        WHEN 1 THEN 1  -- 1待审核 → 1待审核
        WHEN 2 THEN 3  -- 2已充值 → 3已完成
        WHEN 3 THEN 2  -- 3充值中 → 2充值中
        WHEN 4 THEN 6  -- 4审核拒绝 → 6审核拒绝
        WHEN 5 THEN 4  -- 5充值失败 → 4充值失败
        WHEN 6 THEN 4  -- 6部分失败 → 4充值失败
        ELSE 0  -- 默认值，方便比对数据
    END AS status,
    bt.request_id,
    bt.remark AS remark,
    bt.created_at,
    bt.updated_at,
    bt.approved_by,
    bt.approved_at,
    CASE 
        WHEN bt.pay_status = 0 THEN 1  -- pay_status = 0 旧版没有pay_status，代表成功，标记为已支付
        WHEN bt.pay_status = 2 THEN 1  -- pay_status = 2 已充值，标记为已支付
        ELSE 0 
    END AS is_pay
FROM budget_transfers bt
INNER JOIN business_center bc_bt ON bc_bt.bc_id = bt.bc_id  -- 通过 bc_id 找到对应的 name
INNER JOIN wallet w ON w.company_id = bt.company_id
    AND (
        (bc_bt.name IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%CN_USD%')
        OR (bc_bt.name NOT IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%' || bc_bt.name || '%')
    )
WHERE bt.target_id != '' AND bt.target_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM recharge_records_ad rra 
        WHERE rra.request_id = bt.request_id
    )
    AND NOT EXISTS (
        SELECT 1 FROM wallet_recharge wr 
        WHERE wr.request_id = bt.request_id AND wr.recharge_type = 5
    );

-- 2.3 生成广告户流水记录（wallet_recharge_ad）- wallet_recharge_id暂时设为0，后续通过更新关联
INSERT INTO wallet_recharge_ad (
    wallet_recharge_id, advertiser_id, amount, status, request_id, created_at, updated_at
)
SELECT
    0 AS wallet_recharge_id,  -- 暂时设为0，后续通过更新关联
    bt.target_id AS advertiser_id,
    bt.amount AS amount,
    CASE bt.pay_status
        WHEN 0 THEN 3  -- pay_status = 0 旧版没有pay_status，代表成功 → 3已完成
        WHEN 2 THEN 3  -- 2已充值 → 3已完成
        WHEN 3 THEN 2  -- 3充值中 → 2充值中
        WHEN 4 THEN 4  -- 4审核拒绝 → 4充值失败
        WHEN 5 THEN 4  -- 5充值失败 → 4充值失败
        WHEN 6 THEN 4  -- 6部分失败 → 4充值失败
        ELSE 0  -- 默认成功（旧版没有pay_status）
    END AS status,
    bt.request_id,
    bt.created_at,
    bt.updated_at
FROM budget_transfers bt
WHERE bt.target_id != '' AND bt.target_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM recharge_records_ad rra 
        WHERE rra.request_id = bt.request_id
    )
    AND NOT EXISTS (
        SELECT 1 FROM wallet_recharge_ad wra 
        WHERE wra.request_id = bt.request_id 
        AND wra.advertiser_id = bt.target_id
    );

-- =================================================================================
-- 第三步：迁移新版本 budget_transfers 数据（有对应 recharge_records_ad）
-- 一条 budget_transfers → wallet_recharge 一条充值记录
-- recharge_records_ad 多条 → wallet_recharge_ad 多条记录
-- recharge_records_ad 多条 → wallet_recharge 多条使用记录
-- =================================================================================

-- 3.1 生成钱包充值记录（一次充值）
-- 注意：充值记录的状态应该根据所有使用记录的状态来决定
-- 如果所有使用记录都成功（recharge_status = 2），充值记录才算成功（status = 3）
-- 如果有任何使用记录失败，充值记录的状态应该根据 pay_status 来决定
INSERT INTO wallet_recharge (
    wallet_id, from_user_id, advertiser_id, amount, recharge_type, 
    status, request_id, remark, created_at, updated_at, approved_by, approved_at, is_pay
)
SELECT DISTINCT
    w.id AS wallet_id,
    bt.from_user_id,
    '' AS advertiser_id,  -- 新版本可能没有单个 advertiser_id
    bt.amount AS amount,
    CASE 
        WHEN bt.amount > 0 THEN 4  -- 广告账户(充值)
        ELSE 5  -- 广告账户(退款)
    END AS recharge_type,
    -- 充值记录的状态：如果所有使用记录都成功，则充值记录成功；否则根据 pay_status 决定
    CASE 
        -- 如果所有使用记录都成功（recharge_status = 2），充值记录才算成功
        WHEN (
            SELECT COUNT(*) FROM recharge_records_ad rra_check
            WHERE rra_check.request_id = bt.request_id
                AND rra_check.recharge_status != 2  -- 有非成功的记录
        ) = 0 
        AND (
            SELECT COUNT(*) FROM recharge_records_ad rra_check2
            WHERE rra_check2.request_id = bt.request_id
        ) > 0  -- 确保有使用记录
        THEN 3  -- 所有使用记录都成功，充值记录成功
        -- 否则根据 pay_status 决定
        WHEN bt.pay_status = 0 THEN 3  -- pay_status = 0 旧版没有pay_status，代表成功 → 已完成
        WHEN bt.pay_status = 1 THEN 1  -- 1待审核 → 1待审核
        WHEN bt.pay_status = 2 THEN 3  -- 2已充值 → 3已完成
        WHEN bt.pay_status = 3 THEN 2  -- 3充值中 → 2充值中
        WHEN bt.pay_status = 4 THEN 6  -- 4审核拒绝 → 6审核拒绝
        WHEN bt.pay_status = 5 THEN 4  -- 5充值失败 → 4充值失败
        WHEN bt.pay_status = 6 THEN 4  -- 6部分失败 → 4充值失败
        ELSE 0  -- 默认值，方便比对数据
    END AS status,
    bt.request_id,
    bt.remark AS remark,
    bt.created_at,
    bt.updated_at,
    bt.approved_by,
    bt.approved_at,
    CASE 
        -- 如果所有使用记录都成功，标记为已支付
        WHEN (
            SELECT COUNT(*) FROM recharge_records_ad rra_check
            WHERE rra_check.request_id = bt.request_id
                AND rra_check.recharge_status != 2
        ) = 0 
        AND (
            SELECT COUNT(*) FROM recharge_records_ad rra_check2
            WHERE rra_check2.request_id = bt.request_id
        ) > 0
        THEN 1
        WHEN bt.pay_status = 0 THEN 1  -- pay_status = 0 旧版没有pay_status，代表成功，标记为已支付
        WHEN bt.pay_status = 2 THEN 1  -- pay_status = 2 已充值，标记为已支付
        ELSE 0 
    END AS is_pay
FROM budget_transfers bt
INNER JOIN business_center bc_bt ON bc_bt.bc_id = bt.bc_id  -- 通过 bc_id 找到对应的 name
INNER JOIN wallet w ON w.company_id = bt.company_id
    AND (
        (bc_bt.name IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%CN_USD%')
        OR (bc_bt.name NOT IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%' || bc_bt.name || '%')
    )
WHERE EXISTS (
        SELECT 1 FROM recharge_records_ad rra 
        WHERE rra.request_id = bt.request_id
    )
    AND NOT EXISTS (
        SELECT 1 FROM wallet_recharge wr 
        WHERE wr.request_id = bt.request_id
    );

-- 3.2 生成广告户流水记录（wallet_recharge_ad）- wallet_recharge_id暂时设为0，后续通过更新关联
INSERT INTO wallet_recharge_ad (
    wallet_recharge_id, advertiser_id, amount, status, request_id, created_at, updated_at
)
SELECT
    0 AS wallet_recharge_id,  -- 暂时设为0，后续通过更新关联
    rra.advertiser_id,
    rra.amount,
    CASE rra.recharge_status
        WHEN 1 THEN 2  -- 待充值 1 → 2 充值中
        WHEN 2 THEN 3  -- 充值成功 2 → 3 已完成
        WHEN 3 THEN 4  -- 充值失败 3 → 4 充值失败
        WHEN 4 THEN 4  -- 已关闭 4 → 4 充值失败
        ELSE 0  -- 默认值，方便比对数据
    END AS status,
    rra.request_id,
    rra.created_at,
    rra.updated_at
FROM recharge_records_ad rra
WHERE rra.advertiser_id != '' AND rra.advertiser_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM wallet_recharge_ad wra 
        WHERE wra.request_id = rra.request_id 
        AND wra.advertiser_id = rra.advertiser_id
    );

-- 3.3 生成钱包使用记录（每个广告户充值对应一条使用记录）
INSERT INTO wallet_recharge (
    wallet_id, from_user_id, advertiser_id, amount, recharge_type, 
    status, request_id, remark, created_at, updated_at, approved_by, approved_at, is_pay
)
SELECT
    w.id AS wallet_id,
    bt.from_user_id,
    rra.advertiser_id,
    -rra.amount AS amount,  -- 使用记录：金额取负
    5 AS recharge_type,  -- 广告账户(退款/使用)
    CASE rra.recharge_status
        WHEN 1 THEN 2  -- 待充值 1 → 2 充值中
        WHEN 2 THEN 3  -- 充值成功 2 → 3 已完成
        WHEN 3 THEN 4  -- 充值失败 3 → 4 充值失败
        WHEN 4 THEN 4  -- 已关闭 4 → 4 充值失败
        ELSE 0  -- 默认值，方便比对数据
    END AS status,
    rra.request_id AS request_id,  -- 保持原始 request_id
    bt.remark AS remark,
    rra.created_at,
    rra.updated_at,
    bt.approved_by,
    CASE WHEN rra.recharge_status = 2 THEN rra.updated_at ELSE NULL END AS approved_at,
    CASE WHEN rra.recharge_status = 2 THEN 1 ELSE 0 END AS is_pay
FROM recharge_records_ad rra
INNER JOIN budget_transfers bt ON bt.request_id = rra.request_id
INNER JOIN business_center bc_bt ON bc_bt.bc_id = bt.bc_id  -- 通过 bc_id 找到对应的 name
INNER JOIN wallet w ON w.company_id = bt.company_id
    AND (
        (bc_bt.name IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%CN_USD%')
        OR (bc_bt.name NOT IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%' || bc_bt.name || '%')
    )
WHERE rra.advertiser_id != '' AND rra.advertiser_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM wallet_recharge wr 
        WHERE wr.request_id = rra.request_id AND wr.recharge_type = 5 AND wr.advertiser_id = rra.advertiser_id
    );

-- =================================================================================
-- 第四步：更新 wallet_recharge_ad 的 wallet_recharge_id（关联使用记录）
-- =================================================================================

-- 4.1 更新旧版本的 wallet_recharge_ad（关联使用记录）
UPDATE wallet_recharge_ad wra
SET wallet_recharge_id = wr.id
FROM budget_transfers bt
INNER JOIN business_center bc_bt ON bc_bt.bc_id = bt.bc_id  -- 通过 bc_id 找到对应的 name
INNER JOIN wallet w ON w.company_id = bt.company_id
    AND (
        (bc_bt.name IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%CN_USD%')
        OR (bc_bt.name NOT IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%' || bc_bt.name || '%')
    )
INNER JOIN wallet_recharge wr ON wr.wallet_id = w.id
    AND wr.request_id = bt.request_id AND wr.recharge_type = 5
    AND wr.advertiser_id = bt.target_id
WHERE wra.wallet_recharge_id = 0
    AND wra.request_id = bt.request_id
    AND wra.advertiser_id = bt.target_id
    AND NOT EXISTS (
        SELECT 1 FROM recharge_records_ad rra 
        WHERE rra.request_id = bt.request_id
    );

-- 4.2 更新新版本的 wallet_recharge_ad（关联钱包记录）
-- 根据广告户id和country，通过客户编号找到公司，再找到对应的钱包和钱包流水记录
UPDATE wallet_recharge_ad wra
SET wallet_recharge_id = wr.id
FROM recharge_records_ad rra
-- 通过广告户id和country找到客户编号
INNER JOIN th_tiktok_advertisers tta ON tta.advertiser_id = rra.advertiser_id
    AND (
        (rra.country IN ('CN_USD_01', 'CN_USD_02') AND tta.country IN ('CN_USD_01', 'CN_USD_02'))
        OR (rra.country NOT IN ('CN_USD_01', 'CN_USD_02') AND tta.country = rra.country)
    )
-- 通过客户编号找到公司
INNER JOIN sys_dept_customer sdc ON sdc.customer_id = tta.customer_id
-- 通过公司id和country找到钱包
INNER JOIN wallet w ON w.company_id = sdc.dept_id
    AND (
        (rra.country IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%CN_USD_01%' AND w.country LIKE '%CN_USD_02%')
        OR (rra.country NOT IN ('CN_USD_01', 'CN_USD_02') AND w.country LIKE '%' || rra.country || '%')
    )
-- 通过钱包id、广告户id找到对应的钱包流水记录
INNER JOIN wallet_recharge wr ON wr.wallet_id = w.id
    AND wr.advertiser_id = rra.advertiser_id
WHERE wra.wallet_recharge_id = 0
    AND rra.request_id = wra.request_id
    AND rra.advertiser_id = wra.advertiser_id
    AND wr.request_id = rra.request_id;

-- =================================================================================
-- 第五步：更新 wallet 的 amount（根据 wallet_recharge 计算，只更新成功的记录）
-- =================================================================================

UPDATE wallet w
SET amount = COALESCE((
    SELECT SUM(amount)
    FROM wallet_recharge wr
    WHERE wr.wallet_id = w.id
        AND wr.status = 3  -- 只统计已完成状态的记录（status = 3 表示已完成）
), 0),
updated_at = NOW();

COMMIT;

-- =================================================================================
-- 数据验证查询（可选，用于验证迁移结果）
-- =================================================================================
-- ROLLBACK;
-- 验证 wallet 数量
-- SELECT COUNT(*) AS wallet_count FROM wallet;

-- 验证 wallet_recharge 数量
-- SELECT COUNT(*) AS wallet_recharge_count FROM wallet_recharge;

-- 验证 wallet_recharge_ad 数量
-- SELECT COUNT(*) AS wallet_recharge_ad_count FROM wallet_recharge_ad;

-- 验证 wallet 金额总和
-- SELECT SUM(amount) AS total_amount FROM wallet;

-- 验证 wallet_recharge 金额总和
-- SELECT SUM(amount) AS total_recharge_amount FROM wallet_recharge;
