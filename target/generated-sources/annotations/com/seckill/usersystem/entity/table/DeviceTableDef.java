package com.seckill.usersystem.entity.table;

import com.mybatisflex.core.query.QueryColumn;
import com.mybatisflex.core.table.TableDef;

// Auto generate by mybatis-flex, do not modify it.
public class DeviceTableDef extends TableDef {

    public static final DeviceTableDef DEVICE = new DeviceTableDef();

    public final QueryColumn ID = new QueryColumn(this, "id");

    public final QueryColumn OS_NAME = new QueryColumn(this, "os_name");

    public final QueryColumn STATUS = new QueryColumn(this, "status");

    public final QueryColumn USER_ID = new QueryColumn(this, "user_id");

    public final QueryColumn BROWSER = new QueryColumn(this, "browser");

    public final QueryColumn DEVICE_ID = new QueryColumn(this, "device_id");

    public final QueryColumn LOCATION = new QueryColumn(this, "location");

    public final QueryColumn IP_ADDRESS = new QueryColumn(this, "ip_address");

    public final QueryColumn IS_TRUSTED = new QueryColumn(this, "is_trusted");

    public final QueryColumn CREATE_TIME = new QueryColumn(this, "create_time");

    public final QueryColumn DEVICE_NAME = new QueryColumn(this, "device_name");

    public final QueryColumn DEVICE_TYPE = new QueryColumn(this, "device_type");

    public final QueryColumn UPDATE_TIME = new QueryColumn(this, "update_time");

    public final QueryColumn LAST_LOGIN_TIME = new QueryColumn(this, "last_login_time");

    public final QueryColumn LAST_ACTIVE_TIME = new QueryColumn(this, "last_active_time");

    /**
     * 所有字段。
     */
    public final QueryColumn ALL_COLUMNS = new QueryColumn(this, "*");

    /**
     * 默认字段，不包含逻辑删除或者 large 等字段。
     */
    public final QueryColumn[] DEFAULT_COLUMNS = new QueryColumn[]{ID, OS_NAME, STATUS, USER_ID, BROWSER, DEVICE_ID, LOCATION, IP_ADDRESS, IS_TRUSTED, CREATE_TIME, DEVICE_NAME, DEVICE_TYPE, UPDATE_TIME, LAST_LOGIN_TIME, LAST_ACTIVE_TIME};

    public DeviceTableDef() {
        super("", "sys_device");
    }

    private DeviceTableDef(String schema, String name, String alisa) {
        super(schema, name, alisa);
    }

    public DeviceTableDef as(String alias) {
        String key = getNameWithSchema() + "." + alias;
        return getCache(key, k -> new DeviceTableDef("", "sys_device", alias));
    }

}
