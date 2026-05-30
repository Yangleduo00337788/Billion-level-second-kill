package com.seckill.im.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class ImGroupTableDef extends TableDef {

    public static final ImGroupTableDef IM_GROUP = new ImGroupTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn AVATAR = new QueryColumn(this, "avatar");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn OWNER_ID = new QueryColumn(this, "owner_id");

    public final QueryColumn GROUP_NAME = new QueryColumn(this, "group_name");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn MEMBER_COUNT = new QueryColumn(this, "member_count");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, AVATAR, STATUS, OWNER_ID, GROUP_NAME, CREATE_TIME, UPDATE_TIME, MEMBER_COUNT};

    public ImGroupTableDef() {
        super("", "im_group");
    }

    private ImGroupTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public ImGroupTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new ImGroupTableDef("", "im_group", alias));
    }

}
