wrk.headers["Content-Type"] = "application/json"    
local cnt = 0
function request()    
    local body = '{"username": "test%s","password": "testpassword%s"}'
    body = string.format(body, cnt, cnt)
    print(body)
    cnt = cnt + 1
    return wrk.format(nil, nil, nil, body)    
end 