package com.seckill.im.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class ImSessionUserTableDef extends TableDef {

    public static final ImSessionUserTableDef IM_SESSION_USER = new ImSessionUserTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn SESSION_ID = new QueryColumn(this, "session_id");

    public final QueryColumn JOINED_TIME = new QueryColumn(this, "joined_time");

    public final QueryColumn LAST_READ_MSG_SEQ = new QueryColumn(this, "last_read_msg_seq");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, USER_ID, SESSION_ID, JOINED_TIME, LAST_READ_MSG_SEQ};

    public ImSessionUserTableDef() {
        super("", "im_session_user");
    }

    private ImSessionUserTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public ImSessionUserTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new ImSessionUserTableDef("", "im_session_user", alias));
    }

}
