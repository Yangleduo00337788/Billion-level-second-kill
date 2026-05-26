package com.seckill.im.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class ImSessionTableDef extends TableDef {

    public static final ImSessionTableDef IM_SESSION = new ImSessionTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn GROUP_ID = new QueryColumn(this, "group_id");

    public final QueryColumn LAST_MSG_ID = new QueryColumn(this, "last_msg_id");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn LAST_MSG_TIME = new QueryColumn(this, "last_msg_time");

    public final QueryColumn SESSION_TYPE = new QueryColumn(this, "session_type");

    public final QueryColumn LAST_MSG_CONTENT = new QueryColumn(this, "last_msg_content");

    public final QueryColumn LAST_MSG_SENDER_ID = new QueryColumn(this, "last_msg_sender_id");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, STATUS, GROUP_ID, LAST_MSG_ID, CREATE_TIME, UPDATE_TIME, LAST_MSG_TIME, SESSION_TYPE, LAST_MSG_CONTENT, LAST_MSG_SENDER_ID};

    public ImSessionTableDef() {
        super("", "im_session");
    }

    private ImSessionTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public ImSessionTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new ImSessionTableDef("", "im_session", alias));
    }

}
