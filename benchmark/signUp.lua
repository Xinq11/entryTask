--调用注册接口，向数据库导入1000w条测试数据
wrk.headers["Content-Type"] = "application/json"
local cnt = 9988367
local num = 7
function request()
    if (num >= 10)
    then
        num = 0
    end
    local body = '{"username": "test%s","password": "testpassword%s"}'
    body = string.format(body, cnt, num)
    print(body)
    cnt = cnt + 1
    num = num + 1
    return wrk.format(nil, nil, nil, body)
end