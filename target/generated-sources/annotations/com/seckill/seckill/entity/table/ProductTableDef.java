package com.seckill.seckill.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class ProductTableDef extends TableDef {

    public static final ProductTableDef PRODUCT = new ProductTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn PRICE = new QueryColumn(this, "price");

    public final QueryColumn TITLE = new QueryColumn(this, "title");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn MAIN_IMAGE = new QueryColumn(this, "main_image");

    public final QueryColumn SOLD_COUNT = new QueryColumn(this, "sold_count");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn TOTAL_STOCK = new QueryColumn(this, "total_stock");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn DESCRIPTION = new QueryColumn(this, "description");

    public final QueryColumn PRODUCT_NAME = new QueryColumn(this, "product_name");

    public final QueryColumn ORIGINAL_PRICE = new QueryColumn(this, "original_price");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, PRICE, TITLE, STATUS, MAIN_IMAGE, SOLD_COUNT, CREATE_TIME, TOTAL_STOCK, UPDATE_TIME, DESCRIPTION, PRODUCT_NAME, ORIGINAL_PRICE};

    public ProductTableDef() {
        super("", "sk_product");
    }

    private ProductTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public ProductTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new ProductTableDef("", "sk_product", alias));
    }

}
