package com.seckill.seckill.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class InventoryLogTableDef extends TableDef {

    public static final InventoryLogTableDef INVENTORY_LOG = new InventoryLogTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn REMARK = new QueryColumn(this, "remark");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn ORDER_ID = new QueryColumn(this, "order_id");

    public final QueryColumn PRODUCT_ID = new QueryColumn(this, "product_id");

    public final QueryColumn SECKILL_ID = new QueryColumn(this, "seckill_id");

    public final QueryColumn AFTER_STOCK = new QueryColumn(this, "after_stock");

    public final QueryColumn CHANGE_TYPE = new QueryColumn(this, "change_type");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn BEFORE_STOCK = new QueryColumn(this, "before_stock");

    public final QueryColumn CHANGE_AMOUNT = new QueryColumn(this, "change_amount");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, REMARK, USER_ID, ORDER_ID, PRODUCT_ID, SECKILL_ID, AFTER_STOCK, CHANGE_TYPE, CREATE_TIME, BEFORE_STOCK, CHANGE_AMOUNT};

    public InventoryLogTableDef() {
        super("", "sk_inventory_log");
    }

    private InventoryLogTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public InventoryLogTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new InventoryLogTableDef("", "sk_inventory_log", alias));
    }

}
