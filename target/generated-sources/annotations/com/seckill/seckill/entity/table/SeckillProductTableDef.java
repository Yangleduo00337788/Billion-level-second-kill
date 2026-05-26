package com.seckill.seckill.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class SeckillProductTableDef extends TableDef {

    public static final SeckillProductTableDef SECKILL_PRODUCT = new SeckillProductTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn END_TIME = new QueryColumn(this, "end_time");

    public final QueryColumn VERSION = new QueryColumn(this, "version");

    public final QueryColumn PRODUCT_ID = new QueryColumn(this, "product_id");

    public final QueryColumn START_TIME = new QueryColumn(this, "start_time");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn PERSON_LIMIT = new QueryColumn(this, "person_limit");

    public final QueryColumn SECKILL_PRICE = new QueryColumn(this, "seckill_price");

    public final QueryColumn SECKILL_STOCK = new QueryColumn(this, "seckill_stock");

    public final QueryColumn REDIS_STOCK_KEY = new QueryColumn(this, "redis_stock_key");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, STATUS, END_TIME, VERSION, PRODUCT_ID, START_TIME, CREATE_TIME, UPDATE_TIME, PERSON_LIMIT, SECKILL_PRICE, SECKILL_STOCK, REDIS_STOCK_KEY};

    public SeckillProductTableDef() {
        super("", "sk_seckill_product");
    }

    private SeckillProductTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public SeckillProductTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new SeckillProductTableDef("", "sk_seckill_product", alias));
    }

}
