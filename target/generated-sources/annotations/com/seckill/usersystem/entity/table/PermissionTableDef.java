package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class PermissionTableDef extends TableDef {

    public static final PermissionTableDef PERMISSION = new PermissionTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn ICON = new QueryColumn(this, "icon");

    public final QueryColumn PATH = new QueryColumn(this, "path");

    public final QueryColumn SORT = new QueryColumn(this, "sort");

    public final QueryColumn METHOD = new QueryColumn(this, "method");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn PARENT_ID = new QueryColumn(this, "parent_id");

    public final QueryColumn IS_DELETED = new QueryColumn(this, "is_deleted");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn PERMISSION_CODE = new QueryColumn(this, "permission_code");

    public final QueryColumn PERMISSION_NAME = new QueryColumn(this, "permission_name");

    /**
     * 权限类型：1-菜单，2-按钮，3-接口
     */
    public final QueryColumn PERMISSION_TYPE = new QueryColumn(this, "permission_type");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, ICON, PATH, SORT, METHOD, STATUS, PARENT_ID, IS_DELETED, CREATE_TIME, UPDATE_TIME, PERMISSION_CODE, PERMISSION_NAME, PERMISSION_TYPE};

    public PermissionTableDef() {
        super("", "sys_permission");
    }

    private PermissionTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public PermissionTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new PermissionTableDef("", "sys_permission", alias));
    }

}
