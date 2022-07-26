wrk.headers["Content-Type"] = "application/json"
local cnt = 8153443
local num = 3
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