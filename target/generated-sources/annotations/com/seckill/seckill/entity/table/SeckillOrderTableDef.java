package com.seckill.seckill.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class SeckillOrderTableDef extends TableDef {

    public static final SeckillOrderTableDef SECKILL_ORDER = new SeckillOrderTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn ORDER_NO = new QueryColumn(this, "order_no");

    public final QueryColumn PAY_TIME = new QueryColumn(this, "pay_time");

    public final QueryColumn QUANTITY = new QueryColumn(this, "quantity");

    public final QueryColumn PRODUCT_ID = new QueryColumn(this, "product_id");

    public final QueryColumn SECKILL_ID = new QueryColumn(this, "seckill_id");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn PRODUCT_NAME = new QueryColumn(this, "product_name");

    public final QueryColumn TOTAL_AMOUNT = new QueryColumn(this, "total_amount");

    public final QueryColumn PRODUCT_PRICE = new QueryColumn(this, "product_price");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, STATUS, USER_ID, ORDER_NO, PAY_TIME, QUANTITY, PRODUCT_ID, SECKILL_ID, CREATE_TIME, UPDATE_TIME, PRODUCT_NAME, TOTAL_AMOUNT, PRODUCT_PRICE};

    public SeckillOrderTableDef() {
        super("", "sk_order");
    }

    private SeckillOrderTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public SeckillOrderTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new SeckillOrderTableDef("", "sk_order", alias));
    }

}
