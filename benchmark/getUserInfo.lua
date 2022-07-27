-- 获取用户信息接口脚本
local threads = {}
local counter = 1

-- 初始化数据
function init(args)
    thread_success = 0
    thread_fail = 0
end

function setup(thread)
-- 给每个线程设置一个 id 参数
   thread:set("id", counter)
-- 将线程添加到 table 中
   table.insert(threads, thread)
   counter = counter + 1
end

function response(status, headers, body)
-- 每得到一次请求的响应 判断请求是否成功
    if status == 200 then
           if string.find(body, "6") ~= nil then
                thread_success = thread_success + 1
           else
                thread_fail = thread_fail + 1
           end
    else
           thread_fail = thread_fail + 1
    end
end

-- 统计请求成功和失败的数量
function done(summary, latency, requests)
    local total_success = 0
    local total_fail = 0
    for _, thread in pairs(threads) do
        local thread_success = thread:get("thread_success")
        local thread_fail = thread:get("thread_fail")
        total_success = total_success + thread_success
        total_fail = total_fail + thread_fail
    end
    local success_msg = "total_success: %s"
    print(success_msg:format(total_success))
    local fail_msg = "total_fail: %s"
    print(fail_msg:format(total_fail))
end