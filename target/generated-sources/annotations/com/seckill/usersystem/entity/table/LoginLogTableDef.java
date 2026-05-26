package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class LoginLogTableDef extends TableDef {

    public static final LoginLogTableDef LOGIN_LOG = new LoginLogTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn OS_NAME = new QueryColumn(this, "os_name");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn BROWSER = new QueryColumn(this, "browser");

    public final QueryColumn DEVICE_ID = new QueryColumn(this, "device_id");

    public final QueryColumn LOCATION = new QueryColumn(this, "location");

    public final QueryColumn USERNAME = new QueryColumn(this, "username");

    public final QueryColumn IP_ADDRESS = new QueryColumn(this, "ip_address");

    public final QueryColumn LOGIN_TIME = new QueryColumn(this, "login_time");

    public final QueryColumn LOGIN_TYPE = new QueryColumn(this, "login_type");

    public final QueryColumn USER_AGENT = new QueryColumn(this, "user_agent");

    public final QueryColumn DEVICE_TYPE = new QueryColumn(this, "device_type");

    public final QueryColumn FAIL_REASON = new QueryColumn(this, "fail_reason");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, OS_NAME, STATUS, USER_ID, BROWSER, DEVICE_ID, LOCATION, USERNAME, IP_ADDRESS, LOGIN_TIME, LOGIN_TYPE, USER_AGENT, DEVICE_TYPE, FAIL_REASON};

    public LoginLogTableDef() {
        super("", "sys_login_log");
    }

    private LoginLogTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public LoginLogTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new LoginLogTableDef("", "sys_login_log", alias));
    }

}
