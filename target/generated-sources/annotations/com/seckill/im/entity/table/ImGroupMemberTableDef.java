package com.seckill.im.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class ImGroupMemberTableDef extends TableDef {

    public static final ImGroupMemberTableDef IM_GROUP_MEMBER = new ImGroupMemberTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn ROLE = new QueryColumn(this, "role");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn GROUP_ID = new QueryColumn(this, "group_id");

    public final QueryColumn JOINED_TIME = new QueryColumn(this, "joined_time");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, ROLE, USER_ID, GROUP_ID, JOINED_TIME};

    public ImGroupMemberTableDef() {
        super("", "im_group_member");
    }

    private ImGroupMemberTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public ImGroupMemberTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new ImGroupMemberTableDef("", "im_group_member", alias));
    }

}
