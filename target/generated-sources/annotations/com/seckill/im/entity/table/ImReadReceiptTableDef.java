package com.seckill.im.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class ImReadReceiptTableDef extends TableDef {

    public static final ImReadReceiptTableDef IM_READ_RECEIPT = new ImReadReceiptTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn MSG_ID = new QueryColumn(this, "msg_id");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn READ_TIME = new QueryColumn(this, "read_time");

    public final QueryColumn SESSION_ID = new QueryColumn(this, "session_id");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, MSG_ID, USER_ID, READ_TIME, SESSION_ID};

    public ImReadReceiptTableDef() {
        super("", "im_read_receipt");
    }

    private ImReadReceiptTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public ImReadReceiptTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new ImReadReceiptTableDef("", "im_read_receipt", alias));
    }

}
