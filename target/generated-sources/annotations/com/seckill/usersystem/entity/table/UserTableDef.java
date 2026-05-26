package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class UserTableDef extends TableDef {

    public static final UserTableDef USER = new UserTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn EMAIL = new QueryColumn(this, "email");

    public final QueryColumn PHONE = new QueryColumn(this, "phone");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn PASSWORD = new QueryColumn(this, "password");

    public final QueryColumn USERNAME = new QueryColumn(this, "username");

    public final QueryColumn IS_DELETED = new QueryColumn(this, "is_deleted");

    public final QueryColumn LOGIN_TYPE = new QueryColumn(this, "login_type");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn LOGIN_COUNT = new QueryColumn(this, "login_count");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn LAST_LOGIN_IP = new QueryColumn(this, "last_login_ip");

    public final QueryColumn LAST_LOGIN_TIME = new QueryColumn(this, "last_login_time");

    public final QueryColumn PASSWORD_UPDATE_TIME = new QueryColumn(this, "password_update_time");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, EMAIL, PHONE, STATUS, PASSWORD, USERNAME, IS_DELETED, LOGIN_TYPE, CREATE_TIME, LOGIN_COUNT, UPDATE_TIME, LAST_LOGIN_IP, LAST_LOGIN_TIME, PASSWORD_UPDATE_TIME};

    public UserTableDef() {
        super("", "sys_user");
    }

    private UserTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public UserTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new UserTableDef("", "sys_user", alias));
    }

}
