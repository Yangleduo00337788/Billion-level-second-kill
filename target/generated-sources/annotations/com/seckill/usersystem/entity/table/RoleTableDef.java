package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class RoleTableDef extends TableDef {

    public static final RoleTableDef ROLE = new RoleTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn SORT = new QueryColumn(this, "sort");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn ROLE_CODE = new QueryColumn(this, "role_code");

    public final QueryColumn ROLE_NAME = new QueryColumn(this, "role_name");

    public final QueryColumn IS_DELETED = new QueryColumn(this, "is_deleted");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn DESCRIPTION = new QueryColumn(this, "description");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, SORT, STATUS, ROLE_CODE, ROLE_NAME, IS_DELETED, CREATE_TIME, UPDATE_TIME, DESCRIPTION};

    public RoleTableDef() {
        super("", "sys_role");
    }

    private RoleTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public RoleTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new RoleTableDef("", "sys_role", alias));
    }

}
