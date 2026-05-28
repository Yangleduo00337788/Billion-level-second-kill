package com.seckill.im.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class ImMessageTableDef extends TableDef {

    public static final ImMessageTableDef IM_MESSAGE = new ImMessageTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn SEQ = new QueryColumn(this, "seq");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn CONTENT = new QueryColumn(this, "content");

    public final QueryColumn MSG_TYPE = new QueryColumn(this, "msg_type");

    public final QueryColumn SENDER_ID = new QueryColumn(this, "sender_id");

    public final QueryColumn SESSION_ID = new QueryColumn(this, "session_id");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn RECALL_TIME = new QueryColumn(this, "recall_time");

    public final QueryColumn CLIENT_MSG_ID = new QueryColumn(this, "client_msg_id");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, SEQ, STATUS, CONTENT, MSG_TYPE, SENDER_ID, SESSION_ID, CREATE_TIME, RECALL_TIME, CLIENT_MSG_ID};

    public ImMessageTableDef() {
        super("", "im_message");
    }

    private ImMessageTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public ImMessageTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new ImMessageTableDef("", "im_message", alias));
    }

}
